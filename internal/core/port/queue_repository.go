package port

import "time"

type QueueRepository interface {
	CreateQueue(name string) error
	Enqueue(queueName string, message string) error
	Dequeue(queueName string, timeout time.Duration) (string, error)
}
