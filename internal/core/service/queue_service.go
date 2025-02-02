package service

import (
	"errors"
	"time"

	"github.com/ruslan-onishchenko/go-test-task/internal/adapter/repository"
	"github.com/ruslan-onishchenko/go-test-task/internal/core/port"
)

var (
	ErrQueueFull    = errors.New("service: queue is full")
	ErrQueueTimeout = errors.New("service: timeout during dequeue")
	ErrUnexpected   = errors.New("service: unexpected error")
)

type QueueService struct {
	repo port.QueueRepository
}

func NewQueueService(repo port.QueueRepository) *QueueService {
	return &QueueService{repo: repo}
}

func (s *QueueService) Enqueue(queueName string, message string) error {
	err := s.repo.Enqueue(queueName, message)
	if err != nil {
		if errors.Is(err, repository.ErrQueueFull) {
			return ErrQueueFull
		}
		return ErrUnexpected
	}
	return nil
}

func (s *QueueService) Dequeue(queueName string, timeout time.Duration) (string, error) {
	msg, err := s.repo.Dequeue(queueName, timeout)
	if err != nil {
		if errors.Is(err, repository.ErrTimeout) {
			return "", ErrQueueTimeout
		}
		return "", ErrUnexpected
	}
	return msg, nil
}
