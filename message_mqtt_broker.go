package message

// MqttBroker mqtt 消息结构体
type MqttBroker struct {
	exchange    string // Exchange name
	contentType string // 消息参数
	Topic       string `json:"topic" xml:"topic"`     // mqtt topic
	Type        string `json:"type" xml:"type"`       // mqtt type
	Payload     []byte `json:"payload" xml:"payload"` // mqtt value
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

// MessageType 获取消息参数
func (p *MqttBroker) MessageType() string { return MqttBrokerType }

// PackageType 获取消息参数
func (p *MqttBroker) PackageType() string { return p.Type }

// SetPackageType 设置消息参数
func (p *MqttBroker) SetPackageType(packetType string) Message {
	p.Type = packetType
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
