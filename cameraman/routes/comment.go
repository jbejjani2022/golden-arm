package routes

import (
	"context"
	"fmt"
	"golden-arm/internal"
	"golden-arm/schema"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Comment string `json:"comment"`
}

/*
Gets all movie-goer comments / suggestions

	curl -X GET http://localhost:8080/api/comments -H "Authorization: Bearer YOUR API KEY"
*/
func GetComments(c *gin.Context) {
	if !internal.CheckAuthorization(c) {
		c.AbortWithError(http.StatusUnauthorized, internal.ErrUnauthorized)
		return
	}

	var comments []schema.Comment
	db := schema.GetDBConn()
	ctx := context.Background()

	// Fetch all comments from the database
	err := db.NewSelect().
		Model(&comments).
		Order("date DESC").
		Scan(ctx)
	if err != nil {
		fmt.Printf("Error fetching comments: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": comments})
}

/*
Submits new comment to database

	curl -X POST http://localhost:8080/api/comment -H "Content-Type: application/json" -d
	'{
		"name": "Joey B",
		"email": "jb@example.com",
		"comment": "The golden arm should screen Parasite"
	}'
*/
func SubmitComment(c *gin.Context) {
	var newComment CommentRequest
	if err := c.ShouldBindJSON(&newComment); err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}

	comment := schema.Comment{
		ID:      uuid.New(),
		Name:    newComment.Name,
		Email:   newComment.Email,
		Comment: newComment.Comment,
		Date:    time.Now(),
	}

	db := schema.GetDBConn()
	ctx := context.Background()

	// Insert the new comment into the database
	_, err := db.NewInsert().
		Model(&comment).
		Exec(ctx)
	if err != nil {
		fmt.Printf("Error inserting comment: %v\n", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Comment submitted"})
}
