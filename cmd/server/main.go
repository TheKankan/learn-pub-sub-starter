package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	connectionString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	newChannel, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}
	defer newChannel.Close()

	fmt.Println("Connected to RabbitMQ successfully.")

	dataToSend, err := json.Marshal(routing.PlayingState{
		IsPaused: true,
	})
	if err != nil {
		log.Fatalf("could not marshal data: %v", err)
	}
	err = pubsub.PublishJSON(newChannel, routing.ExchangePerilDirect, string(routing.PauseKey), dataToSend)
	if err != nil {
		log.Fatalf("could not publish message: %v", err)
	}
	fmt.Println("Published message to RabbitMQ.")

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")
}
