package functions

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//FeedData struct for store all feed from collection
type FeedData struct {
	ID       bson.ObjectId `bson:"_id" json:"_id"`
	Keyword  string        `bson:"keyword" json:"keyword"`
	FeedTime time.Time     `bson:"feed_time" json:"feed_time"`
}
