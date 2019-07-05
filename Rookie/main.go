package main

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"

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

//ErrorMessage to store error message in JSON//
type ErrorMessage struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

//AllUserData using with GetAllUser function//
type AllUserData struct {
	Count int        `bson:"count" json:"count"`
	Data  []UserData `bson:"data" json:"data"`
}

func main() {

	e := echo.New()
	e.GET("/user/:user_id", GetUser)
	e.GET("/users", GetAllData)

	e.POST("/save", AddData)

	e.PUT("/user/:user_id", UpdateData)

	e.Logger.Fatal(e.Start(":1323"))
}

//AddData :function to get data from user//
func AddData(c echo.Context) error {
	name := c.FormValue("name")
	age := c.FormValue("age")
	email := c.FormValue("email")
	note := c.FormValue("note")
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return err
	}
	///open session to connect database///
	session, err := mgo.Dial(server)
	if err != nil {
		return err
	}
	defer session.Close()

	///access to database and collection to using data///
	a := session.DB(database).C(collection)
	//check email format//
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	///email invalid error message in JSON///
	emailError := &ErrorMessage{
		Code:        "401",
		Description: "Invalid email format",
	}
	///check input data///
	if len(name) == 0 {
		return c.String(http.StatusUnauthorized, "Plase enter your name.")
	} else if len(age) == 0 {
		return c.String(http.StatusUnauthorized, "Plase enter your age.")
	} else if len(email) == 0 {
		return c.String(http.StatusUnauthorized, "Plase enter your email.")
		///if email format invavid///
	} else if !re.MatchString(email) {
		return c.JSON(http.StatusUnauthorized, emailError)
	}
	//souce of image//
	src, err := avatar.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	///set path of image in server///
	fileName := "img/" + avatar.Filename

	//destination to upload image//
	dst, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer dst.Close()

	//copy image from souce to destination//
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	//for get file name to check file type//
	o, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer o.Close()

	//using getFileType//
	contentType, err := GetFileType(o)
	if err != nil {
		return c.HTML(http.StatusBadRequest, "")
	}
	///file type error message in JSON///
	fileError := &ErrorMessage{
		Code:        "401",
		Description: "Invalid file type. Upload .png or .jpg/.jpeg only",
	}
	///check file type///
	if contentType != "image/png" && contentType != "image/jpeg" && contentType != "image/jpg" {
		os.Remove(fileName)
		return c.JSON(http.StatusUnauthorized, fileError)
	}
	///calculate year of birth///
	t := time.Now()
	conAge, err := strconv.Atoi(age)
	if err != nil {
		return err
	}
	///age error message in JSON///
	ageError := &ErrorMessage{
		Code:        "401",
		Description: "Invalid age. You age must in range of 1 - 100",
	}
	///check age validation///
	if conAge <= 0 || conAge > 100 {
		return c.JSON(http.StatusUnauthorized, ageError)
	}
	///calculate year of birth with year now///
	yearOfBirth := t.Year() - conAge
	l, _ := time.LoadLocation("Local")

	///get data from user store into JSON format///
	add := &UserData{
		ID:          bson.NewObjectId(),
		Name:        name,
		Avatarname:  avatar.Filename,
		Avatartype:  contentType,
		Age:         conAge,
		Yearofbirth: yearOfBirth,
		Note:        note,
		Email:       email,
		Createtime:  t.In(l),
		Updatetime:  t.In(l),
	}
	///check error when email exists///
	emailExists := &ErrorMessage{
		Code:        "401",
		Description: "Email is already used. Plase change email",
	}
	///check email with count a found data in database///
	count, err := a.Find(bson.M{"email": &add.Email}).Count()
	if err != nil {
		return err
	}
	///if found count != 0///
	if count > 0 {
		return c.JSON(http.StatusUnauthorized, emailExists)
	}
	///can add data to database///
	a.Insert(add)
	return c.JSON(http.StatusCreated, GetUserData(add.ID))
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

//GetUser : function get user data from GetUserData to show in HTML//
func GetUser(c echo.Context) error {
	id := c.Param("user_id")
	///check if id param send with invaild format (24-digit)///
	if len(id) != 24 {
		return c.HTML(http.StatusUnauthorized, "Invalid ID format. Plase try again")
	}
	///convert string to bson object///
	bsonID := bson.ObjectIdHex(id)

	u := GetUserData(bsonID)
	///when cannot found user with this id///
	findUserError := &ErrorMessage{
		Code:        "401",
		Description: "User not found",
	}
	if u == (UserData{}) {
		return c.JSON(http.StatusUnauthorized, findUserError)
	}

	return c.JSON(http.StatusCreated, u)
}

//GetAllUser : get all user data from database//
func GetAllUser(limit, page int) (*AllUserData, error) {
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
	///variable for store data to show with condition///
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
	if limit == 1 && page <= len(usersData) {
		queryData = append(queryData, usersData[startValue])
	} else if limit == page {
		for i := startValue; i <= limit; i++ {
			///avoid a out of range of slices///
			if i == len(usersData) {
				break
			}
			///query data from userData into queryData///
			queryData = append(queryData, usersData[i])
		}
	} else {
		for i := startValue; i < limit; i++ {
			///avoid a out of range of slices///
			if i == len(usersData) {
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

//GetAllData : get data from GetAllUser to HTML//
func GetAllData(c echo.Context) error {
	///check error, limit value and page value they cannot be 0 or less than///
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		return c.HTML(http.StatusUnauthorized, "Invalid limit value")
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		return c.HTML(http.StatusUnauthorized, "Invalid page value")
	}
	///store value from GetAllUser in u variable///
	u, err := GetAllUser(limit, page)
	if err != nil {
		return c.HTML(http.StatusUnauthorized, "Cannot get user data")
	}
	///return in JSON format in HTML///
	return c.JSON(http.StatusCreated, u)
}

//UpdateData function for get update user data in database//
func UpdateData(c echo.Context) error {
	id := c.Param("user_id")
	///check if id param send with invaild format (24-digit)///
	if len(id) != 24 {
		return c.HTML(http.StatusUnauthorized, "Invalid ID format. Plase try again")
	}
	///convert string to bson object///
	bsonID := bson.ObjectIdHex(id)
	name := c.FormValue("name")
	age := c.FormValue("age")
	note := c.FormValue("note")
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return err
	}
	///open session to connect database///
	session, err := mgo.Dial(server)
	if err != nil {
		return err
	}
	defer session.Close()

	///access to database and collection to using data///
	a := session.DB(database).C(collection)

	//souce of image//
	src, err := avatar.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	///set path of image in server///
	fileName := "img/" + avatar.Filename

	//destination to upload image//
	dst, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer dst.Close()

	//copy image from souce to destination//
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	//for get file name to check file type//
	o, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer o.Close()

	//using getFileType//
	contentType, err := GetFileType(o)
	if err != nil {
		return c.HTML(http.StatusBadRequest, "")
	}
	///file type error message in JSON///
	fileError := &ErrorMessage{
		Code:        "401",
		Description: "Invalid file type. Upload .png or .jpg/.jpeg only",
	}
	///check file type///
	if contentType != "image/png" && contentType != "image/jpeg" && contentType != "image/jpg" {
		os.Remove(fileName)
		return c.JSON(http.StatusUnauthorized, fileError)
	}
	///calculate year of birth///
	t := time.Now()
	conAge, err := strconv.Atoi(age)
	if err != nil {
		return err
	}
	///age error message in JSON///
	ageError := &ErrorMessage{
		Code:        "401",
		Description: "Invalid age. You age must in range of 1 - 100",
	}
	///check age validation///
	if conAge <= 0 || conAge > 100 {
		return c.JSON(http.StatusUnauthorized, ageError)
	}
	///calculate year of birth with year now///
	yearOfBirth := t.Year() - conAge
	l, _ := time.LoadLocation("Local")

	if note == "clean" {
		note = ""
	}

	// update := &UserData{
	// 	Name:        name,
	// 	Avatarname:  avatar.Filename,
	// 	Avatartype:  contentType,
	// 	Age:         conAge,
	// 	Yearofbirth: yearOfBirth,
	// 	Note:        note,
	// 	Updatetime:  t.In(l),
	// }

	a.UpdateId(bsonID, bson.M{"$set": bson.M{
		"name":          name,
		"avatar_name":   avatar.Filename,
		"avatar_type":   contentType,
		"age":           conAge,
		"year_of_birth": yearOfBirth,
		"note":          note,
		"update_time":   t.In(l)}})
	return c.JSON(http.StatusCreated, GetUserData(bsonID))
}
