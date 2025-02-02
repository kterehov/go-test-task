package main

import (
	"github.com/ruslan-onishchenko/go-test-task/internal/adapter/http"
	"github.com/ruslan-onishchenko/go-test-task/internal/adapter/repository"
	"github.com/ruslan-onishchenko/go-test-task/internal/config"
	"github.com/ruslan-onishchenko/go-test-task/internal/core/service"
)

func main() {
	cfg := config.LoadConfig()

	repo := repository.NewQueueRepositoryMemory(cfg.MaxQueues, cfg.MaxMessages)
	queueService := service.NewQueueService(repo)

	// Настройка маршрутов и запуск сервера
	http.SetupRoutes(queueService, cfg.DefaultTimeout)
	http.StartServer(cfg.Port)
}
