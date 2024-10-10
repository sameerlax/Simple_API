package main

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Team struct represents the Team model
type Team struct {
	ID    uint    `json:"id" gorm:"primaryKey"`
	Name  string  `json:"name"`
	Users []*User `json:"users,omitempty"  gorm:"foreignKey:TeamID" ` //One-to-many relationship
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
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
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
	// GET request to fetch all teams Without User data
	e.GET("/service1/teams", func(c echo.Context) error {
		// create a slice to hold multiple teams
		var teams []Team
		// Retrieve all the teams from database
		result := db.Find(&teams)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Team not found"})
		}
		return c.JSON(http.StatusOK, teams)

	})

	// GET request to fetch all teams With User data
	e.GET("/service1/teams/users", func(c echo.Context) error {
		var teams []Team
		result := db.Preload("Users").Find(&teams)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Team,User Not Found"})
		}
		return c.JSON(http.StatusOK, teams)
	})

	// GET request to fetch a specific team by ID
	e.GET("/service1/teams/:id", func(c echo.Context) error {

		var team Team
		idParam := c.Param("id")

		// Convert idParam to uint
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid team ID"})
		}

		if err := db.Find(&team, id).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Team not found"})
		}
		return c.JSON(http.StatusOK, team)
	})

	// POST request to create a new team
	e.POST("/service1/teams", func(c echo.Context) error {
		team := new(Team)
		if err := c.Bind(team); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		// Validation for team
		if team.Name == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Team Name is required"})
		}
		db.Create(team) //Insert the new team into database
		return c.JSON(http.StatusCreated, team)
	})

	// POST request to create multiple users for a team
	e.POST("/service1/teams/:id/users", func(c echo.Context) error {
		idParam := c.Param("id") // Team ID as string

		// Convert idParam to uint
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid team ID"})
		}

		var users []User

		// Bind request body to users
		if err := c.Bind(&users); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid data Str"})
		}

		// Set TeamID for each user
		for i := range users {
			users[i].TeamID = uint(id) // Convert id to uint and assign

		}

		// Insert users into the database
		if err := db.Create(&users).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create users"})
		}

		// Return the created users
		return c.JSON(http.StatusCreated, users)
	})

	// PUT request to update the team with ID
	e.PUT("/service1/teams/:id", func(c echo.Context) error {
		id := c.Param("id")
		//var team Team
		team := Team{}
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
	e.DELETE("/service1/teams/:id", func(c echo.Context) error {
		id := c.Param("id")
		var team Team

		//Attempt to find the team by ID
		if err := db.First(&team, id).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Team not found"})
		}

		// Delete the team from database
		if err := db.Debug().Delete(&team).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete team"})
		}
		// Return a success message
		return c.JSON(http.StatusOK, team)
	})
	// Remove user from team
	e.DELETE("/service1/teams/:Tid/users/:Uid", func(c echo.Context) error {
		var users []User
		teamId := c.Param("Tid")
		userId := c.Param("Uid")
		// Finding user associated with team
		if err := db.Where("team_id = ? AND id = ? ", teamId, userId).First(&users).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Team not found"})
		}
		// If no user found , return a message
		if len(users) == 0 {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}

		// Delete the found user
		if err := db.Delete(&users).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete users"})
		}
		return c.JSON(http.StatusOK, users)
	})

	// USER CODE
	// GET request to fetch all users
	e.GET("/service1/User", func(c echo.Context) error {
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
	e.GET("/service1/User/:id", func(c echo.Context) error {
		id := c.Param("id")
		var user User
		// Convert string id to int
		intID, err := strconv.Atoi(id) // Convert string id to int
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}

		// Retrieve the user by ID from the database
		result := db.First(&user, intID).Error // Use `First` to fetch the user with the given id
		if result.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}

		return c.JSON(http.StatusOK, user) // Return the found user as JSON
	})

	// POST request to create a new user
	e.POST("/service1/User", func(c echo.Context) error {
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
		return c.JSON(http.StatusCreated, user)
	})

	// PUT request to  update the user with ID
	e.PUT("/service1/User/:id", func(c echo.Context) error {

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

		// Save the updated user back to the database
		if err := db.Save(&user).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed To Update user"})
		}

		// Update the user response
		response := fmt.Sprintf("Updated user ID %d: %s, Age: %d, Phone: %s", intID, user.Name, user.Age, user.Phone)
		return c.JSON(http.StatusOK, map[string]string{"message": response})

	})

	// DELETE request to delete  user by ID
	e.DELETE("/service1/User/:id", func(c echo.Context) error {
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
