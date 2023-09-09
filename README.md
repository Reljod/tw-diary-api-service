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
3. Copy config/config.yml and name it `config.local.yml`
4. Update `config.local.yml`. Make sure that DB Creds in `docker-compose.yml`
is same in the config file.
5. Run:
    ```bash
    go mod init
    go run http/server/server.go -config=config.local.yml
    ```
