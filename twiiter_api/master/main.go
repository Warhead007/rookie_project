package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/mgo.v2/bson"

	functions "trainer/twiiter_api/pkg"

	"gopkg.in/mgo.v2"
)

const (
	server     = "localhost:27017"
	database   = "feed"
	collection = "feed_keyword"
)

//FeedData struct for store all feed from collection
type FeedData struct {
	ID       bson.ObjectId `bson:"_id" json:"_id"`
	Keyword  string        `bson:"keyword" json:"keyword"`
	FeedTime time.Time     `bson:"feed_time" json:"feed_time"`
}

func main() {
	//open session to connect database
	session, err := mgo.Dial(server)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	//access to database and collection to using data
	a := session.DB(database).C(collection)
	feedData := []FeedData{}
	//set layout for use in function CalculateTime
	layout := "Mon Jan 02 15:04:05 -0700 2006"
	//set ticker run every 1 minute
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for now := range ticker.C {
			//get all data from collection Feed
			a.Find(nil).All(&feedData)
			if err != nil {
				panic(err)
			}
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
				}
			}
			fmt.Println("Current time:", timeNow)
		}
	}()
	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	ticker.Stop()

}
