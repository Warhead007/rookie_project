package main

import (
	"bytes"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	//get connect with RabbitMQ server
	con, err := amqp.Dial("amqp://guest:guest@localhost/")
	ErrorMsg(err, "Cannot connect to RabbitMQ")
	defer con.Close()

	//create channel to receive data
	ch, err := con.Channel()
	ErrorMsg(err, "Cannot create channel")
	defer ch.Close()

	//create queue variable (must same at sender variable)
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
	//lesson 3
	// q, err := ch.QueueDeclare(
	// 	"",    //name
	// 	false, //durable
	// 	false, //delete when unused
	// 	true,  //exclusive
	// 	false, //no-wait
	// 	nil,   //arguments
	// )
	// ErrorMsg(err, "Failed to create queue")
	//lesson 3
	// err = ch.ExchangeDeclare(
	// 	"logs",   //name
	// 	"fanout", //type
	// 	true,     //durable
	// 	false,    //auto-deleted
	// 	false,    //internal
	// 	false,    //no-wait
	// 	nil,      //arguments
	// )
	// ErrorMsg(err, "Failed to declare exchange")
	//lesson 3
	// err = ch.QueueBind(
	// 	q.Name, //name
	// 	"",     //routing key
	// 	"logs", //exchange
	// 	false,  //no-wait
	// 	nil,    //arguments
	// )
	// ErrorMsg(err, "Failed to use queueBind")

	//consume data to use
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

	//set prefetch
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	ErrorMsg(err, "Cannot set Qos")
	//set to exit receiving
	forever := make(chan bool)

	//get data here
	go func() {
		for d := range msg {
			log.Printf("Received a massage: %s", d.Body)
			d.Ack(false)
			//for delay message simulates realwork (Lesson 2)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
		}
	}()

	//close channel
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

//ErrorMsg is function for show error message all kind of situation
func ErrorMsg(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
