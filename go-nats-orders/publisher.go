package main

import (
	"io/ioutil"
	"log"

	"github.com/nats-io/stan.go"
)

func main() {

	sc, err := stan.Connect("test-cluster", "publisher-1", stan.NatsURL("nats://127.0.0.1:4222"))
	if err != nil {
		log.Fatalf("NATS connect error: %v", err)
	}
	defer sc.Close()


	data, err := ioutil.ReadFile("model.json")
	if err != nil {
		log.Fatalf("Cannot read model.json: %v", err)
	}


	if err := sc.Publish("orders", data); err != nil {
		log.Fatalf("Publish error: %v", err)
	}

	log.Println("Order published successfully!")
}
