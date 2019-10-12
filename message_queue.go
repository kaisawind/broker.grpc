package message

import "fmt"

// MQ MessageQueue消息队列
type MQ interface {
	Publish() error
	Subscribe() error
	Close() error
}

// NewMQ 创建MessageQueue
func NewMQ(c *Config) (MQ, error) {
	switch c.Backend {
	case MessageBackendRabbitMQ:
		return NewRabbitMQ(c)
	case MessageBackendKafka:
	default:
	}
	return nil, fmt.Errorf("invalid message backend '%s'", c.Backend)
}
