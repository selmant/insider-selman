# Project Name

This project is a Docker-based application that includes a webhook service and a MySQL database. The project is written in Go and can be easily started using Docker Compose.

## Requirements

- Docker
- Docker Compose
- Go (Golang)

## Setup

### Step 1: Clone the Repository

<sh>
git clone https://github.com/selmant/insider-selman.git
cd insider-selman
</sh>

### Step 2: Download Required Go Packages

```sh
go mod download
```

### Step 3: Start Docker Services

```sh
docker-compose up --build
```

This command will create and start the necessary Docker containers.

### Step 4: Initialize the Database

When the database container starts, the commands in the `init.sql` file will be executed to create the necessary database and tables.

### Step 5: Swagger Setup

You can use the Swagger interface to document the API. Follow these steps to set up Swagger:

1. Install the `swag` package:

```sh
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Generate Swagger files:

```sh
swag init -g cmd/message_sender.go
```

Swagger UI will be started along with the message sender and will be accessible at [http://localhost:8081/swagger](http://localhost:8081/swagger).

## Usage

The webhook service runs at `http://localhost:8080`. The MySQL database runs at `localhost:3306` and contains a database named `insider`.

### API Endpoints

- `POST /messages/job-state`: Changes the job state of a message sender.
- `GET /messages`: Lists messages.
- `POST /messages/queue`: Creates a new message.


### Environment Variables

The project uses environment variables defined in the `.env` file:

- `MYSQL_HOST`: Database host address
- `MYSQL_PORT`: Database port
- `MYSQL_USER`: Database username
- `MYSQL_PASSWORD`: Database password
- `MYSQL_DATABASE`: Database name
- `WEBHOOK_BASE_URL`: Webhook service URL
- `WEBHOOK_USERNAME`: Webhook username
- `WEBHOOK_PASSWORD`: Webhook password
- `MESSAGE_SENDER_PORT`: Message sender port

## File Structure

- `cmd/`: Contains application entry points.
    - `message_sender.go`: Message sender service.
    - `mock_webhook.go`: Mock webhook service.
- `config/`: Contains configuration files.
    - `init.sql`: Database initialization file.
- `internal/`: Contains internal application packages.
- `pkg/`: Contains external application packages.
- `Dockerfile`: Docker image creation file.
- `docker-compose.yml`: Docker Compose configuration file.
- `.env`: Environment variables file.