package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func exit1(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://bunny:test@my-rabbit:5672/hole")
	exit1(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	exit1(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_direct", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	exit1(err, "Failed to declare an exchange")

	t := time.Now().UTC()
	body := fmt.Sprintf("sending a 🥕 at %s\n", t)
	err = ch.Publish(
		"logs_direct", // exchange
		"carrots",     // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	exit1(err, "Failed to publish")

	body = fmt.Sprintf("sending a 🌿  at %s\n", t)
	err = ch.Publish(
		"logs_direct", // exchange
		"herbs",       // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	exit1(err, "Failed to publish")
}
