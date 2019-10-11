package message

// MQ MessageQueue消息队列
type MQ interface {
	Close()
}

// NewMQ 创建MessageQueue
func NewMQ(c *Config) {
	switch c.Backend {
	case MessageBackendRabbitMQ:
	case MessageBackendKafka:
	default:
	}
}
