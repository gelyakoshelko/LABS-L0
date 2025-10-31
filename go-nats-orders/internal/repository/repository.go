package repository

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
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


type OrderRepository struct {
	db *pgxpool.Pool
}


func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}


func (r *OrderRepository) GetAll(ctx context.Context) map[string]Order {
	cache := make(map[string]Order)

	rows, err := r.db.Query(ctx, "SELECT orders FROM orders")
	if err != nil {
		log.Printf("warning: could not load orders from db: %v\n", err)
		return cache
	}
	defer rows.Close()

	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			log.Printf("error scanning row: %v\n", err)
			continue
		}

		var o Order
		if err := json.Unmarshal([]byte(data), &o); err != nil {
			log.Printf("error unmarshalling order JSON: %v\n", err)
			continue
		}

		cache[o.OrderUID] = o
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows error: %v\n", err)
	}

	return cache
}
