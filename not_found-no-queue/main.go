package main

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func exit1(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	rUSER := "bunny"
	rPASS := "test"
	rHOST := "my-rabbit"
	rPORT := "5672" // string to make it easy dealign with the env
	rVHOST := "hole"

	// read from ENV
	if e := os.Getenv("RABBITMQ_USER"); e != "" {
		rUSER = e
	}
	if e := os.Getenv("RABBITMQ_PASS"); e != "" {
		rPASS = e
	}
	if e := os.Getenv("RABBITMQ_HOST"); e != "" {
		rHOST = e
	}
	if e := os.Getenv("RABBITMQ_PORT"); e != "" {
		rPORT = e
	}
	if e := os.Getenv("RABBITMQ_VHOST"); e != "" {
		rVHOST = e
	}

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		rUSER, rPASS, rHOST, rPORT, rVHOST))
	exit1(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	exit1(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"my-exchange", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(
		"fake-carrots",
		"",            // routing key
		"my-exchange", // exchange
		false,
		nil,
	)
	fmt.Printf("err = %+v\n", err)
}
