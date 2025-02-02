package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/ruslan-onishchenko/go-test-task/internal/core/domain"
)

var (
	ErrQueueFull        = errors.New("queue is full")
	ErrQueueNotFound    = errors.New("queue not found")
	ErrTimeout          = errors.New("timeout")
	ErrMaxQueuesReached = errors.New("maximum number of queues reached")
)

type QueueRepositoryMemory struct {
	queues      map[string]*domain.Queue
	mu          sync.RWMutex // Используем RWMutex для лучшей конкурентности
	maxQueues   int
	maxMessages int
}

func NewQueueRepositoryMemory(maxQueues, maxMessages int) *QueueRepositoryMemory {
	return &QueueRepositoryMemory{
		queues:      make(map[string]*domain.Queue),
		maxQueues:   maxQueues,
		maxMessages: maxMessages,
	}
}

// Вспомогательная функция для получения очереди или её создания
func (r *QueueRepositoryMemory) getOrCreateQueue(queueName string) (*domain.Queue, error) {
	r.mu.RLock()
	queue, exists := r.queues[queueName]
	r.mu.RUnlock()

	// Если очередь не существует, создаем её
	if !exists {
		if err := r.CreateQueue(queueName); err != nil {
			return nil, ErrQueueNotFound
		}

		// Получаем эксклюзивный доступ после создания очереди
		r.mu.RLock()
		queue = r.queues[queueName]
		r.mu.RUnlock()
	}

	return queue, nil
}

// Функция для создания очереди
func (r *QueueRepositoryMemory) CreateQueue(name string) error {
	r.mu.Lock() // Эксклюзивный замок для записи
	defer r.mu.Unlock()

	// Если достигнут предел очередей
	if len(r.queues) >= r.maxQueues {
		return ErrMaxQueuesReached
	}

	// Если очередь не существует, создаем её
	if _, exists := r.queues[name]; !exists {
		r.queues[name] = &domain.Queue{
			Name:     name,
			Messages: make(chan string, r.maxMessages),
		}
	}
	return nil
}

// Функция для добавления сообщения в очередь
func (r *QueueRepositoryMemory) Enqueue(queueName string, message string) error {
	queue, err := r.getOrCreateQueue(queueName)
	if err != nil {
		return err
	}

	// Пытаемся добавить сообщение в очередь
	select {
	case queue.Messages <- message:
		return nil
	default:
		return ErrQueueFull // Если очередь переполнена
	}
}

// Функция для извлечения сообщения из очереди с таймаутом
func (r *QueueRepositoryMemory) Dequeue(queueName string, timeout time.Duration) (string, error) {
	queue, err := r.getOrCreateQueue(queueName)
	if err != nil {
		return "", err
	}

	// Пытаемся извлечь сообщение из очереди с таймаутом
	select {
	case msg := <-queue.Messages:
		return msg, nil
	case <-time.After(timeout):
		return "", ErrTimeout // Если истекло время ожидания
	}
}
