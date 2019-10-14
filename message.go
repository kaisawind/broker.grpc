package message

// Message 消息接口
type Message interface {
	Exchange() string                          // Exchange name
	SetExchange(name string) Message           // set Exchange name
	ContentType() string                       // 获取消息参数
	SetContentType(contentType string) Message // 设置消息参数
	Type() string                              // 协议类型
	SetType(packetType string) Message         // 设置协议类型
	Serialize() ([]byte, error)                // 数据序列化
	Deserialize(buf []byte, opt string) error  // 消息反序列化
}

// MqttBroker mqtt 消息结构体
type MqttBroker struct {
	exchange    string // Exchange name
	contentType string // 消息参数
	Topic       string `json:"topic" xml:"topic"`             // mqtt topic
	PacketType  string `json:"packet_type" xml:"packet_type"` // mqtt type
	Payload     []byte `json:"payload" xml:"payload"`         // mqtt value
}

// Exchange Exchange name
func (p *MqttBroker) Exchange() string { return p.exchange }

// SetExchange 设置Exchange
func (p *MqttBroker) SetExchange(name string) Message {
	p.exchange = name
	return p
}

// ContentType 获取消息参数
func (p *MqttBroker) ContentType() string { return p.contentType }

// SetContentType 设置消息参数
func (p *MqttBroker) SetContentType(contentType string) Message {
	p.contentType = contentType
	return p
}

// Type 获取消息参数
func (p *MqttBroker) Type() string { return p.PacketType }

// SetType 设置消息参数
func (p *MqttBroker) SetType(packetType string) Message {
	p.PacketType = packetType
	return p
}

// Serialize 数据序列化
func (p *MqttBroker) Serialize() ([]byte, error) {
	return Serialize(p, p.contentType)
}

// Deserialize 消息反序列化
func (p *MqttBroker) Deserialize(buf []byte, contentType string) error {
	return Deserialize(buf, contentType, p)
}
