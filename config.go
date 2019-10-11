package message

// Type MQ type
type Type string

// 消息队列类型
const (
	MessageBackendKafka    Type = "kafka"
	MessageBackendRabbitMQ Type = "rabbitmq"
)

// Config 配置参数
type Config struct {
	Backend  Type
	Username string
	Password string
	Hosts    string
}
