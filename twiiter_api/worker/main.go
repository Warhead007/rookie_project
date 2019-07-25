package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	functions "trainer/twiiter_api/pkg"

	"github.com/ChimeraCoder/anaconda"
	"github.com/streadway/amqp"
)

const (
	exchangeNameFromMaster = "mastertoworker"
	exchangeNameToMap      = "workertomap"
	queueNameFromMaster    = "mastertoworker"
	queueNameToMap         = "workertomap"
)

func main() {
	api := anaconda.NewTwitterApiWithCredentials(
		"",
		"",
		"",
		"")

	//parameter value using with search
	v := url.Values{}
	//set search 100 new tweet
	v.Set("count", "100")
	v.Set("result_type", "recent")
	//set layout for readable
	layout := "Mon Jan 02 15:04:05 -0700 2006"
	//connect with rabbitmq
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	functions.FailOnError(err, "Failed to connect to RabbitMQ.")
	defer conn.Close()
	//create channel to rabbitmq
	cha, err := conn.Channel()
	functions.FailOnError(err, "Failed to open a channel.")
	defer cha.Close()
	//declare exchange to connect master
	err = functions.DeclareExchange(cha, exchangeNameFromMaster)
	//declare exchange to connect map
	err = functions.DeclareExchange(cha, exchangeNameToMap)
	functions.FailOnError(err, "Cannot declare exchange.")
	//declare queue from master
	queue, err := functions.DeclareQueue(cha, queueNameFromMaster)
	functions.FailOnError(err, "Cannot declare queue.")
	//bind queue to connect queue with exchange from master
	err = functions.BindQueue(cha, exchangeNameFromMaster, queueNameFromMaster)
	functions.FailOnError(err, "Cannot binding queue.")
	//consume data from queue
	msgs, err := functions.ConsumeData(cha, queue.Name)
	functions.FailOnError(err, "Consume failed.")
	fmt.Println("Worker starting.")
	go func() {
		for d := range msgs {
			//convert data from master for useable
			var feedData functions.FeedData
			json.Unmarshal(d.Body, &feedData)
			//use search tweet data
			searchResult, _ := api.GetSearch(feedData.Keyword+" -filter:retweets", v)
			//print tweet is match with query
			for i, tweet := range searchResult.Statuses {
				diff, timeWithFormat := functions.CalculateTime(layout, tweet.CreatedAt)
				//if different time between tweet time and time now less than. Get that data to use
				if diff <= 15 {
					fmt.Println("Time of tweet: ", timeWithFormat)
					fmt.Println("Different time of tweet and time now: ", int(diff), " minute")
					fmt.Println("Tweet count:", i+1)
					fmt.Println(functions.StoreDataForMap(tweet))
					//convert data to send to map function
					conTweetData, err := json.Marshal(functions.StoreDataForMap(tweet))
					functions.FailOnError(err, "Cannot convert this struct to JSON.")
					err = functions.PublishData(cha, exchangeNameToMap, conTweetData)
				} else {
					break
				}
			}
			fmt.Println("--------------------------------------------------------------------")
		}
	}()

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	fmt.Println("Worker stoped.")
}
