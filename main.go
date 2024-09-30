package main

import (
	"fmt"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Response represents the JSON response structure
type Response struct {
	Message string `json:"message"`
}

// NameRequest represents the structure of the request body for the POST request
type NameRequest struct {
	Name string `json:"name"`
}

// handlePost handles the POST request to /name
func handlePost(c echo.Context) error {
	var nameReq NameRequest
	err := c.Bind(&nameReq) // Use Echo's Bind method to parse JSON
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
	}

	// Respond with the name
	response := fmt.Sprintf("Hello, %s!", nameReq.Name)
	return c.JSON(http.StatusOK, map[string]string{"message": response})
}

func main() {
	// Create a new Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route for handling GET request at "/"
	e.GET("/", func(c echo.Context) error {
		response := Response{
			Message: "Hello, this is your simple API!",
		}
		return c.JSON(http.StatusOK, response)
	})

	// Route for handling POST request at "/name"
	e.POST("/name", handlePost)

	// Start the server on port 8080
	fmt.Println("Server starting on port 8080...")
	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal("Shutting down the server")
	}
}
