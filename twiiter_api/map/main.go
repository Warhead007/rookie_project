package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	functions "trainer/twiiter_api/pkg"

	"github.com/streadway/amqp"
)

const (
	exchangeNameFromWorker = "workertomap"
	queueNameFromWorker    = "workertomap"
)

func main() {
	//connect with rabbitmq
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	functions.FailOnError(err, "Failed to connect to RabbitMQ.")
	defer conn.Close()
	//create channel to rabbitmq
	cha, err := conn.Channel()
	functions.FailOnError(err, "Failed to open a channel.")
	defer cha.Close()
	//declare queue from master
	queue, err := functions.DeclareQueue(cha, queueNameFromWorker)
	functions.FailOnError(err, "Cannot declare queue.")
	//bind queue to connect queue with exchange from master
	err = functions.BindQueue(cha, exchangeNameFromWorker, queueNameFromWorker)
	functions.FailOnError(err, "Cannot binding queue.")
	//consume data from queue
	msgs, err := functions.ConsumeData(cha, queue.Name)
	functions.FailOnError(err, "Consume failed.")
	fmt.Println("Map starting.")
	go func() {
		for d := range msgs {
			var mongoStreams functions.MongoStreams
			json.Unmarshal(d.Body, &mongoStreams)
			fmt.Println(mongoStreams)
		}
	}()
	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	fmt.Println("Map stopped.")
}
