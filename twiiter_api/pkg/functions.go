package functions

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

//CalculateTime function to get time from tweet then convert and calculate time
//between feed_time and time now
func CalculateTime(layout string, timeValue string) (int, string) {
	t, _ := time.Parse(layout, timeValue)
	l, _ := time.LoadLocation("Local")
	timeWithFormat := t.In(l).Format("15:04")
	diff := time.Now().Sub(t).Minutes()
	return int(diff), timeWithFormat
}

//FailOnError function to handle with error by show message
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

//StoreDataForMap function to make data and store tweet data in payload to send to Map function
func StoreDataForMap(tweet anaconda.Tweet) FeedQuery {
	feedQuery := FeedQuery{
		ChannelTypeID:           "twitter",
		ChannelSouceID:          "twitter",
		ChannelClassificationID: "ta",
		ChannelContentID:        "taw",
		SocialMediaID:           tweet.User.IdStr,
		Payload:                 tweet,
	}
	return feedQuery
}

//StoreDataForStream function to make data from map ,combine data with payload and send into database
func StoreDataForStream(tweet anaconda.Tweet) []byte {
	checkTypeStream := ""
	if tweet.ExtendedEntities.Media == nil {
		checkTypeStream = "text"
	} else {
		checkTypeStream = tweet.ExtendedEntities.Media[0].Type
	}
	//payload only
	var payload anaconda.Tweet
	payload = tweet
	//set mongostream data
	mongoStreams := MongoStreams{
		ChannelTypeID:           "twitter",
		ChannelSouceID:          "twitter",
		ChannelClassificationID: "ta",
		ChannelContentID:        "taw",
		SocialMediaID:           tweet.User.IdStr,
		StreamTypeID:            checkTypeStream,
		CreateAt:                time.Now(),
		UpdateAt:                time.Now(),
	}
	//combine mongoStream and payload to JSON
	conPayloadData, _ := json.Marshal(struct {
		MongoStreams
		anaconda.Tweet
	}{mongoStreams, payload})

	return conPayloadData
}
