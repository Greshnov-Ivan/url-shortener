# URL Shortener

🚀 **URL Shortener** is a simple REST API service for shortening long URLs, built with Golang.

## 📌 Features

- Shorten long URLs into unique short links with the ability to set the expiration date of the short link.
- Redirect to the original URLs from shortened links.
- It is possible to track URL browsing statistics.

## 🔧 Installation & Setup

### Prerequisites

- Go 1.23+
- PostgreSQL (if you don't use docker-compose)
- Docker (if you use docker-compose)

### Clone the repository

```sh
git clone https://github.com/Greshnov-Ivan/url-shortener.git
cd url-shortener
```

### Install dependencies

```sh
go mod tidy
```

### Configure migrator

- Copy the example .env file and update it:

```sh
cp .env.example .env
```

### Run database migrations

Run migrator:

```sh
make migr-up
```

The configurator will apply migrations to your database

### Configure server

- Copy the example configuration file and update it with your database credentials and salt for hash_id_configuration:

```sh
cp config/local.yaml.example config/local.yaml
```

- Specify the path (you can use config/local.yaml) to your configuration file in the environment variables:

```sh
export CONFIG_PATH=./config/local.yaml
```

### Start the server

```sh
make us-run
```

The server should now be running on `http://localhost:8080`.

### Run everything you need with Docker Compose

To start the server along with the database and migrations, use:

```sh
docker-compose up --build
```

This will:
- Start a PostgreSQL database container
- Run the migrations
- Launch the server

Once everything is up, the server should be available at http://localhost:8080.

To stop the containers, run:

```sh
docker-compose down
```

## 📡 API Endpoints

### Shorten a URL

**Request:**

```sh
curl -X 'POST' \
  'http://localhost:8080/links' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "expiresAt": "2025-05-15T00:00:00.000000Z",
  "sourceUrl": "https://github.com/Greshnov-Ivan/url-shortener"
}'
```

**Response:**

```json
{
  "shortUrl": "http://localhost:8080/DdgJK"
}
```

### Retrieve Original URL

**Request:**

```sh
curl -X GET http://localhost:8080/DdgJK
```

Redirects to `https://github.com/Greshnov-Ivan/url-shortener`.

## 🏗 Project Structure

```
url-shortener/
│── cmd/                 # Application entry points
│   └── migrator/        # Migration service application entry point
│   └── url-shortener/   # Main service application entry point
│── config/              # Local configuration files (YAML)
│── docs/                # Swagger documentation
│── internal/            # Application logic
│   │── app/             # Application initialization
│   │── config/          # Configuration structure
│   │── entity/          # Domain entities and data structures
│   │── http/            # HTTP layer (handlers, middleware)
│   │   │── dto/         # Data Transfer Objects for HTTP requests and responses 
│   │   │── handlers/    # Route handlers for processing requests
│   │   │── middleware/  # Middleware functions
│   │── lib/             # Utility and helper functions
│   │── mocks/           # Auto-generated mock files for testing
│   │── repository/      # Database access layer
│   │   │── dto/         # Data Transfer Objects for database interactions 
│   │── service/         # Business logic and application rules (there are unit tests)
│── migrations/          # Database migration files
│── tests/               # Functional tests
│── web/                 # HTML pages
│── Makefile             # Build and management commands
│── go.mod               # Go module dependencies
│── go.sum               # Dependency checksums
```

## 🛠 Development & Testing

### Run tests

```sh
go test ./...
```

### Lint code

```sh
golangci-lint run
```

## 👤 Author

Developed by [Greshnov-Ivan](https://github.com/Greshnov-Ivan). Contributions are welcome!

