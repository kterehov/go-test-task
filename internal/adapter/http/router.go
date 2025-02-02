package http

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ruslan-onishchenko/go-test-task/internal/core/service"
)

const QueuePath = "/queue/"

func SetupRoutes(queueService *service.QueueService, defaultTimeout time.Duration) {
	putQueueHandler := NewPutQueueHandler(queueService)
	getQueueHandler := NewGetQueueHandler(queueService, defaultTimeout)

	http.HandleFunc(QueuePath, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			putQueueHandler.ServeHTTP(w, r)
		case http.MethodGet:
			getQueueHandler.ServeHTTP(w, r)
		default:
			http.Error(w, "Method not found", http.StatusMethodNotAllowed)
		}
	})
}

func StartServer(port int) {
	log.Printf("Starting server on port %d...\n", port)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
