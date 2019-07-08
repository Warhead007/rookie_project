package main

import (
	"time"

	model "trainer/Rookie/pkg"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

const (
	server     = "localhost:27017"
	database   = "rookie_trainer"
	collection = "user_data"
)

//UserData to handle data//
type UserData struct {
	ID          bson.ObjectId `bson:"_id" json:"_id"`
	Name        string        `bson:"name" json:"name"`
	Avatarname  string        `bson:"avatar_name" json:"avatar_name"`
	Avatartype  string        `bson:"avatar_type" json:"avatar_type"`
	Age         int           `bson:"age" json:"age"`
	Yearofbirth int           `bson:"year_of_birth" json:"year_of_birth"`
	Note        string        `bson:"note,omitempty" json:"note,omitempty"`
	Email       string        `bson:"email" json:"email"`
	Createtime  time.Time     `bson:"create_time" json:"create_time"`
	Updatetime  time.Time     `bson:"update_time" json:"update_time"`
}

//AllUserData using with GetAllUser function//
type AllUserData struct {
	Count int        `bson:"count" json:"count"`
	Data  []UserData `bson:"data" json:"data"`
}

func main() {

	e := echo.New()
	e.GET("/user/:user_id", model.GetUser)
	e.GET("/users", model.GetAllUser)

	e.POST("/save", model.AddUser)

	e.PUT("/user/:user_id", model.UpdateUser)

	e.DELETE("/user/:user_id", model.DeleteUser)

	e.Logger.Fatal(e.Start(":1323"))
}
