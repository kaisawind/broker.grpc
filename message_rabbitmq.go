package message

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// RabbitMQ RabbitMQ结构体
type RabbitMQ struct {
	config      *Config // 配置信息
	conn        *amqp.Connection
	ch          *amqp.Channel
	consumerTag string
}

// NewRabbitMQ ...
func NewRabbitMQ(c *Config) (MQ, error) {
	rabbitmq, err := ConnectRabbitMQ(c)
	if err != nil {
		logrus.WithError(err).Errorln("create mq connection failed")
		return nil, err
	}
	go func() {
		err := <-rabbitmq.conn.NotifyClose(make(chan *amqp.Error))
		logrus.WithError(err).Errorln("mq closing")
		// reconnect 断线重连(自愈)
		rabbitmq.Reconnect()
	}()

	return rabbitmq, err
}

// ConnectRabbitMQ rabbitmq connect
func ConnectRabbitMQ(c *Config) (*RabbitMQ, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s/", c.Username, c.Password, c.Hosts)
	logrus.Infoln("create mq connection, mq server url : ", url)

	conn, err := amqp.Dial(url)
	if err != nil {
		logrus.WithError(err).Errorf("connect with mq server failed")
		return nil, err
	}
	logrus.Infoln("connect with mq server success url : ", url)

	ch, err := conn.Channel()
	if err != nil {
		logrus.WithError(err).Errorf("failed to create a mq channel")
		return nil, err
	}
	logrus.Infoln("create a mq channel success")

	rabbitmq := &RabbitMQ{
		config: c,
		conn:   conn,
		ch:     ch,
	}
	return rabbitmq, nil
}

// Reconnect 重连
func (mq *RabbitMQ) Reconnect() error {
	rabbitmq, err := ConnectRabbitMQ(mq.config)
	if err != nil {
		logrus.WithError(err).Errorln("reconnect error")
		return mq.Reconnect()
	}

	// 重设属性
	mq.conn = rabbitmq.conn
	mq.ch = rabbitmq.ch

	go func() {
		err := <-rabbitmq.conn.NotifyClose(make(chan *amqp.Error))
		logrus.WithError(err).Errorln("mq closing")
		// reconnect 断线重连(自愈)
		rabbitmq.Reconnect()
	}()
	return nil
}

// Close 断开连接
func (mq *RabbitMQ) Close() error {
	if mq.conn.IsClosed() {
		logrus.Infoln("connection has been closed")
		return nil
	}

	// will close() the deliveries channel
	err := mq.ch.Cancel(mq.consumerTag, true)
	if err != nil {
		logrus.WithError(err).Errorln("channel closed failed")
		return err
	}

	err = mq.conn.Close()
	if err != nil {
		logrus.WithError(err).Errorln("connection closed failed")
		return err
	}
	return nil
}

// Publish 发送数据
func (mq *RabbitMQ) Publish(msg Message) error {

	topic := msg.Topic()

	body, err := msg.Serialize(JSONOption)
	if err != nil {
		logrus.WithError(err).Errorf("serialize message '%s' failed", topic)
		return err
	}

	err = mq.ch.ExchangeDeclare(
		topic,    // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		logrus.Errorln(err, "Failed to declare an exchange")
		return err
	}

	err = mq.ch.Publish(
		topic, // exchange
		"",    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: JSONOption.Format,
			Body:        body,
		})
	return err
}

// Subscribe 订阅数据
func (mq *RabbitMQ) Subscribe() error {
	err := mq.ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		logrus.Errorln(err, "Failed to declare an exchange")
		return err
	}

	q, err := mq.ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		logrus.Errorln(err, "Failed to declare a queue")
		return err
	}

	err = mq.ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	if err != nil {
		logrus.Errorln(err, "Failed to bind a queue")
		return err
	}

	msgs, err := mq.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logrus.Errorln(err, "Failed to register a consumer")
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			logrus.Infoln(" [x] ", d)
		}
	}()

	logrus.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
	return nil
}
