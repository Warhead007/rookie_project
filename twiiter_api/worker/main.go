package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
	functions "trainer/twiiter_api/pkg"

	"github.com/ChimeraCoder/anaconda"
	"github.com/streadway/amqp"
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

	queue, err := cha.QueueDeclare(
		"mastertoworker", //name
		false,            //durable
		false,            //delete when not used
		false,            //exclusive
		false,            //no-wait
		nil,              //args
	)
	functions.FailOnError(err, "Cannot declare queue.")

	err = cha.QueueBind(
		queue.Name, //name
		"",         //rounting key
		"keyword",  //exchange name
		false,      //no-wait
		nil,        //args
	)
	functions.FailOnError(err, "Cannot binding queue.")

	msgs, err := cha.Consume(
		queue.Name, //name
		"",         //consumer
		true,       //auto-ack
		false,      //exclusive
		false,      //no-local
		false,      //no-wait
		nil,        //args
	)
	functions.FailOnError(err, "Consume failed.")
	fmt.Println("Worker starting.")
	go func() {
		for d := range msgs {
			var mongoStream functions.MongoStreams
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
					mongoStream.ChannelTypeID = "Twitter"
					mongoStream.ChannelSouceID = "Twitter"
					mongoStream.ChannelClassificationID = "Twitter"
					mongoStream.ChannelContentID = "twitter"
					mongoStream.SocialMediaID = "twitter"
					mongoStream.CreateAt = time.Now()
					mongoStream.UpdateAt = time.Now()
					fmt.Println("Time of tweet: ", timeWithFormat)
					//fmt.Println("Text of Tweet: ", tweet.Text)
					//check type of tweet
					if tweet.Entities.Media != nil {
						//fmt.Println(tweet.ExtendedEntities.Media[0].Type)
						mongoStream.StreamTypeID = tweet.ExtendedEntities.Media[0].Type
					} else {
						fmt.Println("text")
					}
					mongoStream.Payload = tweet
					fmt.Println("Different time of tweet and time now: ", int(diff), " minute")
					fmt.Println("Tweet count:", i+1)
					fmt.Println(mongoStream)
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
