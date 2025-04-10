package messagesystem

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var ch *amqp.Channel

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

	body := "Hello World!"
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
	go func() {
		for msg := range qmqpMsgs {
			forever <- string(msg.Body)
		}
	}()

	msg := <-forever
	log.Printf("Received a message: %s", msg)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
