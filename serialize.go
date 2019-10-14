package message

import (
	"encoding/json"
	"encoding/xml"
	"errors"
)

// 预定义序列化参数
var (
	JSONContentType = "application/json"
	XMLContentTypen = "application/xml"
	HexContentType  = "application/octet-stream"
)

// Serialize 序列化
func Serialize(msg Message, contentType string) ([]byte, error) {
	switch contentType {
	case JSONContentType:
		return json.Marshal(msg)
	case XMLContentTypen:
		return xml.Marshal(msg)
	default:
		return nil, errors.New("invalid message codec")
	}
}

// Deserialize 反序列化
func Deserialize(buf []byte, contentType string, obj interface{}) error {
	switch contentType {
	case JSONContentType:
		return json.Unmarshal(buf, obj)
	case XMLContentTypen:
		return xml.Unmarshal(buf, obj)
	default:
		return errors.New("invalid message codec")
	}
}
