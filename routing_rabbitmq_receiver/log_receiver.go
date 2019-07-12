package main

import (
	"log"
	"os"

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

	q, err := ch.QueueDeclare(
		"",    //name
		false, //durable
		false, //delete when unused
		true,  //exclusive
		false, //no-wait
		nil,   //argument
	)
	ErrorMsg(err, "Cannot declare queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}

	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, "log_direct", s)
		err = ch.QueueBind(
			q.Name,       //name
			s,            //routing key
			"log_direct", //exchange
			false,        //no-wait
			nil,          //argument
		)
		ErrorMsg(err, "Failed to bind queue")
	}

	msg, err := ch.Consume(
		q.Name, //name
		"",     //cosumer
		true,   //auto-ack
		false,  //exclusive
		false,  //no-local
		false,  //no-wait
		nil,    //argument
	)
	ErrorMsg(err, "Failed to use consumer")

	forever := make(chan bool)

	go func() {
		for d := range msg {
			log.Printf("[x] %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for logs. to exit press CTRL+C")
	<-forever
}

//ErrorMsg is function for show error message all kind of situation
func ErrorMsg(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
