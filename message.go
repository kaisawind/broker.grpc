package message

// Message 消息接口
type Message interface {
	Topic() string        // Exchange name
	SetTopic(name string) // set Exchange name
	Serialize(opt SerializeOption) ([]byte, error)  // 数据序列化
	Deserialize(buf []byte, opt SerializeOption) error	// 消息反序列化
}
