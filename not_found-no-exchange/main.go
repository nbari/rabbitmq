package main

import (
	"fmt"
	"log"

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

	err = ch.QueueBind(
		"carrots",
		"",                      // routing key
		"my-undefined-exchange", // exchange
		false,
		nil,
	)
	fmt.Printf("err = %+v\n", err)
}
