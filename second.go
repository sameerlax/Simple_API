package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type DataRequest struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Phone int    `json:"phone"`
	ID    int    `json:"id"`
}

var users = make(map[int]DataRequest) // In-memory storage for users

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/name", func(c echo.Context) error {
		var DataReq DataRequest
		if err := c.Bind(&DataReq); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
		}
		// Store the user in the map
		users[DataReq.ID] = DataReq
		response := fmt.Sprintf("Hello, %s! Your age: %d Phone: %d", DataReq.Name, DataReq.Age, DataReq.Phone)
		return c.JSON(http.StatusOK, map[string]string{"message": response})
	})

	// PUT method to update a user
	e.PUT("/name/:id", func(c echo.Context) error {
		id := c.Param("id")
		var DataReq DataRequest
		if err := c.Bind(&DataReq); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
		}
		intID, err := strconv.Atoi(id) // Convert string id to int
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}

		// Update the user in the map
		users[intID] = DataReq
		response := fmt.Sprintf("Updated user ID %d: %s, Age: %d, Phone: %d", intID, DataReq.Name, DataReq.Age, DataReq.Phone)
		return c.JSON(http.StatusOK, map[string]string{"message": response})
	})

	// DELETE method to delete a user
	e.DELETE("/name/:id", func(c echo.Context) error {
		id := c.Param("id")
		intID, err := strconv.Atoi(id) // Convert string id to int
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}
		// Simulate deleting the user
		delete(users, intID)
		response := fmt.Sprintf("Deleted user ID %d", intID)
		return c.JSON(http.StatusOK, map[string]string{"message": response})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
