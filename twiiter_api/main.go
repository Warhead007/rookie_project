package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

//UserInformation struct for store data of user tweet at the time
type UserInformation struct {
	UserName     string    `json:"username"`
	TweetMessage []string  `json:"tweet_message"`
	TweetCount   int       `json:"count_message"`
	StartRecord  time.Time `json:"start_record"`
	EndRecord    time.Time `json:"end_record"`
}

//UserList for collect and create struct of UserInformation
type UserList struct {
	UserInformation []UserInformation
}

func main() {
	config := oauth1.NewConfig("", "")
	token := oauth1.NewToken("", "")
	// OAuth1 setting
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)
	//set user to tracking
	userTwitter := []string{"fm91trafficpro", "js100radio", "MorningNewsTV3"}
	//create slices of strcut
	data := []UserInformation{}

	userID := []string{}
	//init twitter id and start time into slice of strcut
	for i := 0; i < len(userTwitter); i++ {
		userName, _, _ := client.Users.Show(&twitter.UserShowParams{
			ScreenName: userTwitter[i],
		})
		add := UserInformation{
			UserName:    userTwitter[i],
			StartRecord: time.Now(),
		}
		//add new struct into slices of strcut
		data = append(data, add)

		userID = append(userID, userName.IDStr)
	}

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		//control pointer to point data in array of strcut
		for i := 0; i < len(userTwitter); i++ {
			//fitter to steaming twitter tweet only
			if tweet.User.IDStr == userID[i] {
				fmt.Println(tweet.User.ScreenName)
				fmt.Println(tweet.Text)
				//add message tweet data into strcut elements
				data[i] = UserInformation{
					UserName:     data[i].UserName,
					TweetMessage: append(data[i].TweetMessage, tweet.Text),
					TweetCount:   data[i].TweetCount + 1,
					StartRecord:  data[i].StartRecord,
				}
			}
		}
	}

	fmt.Println("Starting Stream...")
	fmt.Println(time.Now())

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Follow:        userID,
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}
	//streaming data from twitter
	go demux.HandleChan(stream.Messages)
	//Declare time and set timer to send data
	ticker := time.NewTicker(15 * time.Minute)
	go func() {
		for now := range ticker.C {
			fmt.Println("Send data: ", now)
		}
	}()

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	fmt.Println(time.Now())
	fmt.Println("Summary data : ", data)
	stream.Stop()
	ticker.Stop()

}
