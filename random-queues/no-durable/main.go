package main

import (
	"fmt"
	"log"
	"os"
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

	args := map[string]interface{}{
		//  ----------------------------------------------------------------------------
		// This creates a race condition
		// to test it increment the value of concurrent to something bigger, for example
		//
		// concurrent := 500
		//
		// it will return somethig like:
		//
		// Failed to register a consumer: Exception (404) Reason: "NOT_FOUND - no queue 'carrot-2019-02-25-eb8f43ec-0f23-4b82-b271-041518363b3a' in vhost 'hole'"
		// Failed to register a consumer: Exception (504) Reason: "channel/connection is not open"
		//
		// The reason for this is that the response speed is proportional to the current
		// load, therefore the more concurrent clients, the more it will take time to
		// accept the clients for consuming and when the client finally is ready for
		// consuming the queue it could be expired.
		//
		// A solution for this is to either not expire the queue or better use exclusive queues.
		"x-message-ttl": int32(3000), //
		"x-expires":     int32(8000), // <-- better remove this, increment it or use an exclusive queue
		//  ----------------------------------------------------------------------------
		"x-dead-letter-exchange":    "amq.direct",
		"x-dead-letter-routing-key": "carrot",
	}

	// concurrent := 500
	concurrent := 300

	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, concurrent)

	for i := 0; i < 1000; i++ {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(j int) {
			start := time.Now()
			queueName := fmt.Sprintf("carrot-%s-%s", time.Now().Format("2006-01-02"), uuid.Must(uuid.NewV4()))
			fmt.Printf("Creating queue: %s\n", queueName)
			defer func() {
				<-semaphore
				wg.Done()
			}()
			q, err := ch.QueueDeclare(
				queueName,
				false, // durable
				false, // delete when usused
				false, // exclusive
				false, // no-wait
				args,  // arguments
			)
			exit1(err, "Failed to declare a queue")

			fmt.Printf("elapsed: %v msg: %d\n", time.Since(start), j)

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
		}(i)
	}
	wg.Wait()
}
