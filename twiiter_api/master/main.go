package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2/bson"

	functions "trainer/twiiter_api/pkg"

	"gopkg.in/mgo.v2"
)

const (
	server     = "localhost:27017"
	database   = "feed"
	collection = "feed_keyword"
)

func main() {
	//open session to connect database
	session, err := mgo.Dial(server)
	functions.FailOnError(err, "Cannot connect mongoDB.")
	defer session.Close()
	//access to database and collection to using data
	a := session.DB(database).C(collection)
	//connect with rabbitmq
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	functions.FailOnError(err, "Failed to connect to RabbitMQ.")
	defer conn.Close()
	//create channel to rabbitmq
	cha, err := conn.Channel()
	functions.FailOnError(err, "Failed to open a channel.")
	defer cha.Close()
	//declare exchange to send data to queue
	err = cha.ExchangeDeclare(
		"keyword", // name
		"fanout",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // args
	)
	functions.FailOnError(err, "Cannot declare exchange.")
	//set variable slices of sturct
	feedData := []functions.FeedData{}
	//set layout for use in function CalculateTime
	layout := "Mon Jan 02 15:04:05 -0700 2006"
	//set ticker run every 1 minute
	ticker := time.NewTicker(1 * time.Minute)
	fmt.Println("master starting.")
	go func() {
		for now := range ticker.C {
			//get all data from collection Feed
			a.Find(nil).All(&feedData)
			//set variable for show current time in string
			timeNow := ""
			//loop for get data
			for _, feedData := range feedData {
				//call CalculateTime function
				diff, timeWithFormat := functions.CalculateTime(layout, feedData.FeedTime.Format(layout))
				//set now format for readable
				timeNow = now.Format("15:04")
				//print data and different time between feed_time and time now
				fmt.Println(feedData.Keyword, timeWithFormat, diff)
				//if different time higher than 15 minute
				if diff >= 15 {
					//set feed_time with time now
					a.UpdateId(feedData.ID, bson.M{"$set": bson.M{
						"feed_time": time.Now()}})
					fmt.Println("Update feed time:", feedData.Keyword)
					conFeedData, err := json.Marshal(feedData)
					functions.FailOnError(err, "Cannot convert this struct to JSON.")
					//set publisher
					err = cha.Publish(
						"keyword", //name
						"",        //rounting key
						false,     //mandatory
						false,     //immediate
						amqp.Publishing{
							ContentType: "text/plain",
							Body:        []byte(conFeedData),
						})
					fmt.Println("Send", string(conFeedData), "to worker")
				}
			}
			fmt.Println("Current time:", timeNow)
			fmt.Println("--------------------------------------------------------------------")
		}
	}()
	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	fmt.Println("master stoped.")
	ticker.Stop()
}
