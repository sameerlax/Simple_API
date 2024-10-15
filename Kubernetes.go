package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type RequestBody struct {
	Message string `json:"message" validate:"required"`
}

func main() {

	// Initialize Echo
	e := echo.New()

	e.GET("/service", func(c echo.Context) error {

		return c.JSON(http.StatusOK, "This is From Service:) ")

	})
	e.GET("/service/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, " Hey Sir  :) ")
	})
	e.POST("/service/new", func(c echo.Context) error {
		var body RequestBody
		//Bind the request body to the RequestBody
		if err := c.Bind(&body); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid data"})
		}
		return c.JSON(http.StatusOK, body)
	})
	// Start the server
	e.Logger.Fatal(e.Start(":1323"))
}
