package main

import (
	"log"

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
	q, err := ch.QueueDeclare(
		"hello", //name
		false,   //durable
		false,   //delete when unused
		false,   //exclusive
		false,   //no-wait
		nil,     //arguments
	)
	ErrorMsg(err, "Failed to create queue")

	//set data to send and publish in RabbitMQ server
	body := "Hello world"
	err = ch.Publish(
		"",     //exchange
		q.Name, //routing key
		false,  // mandatory
		false,  // immediate
		//send data in here
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	ErrorMsg(err, "Cannot publish a message")

}

//ErrorMsg is function for show error message all kind of situation
func ErrorMsg(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
