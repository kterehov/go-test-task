package config

import (
	"flag"
	"time"
)

type Config struct {
	Port           int
	MaxQueues      int
	MaxMessages    int
	DefaultTimeout time.Duration
}

func LoadConfig() *Config {
	port := flag.Int("port", 8080, "Port to run the server on")
	maxQueues := flag.Int("max-queues", 100, "Maximum number of queues")
	maxMessages := flag.Int("max-messages", 1000, "Maximum number of messages per queue")
	defaultTimeout := flag.Int("default-timeout", 30, "Default timeout in seconds for GET requests")
	flag.Parse()

	return &Config{
		Port:           *port,
		MaxQueues:      *maxQueues,
		MaxMessages:    *maxMessages,
		DefaultTimeout: time.Duration(*defaultTimeout) * time.Second,
	}
}
