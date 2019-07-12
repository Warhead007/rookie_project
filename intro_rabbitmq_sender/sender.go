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

	body := "Hello world"
	err = ch.Publish(
		"",     //exchange
		q.Name, //routing key
		false,  // mandatory
		false,  // immediate
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
