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

//UserData struct//
type UserData struct {
	ID          bson.ObjectId `bson:"_id"`
	Name        string        `bson:"name" json:"name"`
	Avatarname  string        `bson:"avatar_name" json:"avatar_name"`
	Avatartype  string        `bson:"avatar_type" json:"avatar_type"`
	Age         int           `bson:"age" json:"age"`
	Yearofbirth int           `bson:"year_of_birth" json:"year_of_birth"`
	Note        string        `bson:"note,omitempty" json:"note,omitempty"`
	Email       string        `bson:"email" json:"email"`
	Createtime  string        `bson:"create_time" json:"create_time"`
	Updatetime  string        `bson:"update_time" json:"update_time"`
}

//ErrorMessage struct//
type ErrorMessage struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func main() {

	e := echo.New()
	e.POST("/user", createUser)
	e.POST("/save", addData)

	e.Logger.Fatal(e.Start(":1323"))
}

func createUser(c echo.Context) error {

	return c.JSON(http.StatusCreated, addData(c))
}

///function to get data from user///
func addData(c echo.Context) error {
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
	contentType, err := getFileType(o)
	if err != nil {
		return c.HTML(http.StatusBadRequest, "")
	}
	///file type error message in JSON///
	fileError := &ErrorMessage{
		Code:        "401",
		Description: "Invalid file type. Upload .png or .jpg/.jpeg only",
	}
	///check file type///
	if contentType != "image/png" && contentType != "image/jpg" {
		os.Remove(fileName)
		return c.JSON(http.StatusUnauthorized, fileError)
	}
	///calculate year of birth///
	time := time.Now()
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
	yearOfBirth := time.Year() - conAge
	createTime := time.Format("15:04:05 02-01-2006")
	updateTime := createTime

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
		Createtime:  createTime,
		Updatetime:  updateTime,
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
	return c.JSON(http.StatusCreated, getUserData(add.ID))
}

///function get one user by ID///
func getUserData(id bson.ObjectId) UserData {
	///open session to connect database///
	session, err := mgo.Dial(server)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	///access to database and collection to using data///
	a := session.DB(database).C(collection)
	userdata := UserData{}
	///find ID of user///
	a.Find(bson.M{"_id": id}).One(&userdata)
	if err != nil {
		panic(err)
	}
	///return in JSON format///
	return userdata
}

///function to get file type///
func getFileType(out *os.File) (string, error) {
	///read file in first 512 byte to check file type///
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "buffer incorrect", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
