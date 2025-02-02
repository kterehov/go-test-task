package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ruslan-onishchenko/go-test-task/internal/core/service"
)

// GetQueueHandler обрабатывает запросы на извлечение сообщения из очереди
type GetQueueHandler struct {
	queueService   *service.QueueService
	defaultTimeout time.Duration
}

func NewGetQueueHandler(queueService *service.QueueService, defaultTimeout time.Duration) *GetQueueHandler {
	return &GetQueueHandler{queueService: queueService, defaultTimeout: defaultTimeout}
}

func (h *GetQueueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queueName := extractQueueName(r)
	timeout := parseTimeout(r, h.defaultTimeout)

	msg, err := h.queueService.Dequeue(queueName, timeout)
	if err != nil {
		handleDequeueError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, msg)
}

// handleDequeueError обрабатывает ошибки при получении сообщения из очереди
func handleDequeueError(w http.ResponseWriter, err error) {
	if errors.Is(err, service.ErrQueueTimeout) {
		// Если произошел timeout, возвращаем 404
		http.Error(w, fmt.Sprintf("Timeout while waiting for message: %v", err), http.StatusNotFound)
	} else {
		// Все остальные ошибки - внутренние ошибки
		http.Error(w, fmt.Sprintf("Failed to dequeue message: %v", err), http.StatusInternalServerError)
	}
}

// parseTimeout извлекает параметр timeout из запроса
func parseTimeout(r *http.Request, defaultTimeout time.Duration) time.Duration {
	timeoutParam := r.URL.Query().Get("timeout")
	if timeoutParam == "" {
		return defaultTimeout
	}

	if t, err := strconv.Atoi(timeoutParam); err == nil {
		return time.Duration(t) * time.Second
	}

	return defaultTimeout
}
