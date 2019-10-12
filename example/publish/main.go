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

	forever := make(chan bool)
	var after <-chan time.Time
loop:
	after = time.After(3 * time.Second)
	for {
		logrus.Infoln("每3s执行一次")
		select {
		case <-forever:
		case <-after:
			mq.Publish()
			goto loop
		}
	}

}
