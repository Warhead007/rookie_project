package main

import (
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/labstack/echo"
)

type userData struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Age    int    `json:"age"`
	Note   string `json:"note"`
	Email  string `json:"email"`
}

func main() {

	e := echo.New()
	e.POST("/user", createUser)
	e.POST("/save", getData)

	e.Logger.Fatal(e.Start(":1323"))
}

func createUser(c echo.Context) error {
	u := new(userData)
	if err := c.Bind(u); err != nil {
		return err
	}
	///check require data///
	if u.Age == 0 || u.Name == "" || u.Email == "" || u.Avatar == "" {
		return c.String(http.StatusBadRequest, "Plase enter data")
	}
	return c.JSON(http.StatusCreated, u)
}

func getData(c echo.Context) error {
	name := c.FormValue("name")
	age := c.FormValue("age")
	email := c.FormValue("email")
	note := c.FormValue("note")
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return err
	}
	//souce of image//
	src, err := avatar.Open()
	if err != nil {
		return err
	}
	defer src.Close()

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
	///check file type///
	if contentType != "image/png" && contentType != "image/jpg" {
		os.Remove(fileName)
		return c.String(http.StatusUnauthorized, "File type not correct.")
	}
	//check email format//
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	///check input data///
	if len(name) == 0 {
		return c.String(http.StatusUnauthorized, "Plase enter your name.")
	}
	if len(age) == 0 {
		return c.String(http.StatusUnauthorized, "Plase enter your age.")
	}
	if len(email) == 0 {
		return c.String(http.StatusUnauthorized, "Plase enter your email.")
		///if email format invavid///
	} else if !re.MatchString(email) {
		return c.String(http.StatusUnauthorized, "Plase enter correct format of email.")
	}

	return c.HTML(http.StatusOK, name+" "+age+" "+email+" "+note+" "+contentType+" ")
}

///function to get file type///
func getFileType(out *os.File) (string, error) {
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "buffer incorrect", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
