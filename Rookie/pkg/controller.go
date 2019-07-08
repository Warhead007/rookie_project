package model

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

//ErrorMessage to store error message in JSON//
type ErrorMessage struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

var internalError = &ErrorMessage{
	Code:        "500",
	Description: "Internal Error. Plase try again.",
}

//AddUser :function to get data from user//
func AddUser(c echo.Context) error {
	name := c.FormValue("name")
	age := c.FormValue("age")
	email := c.FormValue("email")
	note := c.FormValue("note")
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return c.String(http.StatusUnauthorized, "Plase upload your profile picture.")
	}
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
		return c.JSON(http.StatusInternalServerError, internalError)
	}
	defer src.Close()
	///set path of image in server///
	fileName := "img/" + avatar.Filename

	//destination to upload image//
	dst, err := os.Create(fileName)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, internalError)
	}
	defer dst.Close()

	//copy image from souce to destination//
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	//for get file name to check file type//
	o, err := os.Open(fileName)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, internalError)
	}
	defer o.Close()

	//using getFileType//
	contentType, err := GetFileType(o)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, internalError)
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

	conAge, yearOfBirth := CalYearofBirth(age)
	///age error message in JSON///
	ageError := &ErrorMessage{
		Code:        "401",
		Description: "Invalid age. You age must in range of 1 - 100",
	}
	///check age validation///
	if conAge <= 0 || conAge > 100 {
		return c.JSON(http.StatusUnauthorized, ageError)
	}
	///add data into AddData function///
	///check error when email exists///
	emailExists := &ErrorMessage{
		Code:        "401",
		Description: "Email is already used. Plase change email",
	}
	///if found count != 0///
	if CountEmail(email) > 0 {
		return c.JSON(http.StatusUnauthorized, emailExists)
	}
	///can add data to database///
	id := AddData(name, avatar.Filename, contentType, conAge, yearOfBirth, note, email)
	return c.JSON(http.StatusCreated, GetUserData(id))
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
		Code:        "404",
		Description: "User not found",
	}
	///if not found any data///
	if u == (UserData{}) {
		return c.JSON(http.StatusNotFound, findUserError)
	}

	return c.JSON(http.StatusCreated, u)
}

//GetAllUser : get data from GetAllUserData to HTML//
func GetAllUser(c echo.Context) error {
	///check error, limit value and page value they cannot be 0 or less than///
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	///store value from GetAllUser in u variable///
	u, err := GetAllUserData(limit, page)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, internalError)
	}
	///return in JSON format in HTML///
	return c.JSON(http.StatusCreated, u)
}

//UpdateUser : function for get update user data in database//
func UpdateUser(c echo.Context) error {
	id := c.Param("user_id")
	///check if id param send with invaild format (24-digit)///
	if len(id) != 24 {
		return c.HTML(http.StatusUnauthorized, "Invalid ID format. Plase try again")
	}
	///convert string to bson object///
	bsonID := bson.ObjectIdHex(id)

	///when not found user with this id///
	findUserError := &ErrorMessage{
		Code:        "404",
		Description: "User not found",
	}
	///if not found any data in database///
	if GetUserData(bsonID) == (UserData{}) {
		return c.JSON(http.StatusNotFound, findUserError)
	}

	name := c.FormValue("name")
	age := c.FormValue("age")
	note := c.FormValue("note")
	avatar, _ := c.FormFile("avatar")
	contentType := ""
	avatarName := ""

	///if user send a new avatar file///
	if avatar != nil {
		//souce of image//
		src, err := avatar.Open()
		if err != nil {
			return c.JSON(http.StatusUnauthorized, internalError)
		}
		defer src.Close()
		///set path of image in server///
		fileName := "img/" + avatar.Filename

		//destination to upload image//
		dst, err := os.Create(fileName)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, internalError)
		}
		defer dst.Close()

		//copy image from souce to destination//
		if _, err = io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusUnauthorized, internalError)
		}
		//for get file name to check file type//
		o, err := os.Open(fileName)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, internalError)
		}
		defer o.Close()

		//using getFileType//
		contentType, err := GetFileType(o)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, internalError)
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
		///set file name///
		avatarName = avatar.Filename
	}
	conAge, yearOfBirth := 0, 0
	if age != "" {
		conAge, yearOfBirth = CalYearofBirth(age)
		///age error message in JSON///
		ageError := &ErrorMessage{
			Code:        "401",
			Description: "Invalid age. You age must in range of 1 - 100",
		}
		///check age validation///
		if conAge <= 0 || conAge > 100 {
			return c.JSON(http.StatusUnauthorized, ageError)
		}
	}
	///using UpdataData to update data in database with this data///
	userID := UpdateData(bsonID, name, conAge, yearOfBirth, avatarName, note, contentType)

	return c.JSON(http.StatusCreated, GetUserData(userID))
}

//DeleteUser : function using with DeleteUserData//
func DeleteUser(c echo.Context) error {
	type status struct {
		Success bool `json:"success"`
	}
	id := c.Param("user_id")
	///check if id param send with invaild format (24-digit)///
	if len(id) != 24 {
		return c.HTML(http.StatusUnauthorized, "Invalid ID format. Plase try again")
	}
	///convert string to bson object///
	bsonID := bson.ObjectIdHex(id)

	///when not found user with this id///
	findUserError := &ErrorMessage{
		Code:        "404",
		Description: "User not found",
	}
	if GetUserData(bsonID) == (UserData{}) {
		return c.JSON(http.StatusNotFound, findUserError)
	}
	DeleteUserData(bsonID)
	return c.JSON(http.StatusCreated, &status{Success: true})
}
