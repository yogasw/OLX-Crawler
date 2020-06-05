package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"regexp"
	"time"
)

var (
	noPhone string
)

func myUsage() {
	fmt.Printf("\nWA Send Message Text Service by arioki1\n\n")
	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
	flag.PrintDefaults()
}

type Message struct {
	Target  string
	Message string
	Image   string
}

func main() {
	server := flag.String("server", "amqp://"+os.Getenv("RABBITMQ_DEFAULT_USER")+":"+os.Getenv("RABBITMQ_DEFAULT_PASS")+"@"+os.Getenv("RABBITMQ_SERVER")+os.Getenv("RABBITMQ_DEFAULT_VHOST"), "server RabbitMQ ex: amqp://xx:xx@xx.com/xx")
	queue := flag.String("queue", os.Getenv("RABBITMQ_DEFAULT_QUEUE"), "Queue RabbitMQ")
	phone := flag.String("phone", os.Getenv("MASTER_PHONE_NUMBER"), "Primary Number Phone")
	flag.Usage = myUsage
	flag.Parse()
	noPhone = *phone
	println(*server)
	//create new RabbitMQ connection
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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	wac := connectionWhatsApp()
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			<-time.After(5 * time.Second)
			var message Message
			json.Unmarshal(d.Body, &message)
			if messageCheck(message.Target) {
				if wac == nil {
					wac = connectionWhatsApp()
					sendMessage(message.Target, message.Message, message.Image, wac)
				} else {
					sendMessage(message.Target, message.Message, message.Image, wac)
				}

			} else {
				log.Printf("Format Target Not Valid, ex : 628xx@s.whatsapp.net or xxx-xxx@g.us")
			}
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func messageCheck(str string) bool {
	var re = regexp.MustCompile(`(?m).*@g\.us|.*@s\.whatsapp\.net`)
	match := re.FindAllString(str, -1)
	return len(match) != 0
}
