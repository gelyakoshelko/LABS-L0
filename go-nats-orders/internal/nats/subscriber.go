package nats

import (
	"encoding/json"
	"log"

	"github.com/nats-io/stan.go"
	"go-nats-orders/internal/repository"
)


type Subscriber struct {
	Conn   stan.Conn
	Repo   *repository.OrderRepository
	Cache  map[string]repository.Order
}


func NewSubscriber(conn stan.Conn, repo *repository.OrderRepository) *Subscriber {
	return &Subscriber{
		Conn:  conn,
		Repo:  repo,
		Cache: make(map[string]repository.Order),
	}
}


func (s *Subscriber) Subscribe(channel string) {
	_, err := s.Conn.Subscribe(channel, func(m *stan.Msg) {
		var o repository.Order
		if err := json.Unmarshal(m.Data, &o); err != nil {
			log.Printf("Invalid JSON received: %v", err)
			return
		}

		s.Cache[o.OrderUID] = o

		log.Printf("Order received and cached: %s", o.OrderUID)
	}, stan.DurableName("durable"))
	if err != nil {
		log.Fatalf("NATS subscribe error: %v", err)
	}
	log.Println("NATS connected and subscribed")
}
