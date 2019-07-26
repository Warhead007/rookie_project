package functions

import (
	"time"

	"github.com/ChimeraCoder/anaconda"
	"gopkg.in/mgo.v2/bson"
)

//FeedData struct for store all feed from collection
type FeedData struct {
	ID       bson.ObjectId `bson:"_id" json:"_id"`
	Keyword  string        `bson:"keyword" json:"keyword"`
	FeedTime time.Time     `bson:"feed_time" json:"feed_time"`
}

//FeedQuery struct for data from worker to map
type FeedQuery struct {
	ChannelTypeID           string         `bson:"channel_type_id" json:"channel_type_id"`
	ChannelSouceID          string         `bson:"channel_souce_id" json:"channel_souce_id"`
	ChannelClassificationID string         `bson:"channel_classification_id" json:"channel_classification_id"`
	ChannelContentID        string         `bson:"channel_content_id" json:"channel_content_id"`
	StreamTypeID            string         `bson:"stream_type_id" json:"stream_type_id"`
	SocialMediaID           string         `bson:"social_media_id" json:"social_media_id"`
	Payload                 anaconda.Tweet `bson:"payload" json:"payload"`
}

//MongoStreams sturct for store data from map function to database
type MongoStreams struct {
	ChannelTypeID           string         `bson:"channel_type_id" json:"channel_type_id"`
	ChannelSouceID          string         `bson:"channel_souce_id" json:"channel_souce_id"`
	ChannelClassificationID string         `bson:"channel_classification_id" json:"channel_classification_id"`
	ChannelContentID        string         `bson:"channel_content_id" json:"channel_content_id"`
	StreamTypeID            string         `bson:"stream_type_id" json:"stream_type_id"`
	SocialMediaID           string         `bson:"social_media_id" json:"social_media_id"`
	CreateAt                time.Time      `bson:"create_at" json:"create_at"`
	UpdateAt                time.Time      `bson:"update_at" json:"update_at"`
	Payload                 anaconda.Tweet `bson:"payload" json:"payload"`
}
