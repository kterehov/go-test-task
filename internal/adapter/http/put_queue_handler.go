package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ruslan-onishchenko/go-test-task/internal/core/service"
)

// PutQueueHandler обрабатывает запросы на добавление сообщения в очередь
type PutQueueHandler struct {
	queueService *service.QueueService
}

func NewPutQueueHandler(queueService *service.QueueService) *PutQueueHandler {
	return &PutQueueHandler{queueService: queueService}
}

func (h *PutQueueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queueName := extractQueueName(r)

	message, err := decodeMessage(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.queueService.Enqueue(queueName, message); err != nil {
		if errors.Is(err, service.ErrQueueFull) {
			http.Error(w, fmt.Sprintf("Queue %s is full: %v", queueName, err), http.StatusBadRequest)
		} else {
			http.Error(w, fmt.Sprintf("Failed to enqueue message: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// extractQueueName извлекает имя очереди из URL
func extractQueueName(r *http.Request) string {
	return r.URL.Path[len(QueuePath):]
}

// decodeMessage декодирует JSON-запрос и возвращает сообщение
func decodeMessage(r *http.Request) (string, error) {
	var body struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Message == "" {
		return "", errors.New("invalid message")
	}
	return body.Message, nil
}
