package main

import (
	"github.com\kaisawind\message"
)

func main() {
	config := message.Config{
		Backend: message.MessageBackendRabbitMQ,
		Username: "user",
		Password: "bitnami",
		Hosts:"localhost:5672",
	}
}