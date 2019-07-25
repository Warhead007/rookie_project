package functions

import (
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
func StoreDataForMap(tweet anaconda.Tweet) MongoStreams {
	mongoStream := MongoStreams{
		ChannelTypeID:           "twitter",
		ChannelSouceID:          "twitter",
		ChannelClassificationID: "ta",
		ChannelContentID:        "taw",
		SocialMediaID:           "twitter",
		CreateAt:                time.Now(),
		UpdateAt:                time.Now(),
		Payload:                 tweet,
	}
	if tweet.Entities.Media != nil {
		mongoStream.StreamTypeID = tweet.ExtendedEntities.Media[0].Type
	} else {
		mongoStream.StreamTypeID = "text"
	}
	return mongoStream
}
