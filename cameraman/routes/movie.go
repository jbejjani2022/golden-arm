package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MovieRequest struct {
	Title     string `json:"title"`
	PosterUrl string `json:"poster_url"`
	Date      string `json:"date"`
}

// Gets movie whose screening date is closest in the future
// Returns an error if all movies are in the past
// e.g. get the current week's screening
func GetMovie(c *gin.Context) {
	// TODO: fetch current movie from database
	var currentMovie = gin.H{
		"title":      "Interstellar",
		"poster_url": "/assets/interstellar.jpg",
		"date":       "2024-12-01",
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": currentMovie})
}

// Adds new movie to database
// e.g. set the upcoming screening
func SetMovie(c *gin.Context) {
	var newMovie MovieRequest
	if err := c.ShouldBindJSON(&newMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: enter new movie into database

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Movie updated successfully"})
}

// Gets all past movies screened
func GetMovieArchive(c *gin.Context) {
	// TODO: fetch movie archive from database
	var movieArchive = []gin.H{
		{
			"title":      "Interstellar",
			"poster_url": "/assets/interstellar.jpg",
			"date":       "2024-12-01",
		},
		{
			"title":      "The Dark Knight",
			"poster_url": "/assets/dark_knight.jpg",
			"date":       "2024-11-24",
		},
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": movieArchive})
}
