package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Comment string `json:"comment"`
}

// Gets all movie-goer comments / suggestions
func GetComments(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": nil})
}

// Submits new comment to database
func SubmitComment(c *gin.Context) {
	var newComment CommentRequest
	if err := c.ShouldBindJSON(&newComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: enter new comment into database

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Comment submitted"})
}
