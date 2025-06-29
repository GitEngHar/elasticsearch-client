## Hi there! 👋

### About this repository
This repository contains a simple Go program for interacting with Elasticsearch on macOS.

### Prerequisites
- Docker
- Go (≥1.18)

### Usage

1. Start Elasticsearch
    ```sh
    ./elastic-search-run.bash
    ```

2. Generate a password
    ```sh
    ./elastic-search-run.bash
    ```

3. Set the environment variable
    ```sh
    export ELASTIC_PASSWD=<YOUR_PASSWORD>
    ```

4. Run the Go program
    ```sh
    go run main.go
    ```

### Example output

    ```sh
    ~/Documents/elasticsearch-client git:[master]
    go run main.go
    Indexed document, status: 200 OK
    Found 1 hits:
    - bob: Hello Elastic 
    ```
