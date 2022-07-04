package main

import (
	"flag"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"time"
)

var config = Config{}
var debugMode *bool
var useMonFile *bool

// declare color for console output
var colorGreen = "\033[32m"
var colorReset = "\033[0m"
var colorRed = "\033[31m"

type HTTPResp struct {
	URL  string
	Resp *http.Response
	Err  error
}

// writeToMonFile writes byte data to mon file
func writeToMonFile(fileName string, data []byte) error {
	f, err := os.OpenFile(fileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}

// getHostFromURL extracts hostname from url
func getHostFromURL(url string) string {
	rp := regexp.MustCompile(`https:\/\/([a-z0-9.\-]+)\/.*`)
	return rp.FindStringSubmatch(url)[1]
}

func checkResponseData(sunResponse *HTTPResp, chatNotifyString *string, verbose bool) {
	if verbose {
		log.Print(colorGreen, sunResponse.URL, " ", sunResponse.Resp.Status)
	}
	// if requests failed
	if sunResponse.Err != nil {

		if verbose {
			log.Print(colorRed, "Monitoring request error", colorReset)
			log.Print(colorRed, "HTTP/3 check error on url: "+getHostFromURL(sunResponse.URL), sunResponse.Err, colorReset)
		}

		*chatNotifyString += fmt.Sprintf("HTTP/3 check failed %s: %v\n", sunResponse.URL, sunResponse.Err)

		return
	} // if err sunResponse

	// if err is nil, then check status code
	if sunResponse.Resp.Status != config.ExpectedStatusCode {

		if verbose {
			log.Print(colorRed, "Monitoring status code", colorReset)
			errString := fmt.Sprintf("HTTP/3 invalid response code: %s  URL:%s\n", sunResponse.Resp.Status, sunResponse.URL)
			log.Print(colorRed, errString, colorReset)
		}

		*chatNotifyString += fmt.Sprintf("HTTP/3 invalid response code %s  %s\n", sunResponse.Resp.Status, sunResponse.URL)
		return
	} // if status code checker

}

// httpCheckWorker is worker process. Need to run from checkQuicSun function
func httpCheckWorker(in <-chan string, resultStatus chan<- *HTTPResp, quicConf quic.Config) {
	roundTripper := &http3.RoundTripper{
		QuicConfig: &quicConf,
	}
	defer roundTripper.Close()

	hClient := &http.Client{
		Transport: roundTripper,
		Timeout:   time.Second * 4,
	}

	for addr := range in {
		if *debugMode {
			fmt.Println(addr)
		}

		rsp, err := hClient.Get(addr)
		// often, we get Client.Timeout error, because UDP datagram was drop. Try to avoid this error
		if err != nil {
			// try to check error type, if it's a Client.Timeout , then try to create
			// one more request, redeclare response and error vars and send to channel
			switch err := err.(type) {
			case net.Error:
				if err.Timeout() {
					log.Printf("This was a net.Error with a Timeout: %v\n", addr)
					rsp, err := hClient.Get(addr)
					resultStatus <- &HTTPResp{addr, rsp, err}
					runtime.Gosched()
					continue
				}
			}
		}
		resultStatus <- &HTTPResp{addr, rsp, err}
		runtime.Gosched()
	}
}

// checkQuicSun is a main checker logic
func checkQuicSun(urlsArr []string) []*HTTPResp {
	var qconf quic.Config
	var quicVer []quic.VersionNumber
	var responses []*HTTPResp

	// add quic version to quic config
	quicVer = append(quicVer, quic.VersionDraft29)
	qconf.Versions = quicVer

	qconf.MaxIdleTimeout = time.Second * 10

	runtime.GOMAXPROCS(0)
	workerInput := make(chan string, 3)
	// need big buffer size(when we are using 100+ urls in the config),
	// because we need to save data before the main tread LOOP will start
	resultWorker := make(chan *HTTPResp, 500)
	// run worker async
	for i := 0; i < config.GoroutinesCount; i++ {
		go httpCheckWorker(workerInput, resultWorker, qconf)
	}

	for _, addr := range urlsArr {
		workerInput <- addr

	}
LOOP:
	for {
		select {
		case r := <-resultWorker:
			responses = append(responses, r)
			if len(responses) == len(urlsArr) {
				break LOOP
			} // end if
		} // end select
	} // end for
	defer close(workerInput)
	defer close(resultWorker)

	return responses
}

func init() {
	var configPath string
	flag.StringVar(&configPath, "c", "", "usage -c config")
	debugMode = flag.Bool("v", false, "usage -v (verbose)")
	useMonFile = flag.Bool("w", false, "usage -w (write data to monfile)")

	flag.Parse()
	config.GetConfig(configPath)
}

func main() {

	sunResponces := checkQuicSun(config.Urls)
	var chatMsg string

	for _, HTTPResponse := range sunResponces {
		checkResponseData(HTTPResponse, &chatMsg, *debugMode)
	}

	if *useMonFile {
		err := writeToMonFile(config.MonFile, []byte(chatMsg))
		if err != nil {
			log.Print("Couldn't write data to the monitoring file", err)
		}
	}
}
