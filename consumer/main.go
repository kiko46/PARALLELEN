package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

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
	mongoURL := os.Getenv("MONGO_URL")
	mongoDB := os.Getenv("MONGODB_DB")
	mongoCollection := os.Getenv("MONGODB_COLLECTION")
	queueNames := strings.Split(os.Getenv("QUEUE_NAMES"), ",")

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

	clientOptions := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	database := client.Database(mongoDB)
	collection := database.Collection(mongoCollection)

	for _, queueName := range queueNames {
		go consumeQueue(queueName, ch, collection)
	}

	log.Printf("Waiting for messages. To exit press CTRL+C")
	select {}
}

func consumeQueue(queueName string, ch *amqp.Channel, collection *mongo.Collection) {
	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register a consumer for queue %s: %v", queueName, err)
	}

	var prices []float64
	for d := range msgs {
		var stock Stock
		if err := json.Unmarshal(d.Body, &stock); err != nil {
			log.Printf("Failed to unmarshal JSON from queue %s: %v", queueName, err)
			continue
		}
		prices = append(prices, stock.Price)
		if len(prices) == 1000 {
			total := 0.0
			for _, price := range prices {
				total += price
			}
			avgPrice := total / float64(len(prices))
			filter := bson.M{"company": stock.Company}
			update := bson.M{"$set": bson.M{"company": stock.Company, "avgPrice": avgPrice}}
			_, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
			if err != nil {
				log.Printf("Failed to update MongoDB: %v", err)
			}
			prices = nil // Reset prices to start aggregating the next 1000 prices
		}
	}
}
