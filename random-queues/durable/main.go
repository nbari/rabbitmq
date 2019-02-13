package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
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

	args := map[string]interface{}{
		"x-message-ttl":             int32(3000),
		"x-expires":                 int32(8000),
		"x-dead-letter-exchange":    "amq.direct",
		"x-dead-letter-routing-key": "carrot",
	}

	concurrent := 50
	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, concurrent)

	for i := 0; i < 1000; i++ {
		semaphore <- struct{}{}
		wg.Add(1)
		go func() {
			queueName := fmt.Sprintf("carrot-%s-%s", time.Now().Format("2006-01-02"), uuid.Must(uuid.NewV4()))
			fmt.Printf("Creating queue: %s\n", queueName)
			defer func() {
				<-semaphore
				wg.Done()
			}()
			q, err := ch.QueueDeclare(
				queueName,
				true,  // durable
				false, // delete when usused
				false, // exclusive
				false, // no-wait
				args,  // arguments
			)
			exit1(err, "Failed to declare a queue")

			_, err = ch.Consume(
				q.Name, // queue
				"",     // consumer
				true,   // auto-ack
				false,  // exclusive
				false,  // no-local
				false,  // no-wait
				nil,    // args
			)
			exit1(err, "Failed to register a consumer")
		}()
	}
	wg.Wait()
}
