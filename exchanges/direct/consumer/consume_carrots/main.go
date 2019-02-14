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
		"logs_direct", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	exit1(err, "Failed to declare an exchange")

	// queue args
	args := map[string]interface{}{
		"x-message-ttl":             int32(43200000),  // 12 hours
		"x-expires":                 int32(345600000), // 4 days
		"x-dead-letter-exchange":    "amq.direct",
		"x-dead-letter-routing-key": "roten-carrots",
	}

	q, err := ch.QueueDeclare(
		"carrots", // name
		false,     // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // no-wait
		args,      // arguments
	)
	exit1(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,        // queue name
		"carrots",     // routing key
		"logs_direct", // exchange
		false,
		nil)
	exit1(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	exit1(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	select {}
}
