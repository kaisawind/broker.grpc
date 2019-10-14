package message

// MessageType 消息类型
const (
	MqttBrokerType = "mqtt-broker"
)

// HandlerFunc Message处理函数
type HandlerFunc func(msg Message, ctx interface{})

// Message 消息接口
type Message interface {
	Exchange() string                                 // Exchange name
	SetExchange(name string) Message                  // set Exchange name
	ContentType() string                              // 获取消息参数
	SetContentType(contentType string) Message        // 设置消息参数
	MessageType() string                              // 消息类型
	PackageType() string                              // 协议类型
	SetPackageType(packetType string) Message         // 设置协议类型
	Serialize() ([]byte, error)                       // 数据序列化
	Deserialize(buf []byte, contentType string) error // 消息反序列化
}
