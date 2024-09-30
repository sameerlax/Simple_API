package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//type team struct {
//	ID   unit   `json:"id" gorm:"primaryKey"`
//	Name string `json:"name"`
//}

// Define the User struct, which will represent a table in PostgreSQL
type User struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Phone string `json:"phone"`
}

var db *gorm.DB
var err error

func initDB() {
	// Set up the PostgreSQL connection
	dsn := "host=localhost user=postgres password=2023 dbname=yourdb port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	// Automatically migrate the schema (create the "users" table)
	db.AutoMigrate(&User{})
}

func main() {
	// Initialize the database connection
	initDB()

	// Initialize Echo
	e := echo.New()

	// GET request to fetch all users
	e.GET("/users/:id", func(c echo.Context) error {
		id := c.Param("id")
		var user User
		// Convert string id to int
		intID, err := strconv.Atoi(id) // Convert string id to int
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}

		// Retrieve the user by ID from the database
		result := db.First(&user, intID) // Use `First` to fetch the user with the given id
		if result.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}

		return c.JSON(http.StatusOK, user) // Return the found user as JSON
	})

	// POST request to create a new user
	e.POST("/User", func(c echo.Context) error {
		user := new(User)
		if err := c.Bind(user); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		db.Create(&user) // Insert the new user into the database
		return c.JSON(http.StatusOK, user)
	})

	e.PUT("/User/:id", func(c echo.Context) error {

		id := c.Param("id")
		var user User
		// Convert string id to int
		intID, err := strconv.Atoi(id) // Convert string id to int
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}

		// Find the user by ID
		result := db.Where("id = ?", intID).First(&user)
		if result.Error != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invaild user ID"})
		}

		// Bind the updated user data from the request body
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
		}

		// Save the updated user back to the database
		if err := db.Save(&user).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed To Update user"})
		}

		// Update the user response

		response := fmt.Sprintf("Updated user ID %d: %s, Age: %d, Phone: %s", intID, user.Name, user.Age, user.Phone)
		return c.JSON(http.StatusOK, map[string]string{"message": response})

	})

	// DELETE request to create a new user
	e.DELETE("/User/:id", func(c echo.Context) error {
		id := c.Param("id")
		var user User
		// Convert string id to int
		intID, err := strconv.Atoi(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}
		// Find the user by ID
		if err := db.Where("id = ?", intID).First(&user).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "User not found "})
		}
		// Delete the user
		if err := db.Delete(&user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "User not found "})
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "User Deleted Successfully"})
	})
	// Start the server
	e.Logger.Fatal(e.Start(":1323"))
}
