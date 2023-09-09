# Twitter Diary API Service

Twitter Diary API Service built with Go's Gin Web Framework

## Development Guide

1. Clone the repo (*I think we all know how to do this LOL*)
2. Spin up a PostgreSQL Server, the recommended way is to use Docker Compose:
    - Download Docker Desktop
    - Run:
        ```bash
        docker compose up -d

        # To stop:
        docker compose down
        ```
3. Download DB Migration/Versioning Tool, `goose` via:
    - `go install github.com/pressly/goose/v3/cmd/goose@latest`  or
    - `brew install goose`
4. If running for the first time, add `.env` file or export `GOOSE_DRIVER` and `GOOSE_DBSTRING`. The instructions are in the `gw` file.<br>Then run:
    ```bash
    chmod 755 ./gw
    ./gw up
    ./gw status
    ```
    
5. Copy config/config.yml and name it `config.local.yml`
6. Update `config.local.yml`. Make sure that DB Creds in `docker-compose.yml`
is same in the config file.
7. Run:
    ```bash
    go mod init
    go run http/server/server.go -config=config.local.yml
    ```

