package main

import (
	"fmt"
	"log"
	"sync"

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

	concurrent := 50
	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, concurrent)

	for i := 0; i < 100; i++ {
		semaphore <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() {
				<-semaphore
				wg.Done()
			}()
			q, err := ch.QueueDeclare(
				"",
				false, // durable
				false, // delete when usused
				true,  // exclusive
				false, // no-wait
				nil,   // arguments
			)
			if err != nil {
				fmt.Errorf("%s: %s", err, "Failed to declare a queue")
			} else {
				fmt.Printf("random queue name: %s\n", q.Name)
			}
		}()
	}

	wg.Wait()
}
