MBUILD_TAGS=-tags 'no_mysql no_sqlite3 no_ydb'
MBINARY=migrator
USBINARY=url-shortener

deps:
	go mod tidy

build-m:
	go build $(MBUILD_TAGS) -o $(MBINARY) cmd/migrator/main.go

build-us:
	go build -o $(USBINARY) cmd/url-shortener/main.go

us-run: build-us
	@export $(shell sed 's/#.*//g' .env | xargs); \
	./$(USBINARY)

migr-up: build-m
	@export $(shell sed 's/#.*//g' .env | xargs); \
	./$(MBINARY)

migr-down: build-m
	@export $(shell sed 's/#.*//g' .env | xargs); \
	export MIGRATIONS_DIRECTION='down'; \
	./$(MBINARY)

create-migration:
ifndef NAME
	$(error "NAME is not set. Use make create-migration NAME=name")
endif
	@export $(shell sed 's/#.*//g' .env | xargs); \
	goose -dir $$MIGRATIONS_PATH create $(NAME) sql

migr-g-up:
	@export $(shell sed 's/#.*//g' .env | xargs); \
	export CS_URL_SHORTENER="postgres://$$POSTGRES_USER:$$POSTGRES_PASSWORD@$$POSTGRES_HOST:$$POSTGRES_PORT/$$POSTGRES_DB?sslmode=disable"; \
	goose -dir $$MIGRATIONS_PATH -table $$MIGRATIONS_TABLE) postgres $$CS_URL_SHORTENER up

migr-g-down:
	@export $(shell sed 's/#.*//g' .env | xargs); \
	export CS_URL_SHORTENER="postgres://$$POSTGRES_USER:$$POSTGRES_PASSWORD@$$POSTGRES_HOST:$$POSTGRES_PORT/$$POSTGRES_DB?sslmode=disable"; \
	goose -dir $$MIGRATIONS_PATH -table $$MIGRATIONS_TABLE) postgres $$CS_URL_SHORTENER down

migr-status:
	@export $(shell sed 's/#.*//g' .env | xargs); \
	export CS_URL_SHORTENER="postgres://$$POSTGRES_USER:$$POSTGRES_PASSWORD@$$POSTGRES_HOST:$$POSTGRES_PORT/$$POSTGRES_DB?sslmode=disable"; \
	goose -dir $$MIGRATIONS_PATH -table $$MIGRATIONS_TABLE) postgres $$CS_URL_SHORTENER status

clean-m:
	rm -f $(MBINARY)

clean-us:
	rm -f $(USBINARY)

sw-init:
	swag init -g cmd/url-shortener/main.go -o docs

lint:
	golangci-lint run ./...

g-mock-service:
	mockgen -source=internal/service/url_shortener_service.go -destination=internal/mocks/mock_service.go -package=mocks

test:
	go test -v ./...

help:
	@echo "A .env file with the required environment variables must be present for most commands to work."
	@echo "Available commands:"
	@echo "  make deps                            - Install dependencies (go mod tidy)"
	@echo "  make build-m                         - Build the migrator binary"
	@echo "  make build-us                        - Build the url-shortener binary"
	@echo "  make migr-up                         - Apply migrations using the migrator binary"
	@echo "  make migr-down                       - Roll back migrations using the migrator binary"
	@echo "  make create-migration NAME=name      - Create a new migration"
	@echo "  make migr-g-up                       - Apply migrations using goose"
	@echo "  make migr-g-down                     - Roll back the last migration using goose"
	@echo "  make migr-g-status                   - Check migration status"
	@echo "  make clean-m                         - Remove the migrator binary file"
	@echo "  make clean-us                        - Remove the url-shortener binary file"
	@echo "  make sw-init                         - Init swagger documentation"
	@echo "  make lint                            - Run the linter (golangci-lint)"
	@echo "  make g-mock-service                  - Generate mock_service"
	@echo "  make test                            - Run tests"
