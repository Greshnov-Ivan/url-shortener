MBUILD_TAGS=-tags 'no_mysql no_sqlite3 no_ydb'
MBINARY=migrator
USBINARY=url-shortener
MIGRATIONS_PATH=migrations/
MIGRATIONS_TABLE=migrations
#DB_CONNECTION_STRING?="host=localhost port=5432 dbname=url_shortener user=postgres password=postgres sslmode=disable"

deps:
	go mod tidy

build-m:
	go build $(MBUILD_TAGS) -o $(MBINARY) cmd/migrator/main.go

build-us:
	go build -o $(USBINARY) cmd/url-shortener/main.go

# Применение миграций через бинарник
migrate-up: build-m
	./$(MBINARY) -connection-string=$(DB_DB_CONNECTION_STRING) -migrations-path=$(MIGRATIONS_PATH)

migrate-down: build-m
	./$(MBINARY) -connection-string=$(DB_DB_CONNECTION_STRING) -migrations-path=$(MIGRATIONS_PATH) down

# Создание новой миграции с именем NAME
create-migration:
ifndef NAME
	$(error "name не задан. Используйте make create-migration NAME=название")
endif
	goose -dir $(MIGRATIONS_PATH) create $(NAME) sql

# Применение миграций напрямую через goose
migr-up:
	goose -dir $(MIGRATIONS_PATH) -table $(MIGRATIONS_TABLE) postgres $(DB_CONNECTION_STRING) up

# Откат миграции через goose
migr-down:
	goose -dir $(MIGRATIONS_PATH) -table $(MIGRATIONS_TABLE) postgres $(DB_CONNECTION_STRING) down

# Статус миграций
migr-status:
	goose -dir $(MIGRATIONS_PATH) -table $(MIGRATIONS_TABLE) postgres $(DB_CONNECTION_STRING) status

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
	@echo "Доступные команды:"
	@echo "  make deps                            - Установить зависимости (go mod tidy)"
	@echo "  make build-m                         - Собрать бинарник migrator"
	@echo "  make build-us                        - Собрать бинарник url-shortener"
	@echo "  make migrate-up                      - Применить миграции через бинарник migrator"
	@echo "  make migrate-down                    - Откатить миграции через бинарник migrator"
	@echo "  make create-migration NAME=название  - Создать новую миграцию"
	@echo "  make migr-up                         - Применить миграции через goose"
	@echo "  make migr-down                       - Откатить последнюю миграцию через goose"
	@echo "  make migr-status                     - Проверить статус миграций"
	@echo "  make clean-m                         - Удалить бинарный файл migrator"
	@echo "  make clean-us                        - Удалить бинарный файл url-shortener"
	@echo "  make lint                            - Запустить линтер (golangci-lint)"
	@echo "  make test                            - Запустить тесты"
