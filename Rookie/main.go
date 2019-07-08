package main

import (
	model "trainer/Rookie/pkg"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/user/:user_id", model.GetUser)
	e.GET("/users", model.GetAllUser)

	e.POST("/save", model.AddUser)

	e.PUT("/user/:user_id", model.UpdateUser)

	e.DELETE("/user/:user_id", model.DeleteUser)

	e.Logger.Fatal(e.Start(":1323"))
}
