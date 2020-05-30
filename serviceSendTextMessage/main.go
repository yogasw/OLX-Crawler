package main

import (
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
	"github.com/streadway/amqp"
	"log"
	"os"
	"regexp"
	"time"
)

var (
	noPhone string
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func myUsage() {
	fmt.Printf("\nWA Send Message Text Service by arioki1\n\n")
	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
	flag.PrintDefaults()
}

type Message struct {
	Target  string
	Message string
}

func main() {
	server := flag.String("server", "amqp://pcjfkwqt:QIUWu8D3pwgJCibHDLDx5q-QYa1KyHLc@moose.rmq.cloudamqp.com/pcjfkwqt", "server RabbitMQ ex: amqp://xx:xx@xx.com/xx")
	queue := flag.String("queue", "wa-text", "Queue RabbitMQ")
	phone := flag.String("phone", "default", "Primary Number Phone")
	flag.Usage = myUsage
	flag.Parse()
	noPhone = *phone

	//create new WhatsApp connection
	wac, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
		return
	}

	err = login(wac)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
		return
	}

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

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			<-time.After(5 * time.Second)
			var message Message
			json.Unmarshal(d.Body, &message)
			if messageCheck(message.Target) {
				sendMessage(message.Target, message.Message, wac)
			} else {
				log.Printf("Format Target Not Valid, ex : 628xx@s.whatsapp.net or xxx-xxx@g.us")
			}
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func login(wac *whatsapp.Conn) error {
	//load saved session
	session, err := readSession()
	if err == nil {
		//restore session
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring failed: %v\n", err)
		}
	} else {
		//no saved session -> regular login
		qr := make(chan string)
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()
		session, err = wac.Login(qr)
		if err != nil {
			return fmt.Errorf("error during login: %v\n", err)
		}
	}

	//save session
	err = writeSession(session)
	if err != nil {
		return fmt.Errorf("error saving session: %v\n", err)
	}
	return nil
}

func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(os.TempDir() + "/WhatsAppSession/" + noPhone + ".gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func writeSession(session whatsapp.Session) error {
	if _, err := os.Stat(os.TempDir() + "/WhatsAppSession/"); os.IsNotExist(err) {
		os.Mkdir(os.TempDir()+"/WhatsAppSession/", 0700)
	}
	file, err := os.Create(os.TempDir() + "/WhatsAppSession/" + noPhone + ".gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}

func sendMessage(remoteJid string, message string, wac *whatsapp.Conn) {
	previousMessage := "😘"
	quotedMessage := proto.Message{
		Conversation: &previousMessage,
	}

	ContextInfo := whatsapp.ContextInfo{
		QuotedMessage:   &quotedMessage,
		QuotedMessageID: "",
		Participant:     "", //Whot sent the original message
	}

	msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: remoteJid,
		},
		ContextInfo: ContextInfo,
		Text:        message,
	}

	msgId, err := wac.Send(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error sending message: %v", err)
		//os.Exit(1)
		<-time.After(5 * time.Second)
		wac.Send(msg)
	} else {
		fmt.Println("Message Sent -> ID : " + msgId)
	}
}

func messageCheck(str string) bool {
	var re = regexp.MustCompile(`(?m).*@g\.us|.*@s\.whatsapp\.net`)
	match := re.FindAllString(str, -1)
	return len(match) != 0
}
