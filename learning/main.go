package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
)

func main() {
	////new echo to e variable////
	e := echo.New()
	type User struct {
		Name  string `json:"name" xml:"name" form:"name" query:"name"`
		Email string `json:"email" xml:"email" form:"email" query:"email"`
	}
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello echo")
	// })

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	///get data from form or curl///
	// e.POST("/users", saveUser)
	e.POST("/save", save)

	///function to handing request///
	e.POST("/users", func(c echo.Context) error {
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}
		///Set output when json is valid///
		return c.JSON(http.StatusCreated, u)

	})

	///get data from ""//
	e.GET("/users/:id", getUser)
	e.GET("/show", show)

	///set directory of static file///
	e.Static("/static", "static")

	///get file from param 2 when going to param 1///
	e.File("/pic", "paradrop.png")

	// e.PUT("/users:id", updateUser)

	// e.DELETE("/users/:id", deleteUser)

	///filter group///
	g := e.Group("/boss")
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "Ham" && password == "password" {
			return true, nil
		}
		return false, nil
	}))

	////run in this port////
	e.Logger.Fatal(e.Start(":1323"))
}

////create function for get text from param id////
func getUser(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)

}

////create function for get text from param team and member
func show(c echo.Context) error {
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:"+team+", member:"+member)

}

////create function to save text from curl////
func save(c echo.Context) error {
	///save text only///
	// name := c.FormValue("name")
	// email := c.FormValue("email")
	// return c.String(http.StatusOK, "name: "+name+" email: "+email)

	///save mutiple data///
	//name := c.FormValue("name")
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

	o, err := os.Open("img/" + avatar.Filename)
	if err != nil {
		return err
	}
	defer o.Close()

	//destination to upload image//
	// dst, err := os.Create(avatar.Filename)
	// if err != nil {
	// 	return err
	// }
	// defer dst.Close()

	// //copy image from souce to destination//
	// if _, err = io.Copy(dst, src); err != nil {
	// 	return err
	// }
	contentType, err := getFlieType(o)
	if err != nil {
		return c.HTML(http.StatusBadRequest, "File incorrect!")
	}

	return c.HTML(http.StatusOK, " "+contentType)

}

func getFlieType(out *os.File) (string, error) {
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "File incorrect!", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
