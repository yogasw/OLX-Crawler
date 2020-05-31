package main

import (
	"encoding/json"
	"flag"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type Message struct {
	Target  string
	Message string
	Image   string
}

func main() {
	message := flag.String("message", "default", "message")
	image := flag.String("image", "https://api.cloudamqp.com/img/lemur_256.png", "image")
	server := flag.String("server", "amqp://"+os.Getenv("RABBITMQ_DEFAULT_USER")+":"+os.Getenv("RABBITMQ_DEFAULT_PASS")+"@rabbitmq"+os.Getenv("RABBITMQ_DEFAULT_VHOST"), "server RabbitMQ ex: amqp://xx:xx@xx.com/xx")
	queue := flag.String("queue", os.Getenv("RABBITMQ_DEFAULT_QUEUE"), "Queue RabbitMQ")
	target := flag.String("target", os.Getenv("TARGET_WA_MESSAGE"), "ex : 628xx@s.whatsapp.net or xxx-xxx@g.us")

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
		Image: *image,
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
