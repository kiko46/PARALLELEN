package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Stock struct {
	Company   string  `json:"company"`
	EventType string  `json:"eventType"`
	Price     float64 `json:"price"`
}

func main() {
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	queueName := os.Getenv("QUEUE_NAME")
	mongoURL := os.Getenv("MONGO_URL")

	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	clientOptions := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("stockmarket").Collection("stocks")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var stocks []Stock
			if err := json.Unmarshal(d.Body, &stocks); err != nil {
				log.Printf("Failed to unmarshal JSON: %v", err)
				continue
			}

			var total float64
			for _, stock := range stocks {
				total += stock.Price
			}
			avgPrice := total / float64(len(stocks))

			filter := bson.M{"company": stocks[0].Company}
			update := bson.M{
				"$set": bson.M{
					"company":  stocks[0].Company,
					"avgPrice": avgPrice,
				},
			}
			_, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
			if err != nil {
				log.Printf("Failed to update MongoDB: %v", err)
			}
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")
	<-forever
}
