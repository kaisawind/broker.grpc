package message

import (
	"encoding/json"
	"errors"
)

// 预定义序列化参数
var (
	JSONOption = SerializeOption{Format: "application/json"}
	XMLOption  = SerializeOption{Format: "application/xml"}
	TextOption = SerializeOption{Format: "text/plain"}
)

// SerializeOption 序列化参数
type SerializeOption struct {
	Format string
}

// Serialize 序列化
func Serialize(msg Message, opt SerializeOption) ([]byte, error) {
	switch opt {
	case JSONOption:
		return json.Marshal(msg)
	default:
		return nil, errors.New("invalid message codec")
	}
}

// Deserialize 反序列化
func Deserialize(buf []byte, opt SerializeOption, obj interface{}) error {
	switch opt {
	case JSONOption:
		return json.Unmarshal(buf, obj)
	default:
		return errors.New("invalid message codec")
	}
}
