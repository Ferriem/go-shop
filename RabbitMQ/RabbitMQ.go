package RabbitMQ

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/Ferriem/go-web/Product/datamodels"
	"github.com/Ferriem/go-web/Product/services"
	"github.com/streadway/amqp"
)

const MQURL = "amqp://user:user@localhost:5672/ferriem"

type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string
	Exchange  string
	Key       string
	Mqurl     string
	sync.Mutex
}

// create a instance of RabbitMQ
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
}

// close the channel and connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

// fail on error
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatal("%s: %s", message, err)
		panic(fmt.Sprintf("%s: %s", message, err))
	}
}

// create a simple mode of RabbitMQ
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error

	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect to RabbitMQ")

	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")

	return rabbitmq
}

// publish a message
func (r *RabbitMQ) PublishSimple(message string) error {
	r.Lock()
	defer r.Unlock()
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		// durable
		false,
		// delete when unused
		false,
		// exclusive
		false,
		// no-wait
		false,
		nil,
	)

	if err != nil {
		return err
	}

	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		// mandatory
		false,
		// immediate
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	return nil
}

func (r *RabbitMQ) ConsumeSimple(orderService services.IOrderService, productService services.IProductService) {
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	msgs, err := r.channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)

	r.channel.Qos(
		1,
		0,
		false,
	)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			message := &datamodels.Message{}
			err := json.Unmarshal([]byte(d.Body), &message)
			if err != nil {
				fmt.Println(err)
			}
			_, err = orderService.InsertOrderByMessage(message)
			if err != nil {
				fmt.Println(err)
			}

			err = productService.SubNumberOne(message.ProductID)
			if err != nil {
				fmt.Println(err)
			}
			//if true ack all non_ack message
			//if false ack only this message
			d.Ack(false)
		}
	}()

	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// create a subscription mod of RabbitMQ
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, "")
	var err error

	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect to RabbitMQ")

	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")

	return rabbitmq
}

// Publish in the subscription mode
func (r *RabbitMQ) PublishPub(message string) {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	r.failOnErr(err, "Failed to publish a message")
}

// Consume in the subscription mode
func (r *RabbitMQ) ReceiveSub() {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	err = r.channel.QueueBind(
		q.Name,
		"",
		r.Exchange,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to bind a queue")

	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// create a routing mode of RabbitMQ
func NewRabbitMQRouting(exchangeName, routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error

	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect to RabbitMQ")

	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")

	return rabbitmq
}

// Publish in the routing mode
func (r *RabbitMQ) PublishRouting(message string) {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// Consume in the routing mode
func (r *RabbitMQ) ReceiveRouting() {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to bind a queue")

	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// create a topic mode of RabbitMQ
func NewRabbitMQTopic(exchangeName, routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error

	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect to RabbitMQ")

	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")

	return rabbitmq
}

// Publish in the topic mode
func (r *RabbitMQ) PublishTopic(message string) {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// Consume in the topic mode
func (r *RabbitMQ) ReceiveTopic() {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to bind a queue")

	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
