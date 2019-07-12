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

	err = ch.ExchangeDeclare(
		"log_direct", //name
		"direct",     //type
		true,         //durable
		false,        //delete when unused
		false,        //internal
		false,        //no-wait
		nil,          //argments
	)
	ErrorMsg(err, "Cannot declare exchange")

	body := bodyFrom(os.Args)
	err = ch.Publish(
		"log_direct",          //exchange
		severityFrom(os.Args), //routing key
		false,                 //mandatory
		false,                 //immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
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

//severityForm function to get string from cmd to using in routing key (lesson 4)
func severityFrom(args []string) string {
	s := ""
	if (len(args) < 2) || os.Args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}
	return s
}
