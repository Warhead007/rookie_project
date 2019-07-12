package main

import (
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

func main() {
	//get connection with RabbitMQ server
	con, err := amqp.Dial("amqp://guest:guest@localhost/")
	ErrorMsg(err, "Cannot connect to RabbitMQ")
	defer con.Close()

	//open channel to send data
	ch, err := con.Channel()
	ErrorMsg(err, "Cannot create channel")
	defer ch.Close()

	//create queue variable
	//lesson 1
	// q, err := ch.QueueDeclare(
	// 	"hello", //name
	// 	false,   //durable
	// 	false,   //delete when unused
	// 	false,   //exclusive
	// 	false,   //no-wait
	// 	nil,     //arguments
	// )
	//lesson 2
	q, err := ch.QueueDeclare(
		"hello_1", //name
		true,      //durable
		false,     //delete when unused
		false,     //exclusive
		false,     //no-wait
		nil,       //arguments
	)
	ErrorMsg(err, "Failed to create queue")

	//set data to send and publish in RabbitMQ server
	// body := "Hello world"

	//use message from user (lesson 2)
	body := bodyFrom(os.Args)
	err = ch.Publish(
		"",     //exchange
		q.Name, //routing key
		false,  // mandatory
		false,  // immediate
		//send data in here
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	ErrorMsg(err, "Cannot publish a message")
	//print in commandline
	log.Printf(" [x] Sent %s", body)
}

//ErrorMsg is function for show error message all kind of situation
func ErrorMsg(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

//bodyFrom function to get string from cmd to send (lesson 2)
func bodyFrom(args []string) string {
	s := ""
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
