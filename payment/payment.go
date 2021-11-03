package main

import (
	"encoding/json"
	"fmt"
	"log"
	"payment/queue"
	"time"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)


type Order struct {
	Uuid      string    `json:"uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	ProductId string    `json:"product_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at,string"`
}


func init() {
	// Carregando arquivo .env da pasta payment
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Falha ao carregar arquivo .env")
	}
}


func notifyPaymentProcessed(order Order, ch *amqp.Channel) {
	json, _ := json.Marshal(order)
	queue.Publisher(json, "payment_ex", "", ch)
	fmt.Println(string(json))
}


func main() {
	in := make(chan []byte)

	connection := queue.Connect()
	queue.Consumer("order_queue", connection, in)

	var order Order

	for payload := range in {
		json.Unmarshal(payload, &order)
		order.Status = "Aprovado"

		notifyPaymentProcessed(order, connection)

	}
}
