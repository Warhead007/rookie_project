package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	con, err := amqp.Dial("amqp://guest:guest@localhost/")
	ErrorMsg(err, "Cannot connect to RabbitMQ")
	defer con.Close()

	ch, err := con.Channel()
	ErrorMsg(err, "Cannot create channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", //name
		false,   //durable
		false,   //delete when unused
		false,   //exclusive
		false,   //no-wait
		nil,     //arguments
	)
	ErrorMsg(err, "Failed to create queue")

	msg, err := ch.Consume(
		q.Name, //Name
		"",     //Comsumer
		true,   //auto-ack
		false,  //exclusive
		false,  //noLocal
		false,  //noWait
		nil,    //args
	)
	ErrorMsg(err, "Failed to comsume")

	forever := make(chan bool)

	go func() {
		for d := range msg {
			log.Printf("Received a massage: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

//ErrorMsg is function for show error message all kind of situation
func ErrorMsg(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
