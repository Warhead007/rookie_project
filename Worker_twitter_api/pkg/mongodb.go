package functions

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	server     = "localhost:27017"
	database   = "feed"
	collection = "feed_keyword"
)

//GetAllFeed function to take all feed data to
func GetAllFeed() []FeedData {
	//open session to connect database
	session, err := mgo.Dial(server)
	FailOnError(err, "Cannot connect mongoDB.")
	defer session.Close()
	//access to database and collection to using data
	a := session.DB(database).C(collection)
	feedData := []FeedData{}
	a.Find(nil).All(&feedData)
	return feedData
}

//UpdateFeedTime function for update feed time to keywork
func UpdateFeedTime(id bson.ObjectId) {
	//open session to connect database
	session, err := mgo.Dial(server)
	FailOnError(err, "Cannot connect mongoDB.")
	defer session.Close()
	//access to database and collection to using data
	a := session.DB(database).C(collection)
	//set feed_time with time now
	a.UpdateId(id, bson.M{"$set": bson.M{
		"feed_time": time.Now()}})
}

//AddDataStream function to get data from stream into database
func AddDataStream(data interface{}) {
	//open session to connect database
	session, err := mgo.Dial(server)
	FailOnError(err, "Cannot connect mongoDB.")
	defer session.Close()
	//access to database and collection to using data
	a := session.DB("mongo_streams").C("twitter_data")
	a.Insert(&data)
}
