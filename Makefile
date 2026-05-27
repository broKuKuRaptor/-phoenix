.PHONY: build test clean run

# Сборка сервиса
build:
	go build -o bin/accounts ./cmd/accounts/

# Запуск тестов
test:
	go test ./... -v -count=1

# Запуск сервиса (SQLite в памяти)
run: build
	./bin/accounts --db sqlite://:memory:

# Очистка
clean:
	rm -rf bin/
