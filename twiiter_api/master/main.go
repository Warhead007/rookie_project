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

	functions "trainer/twiiter_api/pkg"
)

const (
	exchangeName = "ha_twfeed"
	rountingKey  = "ha_twfeed.tweet.add"
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
	//declare exchange to send data to queue
	err = functions.DeclareExchange(cha, exchangeName)
	functions.FailOnError(err, "Cannot declare exchange.")
	//set layout for use in function CalculateTime
	layout := "Mon Jan 02 15:04:05 -0700 2006"
	//set ticker run every 1 minute
	ticker := time.NewTicker(1 * time.Minute)
	fmt.Println("master starting.")
	go func() {
		for now := range ticker.C {
			feedData := functions.GetAllFeed()
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
					//call function update time to current time
					functions.UpdateFeedTime(feedData.ID)
					fmt.Println("Update feed time:", feedData.Keyword)
					conFeedData, err := json.Marshal(feedData)
					functions.FailOnError(err, "Cannot convert this struct to JSON.")
					//set publisher
					err = functions.PublishData(cha, exchangeName, rountingKey, conFeedData)
					fmt.Println("Send", string(feedData.Keyword), "to worker")
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
	fmt.Println("master stopped.")
	ticker.Stop()
}
