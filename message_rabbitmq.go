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
func NewRabbitMQ(c *Config) (*RabbitMQ, error) {
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

	go func() {
		err := <-conn.NotifyClose(make(chan *amqp.Error))
		logrus.WithError(err).Errorln("mq closing")
		// reconnect 断线重连(自愈)
		rabbitmq.Reconnect()
	}()

	return rabbitmq, err
}

// Reconnect 重连
func (mq *RabbitMQ) Reconnect() error {
	rabbitmq, err := NewRabbitMQ(mq.config)
	if err != nil {
		logrus.WithError(err).Errorln("reconnect error")
		return mq.Reconnect()
	}

	mq = rabbitmq
	return nil
}

// Close 断开连接
func (mq *RabbitMQ) Close() {
	if mq.conn.IsClosed() {
		logrus.Infoln("connection has been closed")
		return
	}

	// will close() the deliveries channel
	err := mq.ch.Cancel(mq.consumerTag, true)
	if err != nil {
		logrus.WithError(err).Errorln("channel closed failed")
		return
	}

	err = mq.conn.Close()
	if err != nil {
		logrus.WithError(err).Errorln("connection closed failed")
		return
	}
}
