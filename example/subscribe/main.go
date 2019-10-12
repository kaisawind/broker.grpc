package main

import (
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

	mq.Subscribe()
}
