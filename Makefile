.SILENT:

.PHONY: deps build run migrate migrate-down migrate-status up down clean

APP_NAME = main
BUILD_DIR = build
MIGRATIONS_DIR = migrations
CONFIG_PATH = ./config.yaml
GOOSE_DRIVER = postgres
GOOSE_DB_STRING = postgres://postgres:2020@localhost:5432/test_user_profile?sslmode=disable

# Установка зависимостей
deps:
	go mod tidy

# Сборка бинарника
build: deps
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/main.go

# Запуск приложения (без Docker)
run: deps
	go run cmd/main.go

# Применить миграции
migrate:
	goose -dir $(MIGRATIONS_DIR)/pgsql $(GOOSE_DRIVER) "$(GOOSE_DB_STRING)" up

# Откатить последнюю миграцию
migrate-down:
	goose -dir $(MIGRATIONS_DIR)/pgsql $(GOOSE_DRIVER) "$(GOOSE_DB_STRING)" down

# Посмотреть статус миграций
migrate-status:
	goose -dir $(MIGRATIONS_DIR)/pgsql $(GOOSE_DRIVER) "$(GOOSE_DB_STRING)" status

# Собрать и запустить контейнеры
up:
	docker-compose up --build

# Остановить и удалить контейнеры + volume
down:
	docker-compose down -v

#Очистить собранный бинарник
clean:
	if exist $(BUILD_DIR) rmdir /s /q $(BUILD_DIR)