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

//AddData : add new user data into database//
func AddData(name string, avatarName string, avatarType string, age int, yearOfBirth int, note string, email string) bson.ObjectId {
	///open session to connect database///
	session, err := mgo.Dial(server)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	///access to database and collection to using data///
	a := session.DB(database).C(collection)

	t := time.Now()
	l, _ := time.LoadLocation("Local")

	add := &UserData{
		ID:          bson.NewObjectId(),
		Name:        name,
		Avatarname:  avatarName,
		Avatartype:  avatarType,
		Age:         age,
		Yearofbirth: yearOfBirth,
		Note:        note,
		Email:       email,
		Createtime:  t.In(l),
		Updatetime:  t.In(l),
	}
	a.Insert(add)
	return add.ID
}

//GetUserData : function get one user by ID//
func GetUserData(id bson.ObjectId) UserData {
	///open session to connect database///
	session, err := mgo.Dial(server)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	///access to database and collection to using data///
	a := session.DB(database).C(collection)
	userdata := UserData{}
	///query user data with ID///
	a.Find(bson.M{"_id": id}).One(&userdata)
	if err != nil {
		panic(err)
	}
	///return in JSON format///
	return userdata
}

//DeleteUserData : function to delete user by id//
func DeleteUserData(id bson.ObjectId) error {
	///open session to connect database///
	session, err := mgo.Dial(server)
	if err != nil {
		return err
	}
	defer session.Close()
	///access to database and collection to using data///
	a := session.DB(database).C(collection)
	///delete data by user id///
	err = a.Remove(bson.M{"_id": id})
	if err != nil {
		return err
	}
	return err
}

//GetAllUserData : get all user data from database //
func GetAllUserData(limit, page int) (*AllUserData, error) {
	///open session to connect database///
	session, err := mgo.Dial(server)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	///access to database and collection to using data///
	a := session.DB(database).C(collection)
	///variable for store all data of user///
	usersData := []UserData{}
	///variable for store data to show with condition limit and page///
	queryData := []UserData{}
	///query all of user data///
	a.Find(nil).Sort("-create_time").All(&usersData)
	if err != nil {
		panic(err)
	}
	///count all data in database///
	count, err := a.Find(nil).Count()
	///start point to query data from condition///
	startValue := 0
	///check page///
	if page > 1 {
		///start point changed up to page///
		startValue = limit * (page - 1)
	}
	///if limit is 1 and page not higher than len of userData (Avoid index out of range)///
	if limit == 1 && page <= len(usersData) {
		queryData = append(queryData, usersData[startValue])
	} else {
		for i := startValue; i < startValue+limit; i++ {
			///avoid a out of range of slices///
			if i >= len(usersData) {
				break
			}
			///query data from userData into queryData///
			queryData = append(queryData, usersData[i])
		}
	}
	///store all data to show in show variable///
	show := &AllUserData{
		Count: count,
		Data:  queryData,
	}
	///return in JSON format///
	return show, err
}

//UpdateData : update user data by user id//
func UpdateData(id bson.ObjectId, name string, age int, yearOfBirth int, avatarName string, note string, avatarType string) bson.ObjectId {
	///open session to connect database///
	session, err := mgo.Dial(server)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	///access to database and collection to using data///
	a := session.DB(database).C(collection)
	///calculate year of birth///
	t := time.Now()
	l, _ := time.LoadLocation("Local")

	///if user change data///
	if name != "" {
		a.UpdateId(id, bson.M{"$set": bson.M{
			"name":        name,
			"update_time": t.In(l)}})
	}
	if note != "" {
		a.UpdateId(id, bson.M{"$set": bson.M{
			"note":        note,
			"update_time": t.In(l)}})
	} else if note == "clean" {
		///if user input in note "clean". note field will be delete///
		a.UpdateId(id, bson.M{"$unset": bson.M{"note": ""}})
		a.UpdateId(id, bson.M{"$set": bson.M{"update_time": t.In(l)}})
	}
	///if user send a new avatar file///
	if avatarName != "" && avatarType != "" {

		a.UpdateId(id, bson.M{"$set": bson.M{
			"avatar_name": avatarName,
			"avatar_type": avatarType,
			"update_time": t.In(l)}})
	}
	if age != 0 {
		a.UpdateId(id, bson.M{"$set": bson.M{
			"age":           age,
			"year_of_birth": yearOfBirth,
			"update_time":   t.In(l)}})
	}
	return id
}

//CountEmail function to check email exists in database//
func CountEmail(email string) int {
	///open session to connect database///
	session, err := mgo.Dial(server)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	///access to database and collection to using data///
	a := session.DB(database).C(collection)
	///check email with count a found data in database///
	count, _ := a.Find(bson.M{"email": email}).Count()
	return count
}

//GetFileType : function to get file type
func GetFileType(out *os.File) (string, error) {
	///read file in first 512 byte to check file type///
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "buffer incorrect", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

//CalYearofBirth : function to convert age and calculate year of birth//
func CalYearofBirth(age string) (int, int) {
	///calculate year of birth///
	t := time.Now()
	conAge, _ := strconv.Atoi(age)
	///calculate year of birth with year now///
	yearOfBirth := t.Year() - conAge

	return conAge, yearOfBirth
}
