package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var validate *validator.Validate

// Helper method to set `db` to SQLite connection so we can make queries to `hackathons.db`
func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("hackathons.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&Hackathon{}) // keeps our schema up to date
	if err != nil {
		return
	}

	db = database
}

// Defines what a `Hackathon` is
type Hackathon struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Date     string `json:"date"`
	Url      string `json:"url"`
	Location string `json:"location"`
}

type HackathonInput struct {
	Name     string `json:"name" validate:"required"`
	Date     string `json:"date" validate:"required"`
	Url      string `json:"url" validate:"required"`
	Location string `json:"location" validate:"required"`
}

type HackathonDateInput struct {
	Date string `json:"date" validate:"required"`
}

// Get all hackathons
func getHackathons(c *gin.Context) {
	// Code goes here
	var hackathons []Hackathon
	db.Find(&hackathons)
	c.IndentedJSON(http.StatusOK, gin.H{"data": hackathons})
}

// GET /hackathons/:id
// Get hackathon by ID
func getHackathonById(c *gin.Context) {
	// Code goes here
	var hackathon Hackathon

	// error value should be nil (nothing went wrong)
	err := db.Where("id = ?", c.Param("id")).First(&hackathon).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hackathon not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": hackathon})

}

// POST /hackathons
// Create a hackathon
func createHackathon(c *gin.Context) {
	var input HackathonInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = validate.Struct(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var hackathons []Hackathon

	hackathon := Hackathon{
		Id:       int(db.Find(&hackathons).RowsAffected) + 1,
		Name:     input.Name,
		Date:     input.Date,
		Url:      input.Url,
		Location: input.Location,
	}
	db.Create(&hackathon)
	c.JSON(http.StatusOK, gin.H{"data": hackathon})

}

// PATCH /hackathons/:id
// Update a hackathon
func updateHackathon(c *gin.Context) {
	var hackathon Hackathon
	err := db.Where("id = ?", c.Param("id")).First(&hackathon).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hackathon not found!"})
		return
	}

	var input HackathonInput

	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = validate.Struct(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&hackathon).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": hackathon})

}

func updateHackathonDate(c *gin.Context) {
	var hackathon Hackathon
	err := db.Where("id = ?", c.Param("id")).First(&hackathon).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hackathon not found!"})
		return
	}

	var input HackathonDateInput

	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = validate.Struct(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	db.Model(&hackathon).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": hackathon})

}

// DELETE /hackathons/:id
// Delete a hackathon
func deleteHackathon(c *gin.Context) {
	// Code goes here
	var hackathon Hackathon

	err := db.Where("id = ?", c.Param("id")).First(&hackathon).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hackathon not found!"})
		return
	}

	db.Delete(&hackathon)
	c.JSON(http.StatusOK, gin.H{"data": true})

}

func main() {
	router := gin.Default()
	validate = validator.New()

	ConnectDatabase()

	router.GET("/hackathons", getHackathons)
	router.GET("/hackathons/:id", getHackathonById)
	router.POST("/hackathons", createHackathon)
	router.PATCH("/hackathons/:id", updateHackathon)
	router.DELETE("/hackathons/:id", deleteHackathon)
	router.PATCH("/hackathons/:id/date", updateHackathonDate)

	// get hackathon by name

	router.Run("localhost:8080")
}
