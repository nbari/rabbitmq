package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {

	type Queue struct {
		Name  string `json:name`
		VHost string `json:vhost`
	}

	manager := "http://127.0.0.1:15672/api/queues/"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", manager, nil)
	req.SetBasicAuth("brokeruser", "brokeruser-nbari")
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
