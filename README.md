# Simple HTTP/3 Checker

Quic Checker is a Go-based tool for checking the status of websites using
HTTP/3 (QUIC) protocol. It utilizes a worker pool to efficiently handle
multiple site status checks concurrently.

## Features

- Concurrent site status checking using a worker pool
- Configurable number of workers
- Customizable expected status codes for each URL
- Uses QUIC protocol for site status checking

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/13excite/quic-checker.git
    ```

2. Navigate to the project directory:

    ```sh
    cd quic-checker
    ```

3. Install dependencies:

    ```sh
    go mod tidy
    ```

## Usage

### Build

```shell
make build
```

### Run

```shell
./quic-checker -c ./test.yml -v
```
