package main

import (
	"encoding/json"
	"flag"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type Message struct {
	Target  string
	Message string
}

func main() {
	message := flag.String("message", "default", "message")
	target := flag.String("target", "6282329949292-1590306644@g.us", "ex : 628xx@s.whatsapp.net or xxx-xxx@g.us")
	server := flag.String("server", "amqp://pcjfkwqt:QIUWu8D3pwgJCibHDLDx5q-QYa1KyHLc@moose.rmq.cloudamqp.com/pcjfkwqt", "RabbitMQ Server")
	queue := flag.String("queue", "wa-text", "Queue RabbitMQ")
	flag.Parse()

	conn, err := amqp.Dial(*server)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		*queue, // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")
	data := Message{
		Target:  *target,
		Message: *message,
	}
	body, err := json.Marshal(data)
	failOnError(err, "Error")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}
