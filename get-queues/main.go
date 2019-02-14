package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	rUSER := "bunny"
	rPASS := "test"
	rHOST := "my-rabbit"
	rPORT := "15672" // string to make it easy dealign with the env

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

	type Queue struct {
		Name  string `json:name`
		VHost string `json:vhost`
	}

	manager := fmt.Sprintf("http://%s:%s/api/queues/", rHOST, rPORT)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", manager, nil)
	req.SetBasicAuth(rUSER, rPASS)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	value := make([]Queue, 0)
	json.NewDecoder(resp.Body).Decode(&value)
	for _, queue := range value {
		fmt.Printf("vhost: %s queue: %s\n", queue.VHost, queue.Name)
	}
}
