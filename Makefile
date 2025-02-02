# Makefile

# Параметры
PORT := 8080
MAX_QUEUES := 100
MAX_MESSAGES := 1000
DEFAULT_TIMEOUT := 30

# Цель по умолчанию
.PHONY: run
run:
	go run cmd/app/main.go --port $(PORT) --max-queues $(MAX_QUEUES) --max-messages $(MAX_MESSAGES) --default-timeout $(DEFAULT_TIMEOUT)
