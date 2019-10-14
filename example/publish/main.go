package main

import (
	"time"

	"github.com/kaisawind/message"
	"github.com/sirupsen/logrus"
)

func main() {
	config := &message.Config{
		Backend:  message.MessageBackendRabbitMQ,
		Username: "user",
		Password: "bitnami",
		Hosts:    "localhost:5672",
	}

	mq, err := message.NewMQ(config)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	logrus.Infoln("NewMQ success")
	defer mq.Close()

	var after <-chan time.Time
loop:
	after = time.After(3 * time.Millisecond)
	for {
		logrus.Infoln("每3s执行一次")
		select {
		case <-after:
			timenow := time.Now().String()
			msg := &message.MqttBroker{
				Topic:   "/bkvvg3eegkqnbutdv2q0/test01/Clients",
				Payload: []byte(timenow),
			}
			msg.SetExchange("logs").
				SetContentType(message.JSONContentType).
				SetPackageType("mqtt.publish")
			mq.Publish(msg)
			goto loop
		}
	}

}
