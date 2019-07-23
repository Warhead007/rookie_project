package main

import (
	"fmt"
	"net/url"
	functions "trainer/twiiter_api/pkg"

	"github.com/ChimeraCoder/anaconda"
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
	//use search tweet data
	searchResult, _ := api.GetSearch("bnk48 -filter:retweets", v)
	//set layout for readable
	layout := "Mon Jan 02 15:04:05 -0700 2006"
	//print tweet is match with query
	for _, tweet := range searchResult.Statuses {
		//call CalculateTime function
		diff, timeWithFormat := functions.CalculateTime(layout, tweet.CreatedAt)
		//if different time between tweet time and time now less than. Get that data to use
		if diff <= 15 {
			fmt.Println("Time of tweet: ", timeWithFormat)
			fmt.Println(tweet.Text)
			fmt.Println("Text of Tweet: ", tweet.Entities.Media)
			fmt.Println("Different time of tweet and time now: ", int(diff), " minute")
		} else {
			break
		}
	}

}
