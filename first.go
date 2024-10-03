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

// Team struct represents the Team model
type Team struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Users []User `json:"users" gorm:"foreignKey:TeamID"` // One-to-many relationship
}

// Define the User struct, which will represent a table in PostgreSQL
type User struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Phone  string `json:"phone"`
	TeamID uint   `json:"team_id"` //foreign key
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
	db.AutoMigrate(&Team{}, &User{})
	//db.AutoMigrate(&User{})
}

func main() {
	// Initialize the database connection
	initDB()

	// Initialize Echo
	e := echo.New()

	// TEAM CODE
	// GET request to fetch all teams
	e.GET("/teams", func(c echo.Context) error {
		var teams []Team
		if err := db.Preload("Users").Find(&teams).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch teams"})
		}
		return c.JSON(http.StatusOK, teams)
	})

	// GET request to fetch a specific team by ID
	e.GET("/teams/:id", func(c echo.Context) error {
		id := c.Param("id")
		var team Team
		if err := db.First(&team, id).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Team not found"})
		}
		return c.JSON(http.StatusOK, team)
	})

	// POST request to create a new team
	e.POST("/teams", func(c echo.Context) error {
		team := new(Team)
		if err := c.Bind(team); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		db.Create(team) //Insert the new team into database
		return c.JSON(http.StatusCreated, team)
	})

	// PUT request to update the team with ID
	e.PUT("/teams/:id", func(c echo.Context) error {
		id := c.Param("id")
		var team Team

		// Attempt to find the team by ID
		if err := db.First(&team, id).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Team not found"})
		}

		// Bind the updated team data from the request body
		if err := c.Bind(&team); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
		}

		// Validation : check if the team name is valid
		if team.Name == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Team name is required"})
		}

		if len(team.Name) < 3 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Team name must be at least 3 characters long"})
		}

		// Save the updated team back to the database
		if err := db.Save(&team).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update team"})
		}

		// Return the updated team as JSON
		return c.JSON(http.StatusOK, team)
	})

	// DELETE request to delete  a team by ID
	e.DELETE("/teams/:id", func(c echo.Context) error {
		id := c.Param("id")
		var team Team

		//Attempt to find the team by ID
		if err := db.First(&team, id).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Team not found"})
		}
		// Delete the team from database
		if err := db.Delete(&team).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete team"})
		}
		// Return a success message
		return c.JSON(http.StatusOK, team)
	})
	// USER CODE
	// GET request to fetch all users
	e.GET("/User", func(c echo.Context) error {
		// Create a slice to hold multiple users
		var user []User
		// Retrieve all users from the database
		result := db.Find(&user)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users"})
		}
		// Return the list of users as JSON
		return c.JSON(http.StatusOK, user)
	})

	// GET request to fetch  user WITH ID
	e.GET("/User/:id", func(c echo.Context) error {
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
		// Manual Validation
		if user.Name == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user name"})
		}
		if user.Age < 0 || user.Age > 150 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user age"})
		}
		if user.Phone == "" || len(user.Phone) < 10 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user phone"})
		}
		for _, char := range user.Phone {
			if char < '0' || char > '9' {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user phone"})
			}
		}

		db.Create(&user) // Insert the new user into the database
		return c.JSON(http.StatusOK, user)
	})

	// PUT request to  update the user with ID
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
		// Manual Validation
		if user.Name == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user name"})
		}
		if user.Age < 0 || user.Age > 150 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user age"})
		}
		if user.Phone == "" || len(user.Phone) < 10 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user phone"})
		}
		for _, char := range user.Phone {
			if char < '0' || char > '9' {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user phone"})
			}
		}
		// Update the user response
		response := fmt.Sprintf("Updated user ID %d: %s, Age: %d, Phone: %s", intID, user.Name, user.Age, user.Phone)
		return c.JSON(http.StatusOK, map[string]string{"message": response})

	})

	// DELETE request to delete  user by ID
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
		return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
	})
	// Start the server
	e.Logger.Fatal(e.Start(":1323"))

}
