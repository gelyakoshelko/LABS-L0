package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type Order struct {
	OrderUID        string      `json:"order_uid"`
	TrackNumber     string      `json:"track_number"`
	Entry           string      `json:"entry"`
	Delivery        interface{} `json:"delivery"`
	Payment         interface{} `json:"payment"`
	Items           interface{} `json:"items"`
	Locale          string      `json:"locale"`
	InternalSig     string      `json:"internal_signature"`
	CustomerID      string      `json:"customer_id"`
	DeliveryService string      `json:"delivery_service"`
	Shardkey        string      `json:"shardkey"`
	SMID            int         `json:"sm_id"`
	DateCreated     string      `json:"date_created"`
	OofShard        string      `json:"oof_shard"`
}

var ordersMap map[string]Order

func main() {

	file, err := os.ReadFile("model.json")
	if err != nil {
		log.Fatalf("Cannot read model.json: %v", err)
	}

	var order Order
	if err := json.Unmarshal(file, &order); err != nil {
		log.Fatalf("Invalid JSON in model.json: %v", err)
	}

	ordersMap = make(map[string]Order)
	ordersMap[order.OrderUID] = order
	log.Printf("Cache loaded: 1 order")


	http.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/orders/")
		if o, ok := ordersMap[id]; ok {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(o)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "order not found"})
		}
	})


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {

			http.ServeFile(w, r, "./static"+r.URL.Path)
			return
		}
		http.ServeFile(w, r, "./static/index.html")
	})

	log.Println("HTTP server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
