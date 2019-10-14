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

	forever := make(chan bool)

	mq.Subscribe("logs", "", ProcessMessage, nil)

	logrus.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

// ProcessMessage ...
func ProcessMessage(msg message.Message, ctx interface{}) {
	logrus.Infoln(msg)
	logrus.Infoln(ctx)
}
