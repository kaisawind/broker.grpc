package message

// Factory 消息处理工厂
type Factory interface {
	Create(messageType string) Message
}

// DefaultFactory 默认工厂
type DefaultFactory struct{}

// Create 默认工厂的创建
func (p *DefaultFactory) Create(messageType string) Message {
	switch messageType {
	case MqttBrokerType:
		fallthrough
	default:
		return &MqttBroker{}
	}
}
