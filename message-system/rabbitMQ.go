package messagesystem

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var ch *amqp.Channel

var wg sync.WaitGroup

func RabbitMQPublish() {
	// aurora.coder.zqf@gmail.com
	username := "guest"
	password := "guest"
	dialStr := fmt.Sprintf("amqp://%s:%s@localhost:5672/", username, password)
	conn, err := amqp.Dial(dialStr)
	defer conn.Close()
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, _ = conn.Channel()
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "张庆风爱杨诗颖!"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}

// 消费者
func RabbitMQConsume() {
	username := "guest"
	password := "guest"
	dialStr := fmt.Sprintf("amqp://%s:%s@localhost:5672/", username, password)
	conn, err := amqp.Dial(dialStr)
	defer conn.Close()
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, _ = conn.Channel()
	defer ch.Close() // RabbitMQ 通道
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	qmqpMsgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan string)
	defer close(forever)

	wg.Add(1)
	go func() {
		for msg := range qmqpMsgs {
			forever <- string(msg.Body)
		}
	}()

	msg := <-forever
	wg.Done()
	log.Printf("Received a message: %s", msg)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
