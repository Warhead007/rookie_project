package model

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	server     = "localhost:27017"
	database   = "rookie_trainer"
	collection = "user_data"
)

//UserData to handle data
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

//AllUserData using with GetAllUser function
type AllUserData struct {
	Count int        `bson:"count" json:"count"`
	Data  []UserData `bson:"data" json:"data"`
}

//AddData : add new user data into database
func (u UserData) AddData() bson.ObjectId {
	//open session to connect database
	session, err := mgo.Dial(server)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	//access to database and collection to using data
	a := session.DB(database).C(collection)

	a.Insert(u)
	return u.ID
}

//GetUserData : function get one user by ID
func GetUserData(id bson.ObjectId) UserData {
	//open session to connect database
	session, err := mgo.Dial(server)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	//access to database and collection to using data
	a := session.DB(database).C(collection)
	userdata := UserData{}
	//query user data with ID
	a.Find(bson.M{"_id": id}).One(&userdata)
	if err != nil {
		panic(err)
	}
	//return in JSON format
	return userdata
}

//DeleteUserData : function to delete user by id
func DeleteUserData(id bson.ObjectId) error {
	//open session to connect database
	session, err := mgo.Dial(server)
	if err != nil {
		return err
	}
	defer session.Close()
	//access to database and collection to using data
	a := session.DB(database).C(collection)
	//delete data by user id
	err = a.Remove(bson.M{"_id": id})
	if err != nil {
		return err
	}
	return err
}

//GetAllUserData : get all user data from database
func GetAllUserData(limit, page int) (*AllUserData, error) {
	//open session to connect database
	session, err := mgo.Dial(server)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	//access to database and collection to using data
	a := session.DB(database).C(collection)
	//variable for store all data of user
	usersData := []UserData{}
	//variable for store data to show with condition limit and page
	queryData := []UserData{}
	//query all of user data
	a.Find(nil).Sort("-create_time").All(&usersData)
	if err != nil {
		panic(err)
	}
	//count all data in database
	count, err := a.Find(nil).Count()
	//start point to query data from condition
	startValue := 0
	//check page
	if page > 1 {
		//start point changed up to page
		startValue = limit * (page - 1)
	}
	//if limit is 1 and page not higher than len of userData (Avoid index out of range)
	if limit == 1 && page <= len(usersData) {
		queryData = append(queryData, usersData[startValue])
	} else {
		for i := startValue; i < startValue+limit; i++ {
			//avoid a out of range of slices
			if i >= len(usersData) {
				break
			}
			//query data from userData into queryData
			queryData = append(queryData, usersData[i])
		}
	}
	//store all data to show in show variable
	show := &AllUserData{
		Count: count,
		Data:  queryData,
	}
	//return in JSON format
	return show, err
}

//UpdateData : update user data by user id
func (u UserData) UpdateData() {
	//open session to connect database
	session, err := mgo.Dial(server)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	//access to database and collection to using data
	a := session.DB(database).C(collection)
	//calculate year of birth
	t := time.Now()
	l, _ := time.LoadLocation("Local")

	//if user change data
	if u.Name != "" {
		a.UpdateId(u.ID, bson.M{"$set": bson.M{
			"name":        u.Name,
			"update_time": t.In(l)}})
	}
	if u.Note != "" {
		a.UpdateId(u.ID, bson.M{"$set": bson.M{
			"note":        u.Note,
			"update_time": t.In(l)}})
	}
	if u.Note == "clean" {
		//if user input in note "clean". note field will be delete
		a.UpdateId(u.ID, bson.M{"$unset": bson.M{"note": ""}})
		a.UpdateId(u.ID, bson.M{"$set": bson.M{"update_time": t.In(l)}})
	}
	//if user send a new avatar file
	if u.Avatarname != "" && u.Avatartype != "" {

		a.UpdateId(u.ID, bson.M{"$set": bson.M{
			"avatar_name": u.Avatarname,
			"avatar_type": u.Avatartype,
			"update_time": t.In(l)}})
	}
	if u.Age != 0 {
		a.UpdateId(u.ID, bson.M{"$set": bson.M{
			"age":           u.Age,
			"year_of_birth": u.Yearofbirth,
			"update_time":   t.In(l)}})
	}
}

//CountEmail function to check email exists in database
func CountEmail(email string) int {
	//open session to connect database
	session, err := mgo.Dial(server)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	//access to database and collection to using data
	a := session.DB(database).C(collection)
	//check email with count a found data in database
	count, _ := a.Find(bson.M{"email": email}).Count()
	return count
}

//GetFileType : function to get file type
func GetFileType(out *os.File) (string, error) {
	//read file in first 512 byte to check file type
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "buffer incorrect", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

//CalYearofBirth : function to convert age and calculate year of birth
func CalYearofBirth(age string) (int, int) {
	//calculate year of birth
	t := time.Now()
	conAge, _ := strconv.Atoi(age)
	//calculate year of birth with year now
	yearOfBirth := t.Year() - conAge

	return conAge, yearOfBirth
}
