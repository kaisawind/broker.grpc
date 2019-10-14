package message

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// RabbitMQ RabbitMQ结构体
type RabbitMQ struct {
	config   *Config // 配置信息
	conn     *amqp.Connection
	ch       *amqp.Channel
	consumer *rabbitConsumer
	factory  Factory // 消息创建工厂函数
}

// rabbitConsumer 订阅者数据
type rabbitConsumer struct {
	Exchange string      // 路由名
	Queue    string      // 队列名
	Handler  HandlerFunc // 消息处理函数
	Context  interface{} // Context
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

	// 重新订阅数据
	if mq.consumer != nil {
		err = mq.Subscribe(mq.consumer.Exchange, mq.consumer.Queue, mq.consumer.Handler, mq.consumer.Context)
		if err != nil {
			logrus.WithError(err).Errorln("resubscribe error")
		}
		logrus.Infoln("resubscribe success")
	}
	return nil
}

// Close 断开连接
func (mq *RabbitMQ) Close() error {
	if mq.conn.IsClosed() {
		logrus.Infoln("connection has been closed")
		return nil
	}

	// will close() the deliveries channel
	err := mq.ch.Close()
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

// SetMessageFactory 消息工厂函数
func (mq *RabbitMQ) SetMessageFactory(factory Factory) {
	mq.factory = factory
}

// Publish 发送数据
func (mq *RabbitMQ) Publish(msg Message) error {

	exchange := msg.Exchange()
	contentType := msg.ContentType()
	messageType := msg.MessageType()

	body, err := msg.Serialize()
	if err != nil {
		logrus.WithError(err).Errorf("serialize message '%s' failed", exchange)
		return err
	}

	err = mq.ch.ExchangeDeclare(
		exchange, // name
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
		exchange, // exchange
		"",       // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: contentType,
			Timestamp:   time.Now(),
			Type:        messageType,
			Body:        body,
		})
	return err
}

// Subscribe 订阅数据
// exchange: exchange名，将从exchange获取数据
// queue: 队列名，当为""时，会是随机队列名(amq.gen-_JNtxpzXTx1Ic5AC0c-TvA)
// ctx: context
func (mq *RabbitMQ) Subscribe(exchange string, queue string, handler HandlerFunc, ctx interface{}) error {
	err := mq.ch.ExchangeDeclare(
		exchange, // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		logrus.Errorln(err, "Failed to declare an exchange", exchange)
		return err
	}

	q, err := mq.ch.QueueDeclare(
		queue, // name
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
		q.Name,   // queue name
		"",       // routing key is ignore when use fanout
		exchange, // exchange
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

	go func() {
		defer logrus.Infoln("go function exited")
		for d := range msgs {
			if mq.factory == nil {
				logrus.Infoln("factory is nil, set factory to default")
				mq.factory = &DefaultFactory{}
			}
			m := mq.factory.Create(d.Type)
			err := m.Deserialize(d.Body, d.ContentType)
			if err != nil {
				logrus.WithError(err).Errorln("Deserialize", d.Exchange, d.ConsumerTag, d.Type)
			} else {
				if handler != nil {
					handler(m, ctx)
				} else {
					logrus.Warningln("message handler is nil")
				}
			}
		}
	}()

	mq.consumer = &rabbitConsumer{
		Exchange: exchange,
		Queue:    queue,
		Handler:  handler,
		Context:  ctx,
	}
	logrus.Infoln("Subscribe success", exchange, queue)
	return nil
}
