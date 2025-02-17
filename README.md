# URL Shortener

🚀 **URL Shortener** is a simple REST API service for shortening long URLs, built with Golang.

## 📌 Features

- Shorten long URLs into unique short links with the ability to set the expiration date of the short link.
- Redirect to the original URLs from shortened links.
- It is possible to track URL browsing statistics.

## 🔧 Installation & Setup

### Prerequisites

- Go 1.18+
- PostgreSQL (or another supported database)

### Clone the repository

```sh
git clone https://github.com/Greshnov-Ivan/url-shortener.git
cd url-shortener
```

### Install dependencies

```sh
go mod tidy
```

### Configure environment

- Copy the example configuration file and update it with your database credentials and salt for hash_id_configuration:

```sh
cp config/local.yaml.example config/local.yaml
```

- Specify the path (you can use config/local.yaml) to your configuration file in the environment variables:

```sh
export CONFIG_PATH=./config/local.yaml
```

### Run database migrations

- Set the DB connection string in the environment variables:

```sh
export DB_CONNECTION_STRING="host=localhost port=5432 dbname=url_shortener user=yourName password=yourPassword sslmode=disable"
```

- Run migrator:

```sh
make migrate-up
```

### Start the server

```sh
go run cmd/url-shortener/main.go
```

The server should now be running on `http://localhost:8080`.

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
│   │── lib/             # Utility and helper functions
│   │── mocks/           # Auto-generated mock files for testing
│   │── repository/      # Database access layer
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

