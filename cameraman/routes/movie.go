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

type MovieRequest struct {
	Title     string    `json:"title"`
	Date      time.Time `json:"date"`
	PosterUrl string    `json:"poster_url"`
	MenuUrl   string    `json:"menu_url"`
}

/*
Adds new movie to database
e.g. set the upcoming screening

	curl -X POST http://localhost:8080/api/movie -H "Authorization: Bearer YOUR API KEY" \
	-H "Content-Type: application/json" -d
	'{
		"title": "Interstellar",
		"date": "2025-01-10T00:00:00Z",
		"poster_url": "https://example.com/poster.jpg",
		"menu_url": "https://example.com/menu.jpg"
	}'
*/
func AddMovie(c *gin.Context) {
	if !internal.CheckAuthorization(c) {
		c.AbortWithError(http.StatusUnauthorized, internal.ErrUnauthorized)
		return
	}

	var newMovie MovieRequest
	if err := c.ShouldBindJSON(&newMovie); err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}

	movie := schema.Movie{
		ID:        uuid.New(),
		Title:     newMovie.Title,
		Date:      newMovie.Date,
		PosterURL: newMovie.PosterUrl,
		MenuURL:   newMovie.MenuUrl,
	}

	db := schema.GetDBConn()
	ctx := context.Background()

	// Perform an upsert operation based on the Date field
	// e.g. if a movie already exists with identical Date value, update the remaining fields
	err := db.NewInsert().
		Model(&movie).
		On("CONFLICT (date) DO UPDATE").
		Set("title = EXCLUDED.title").
		Set("poster_url = EXCLUDED.poster_url").
		Set("menu_url = EXCLUDED.menu_url").
		Returning("id").
		Scan(ctx, &movie.ID)

	if err != nil {
		fmt.Printf("Error adding movie to database: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Movie added successfully"})
}

/*
Gets movie closest in the future; e.g. get the upcoming screening info
If none, gets most recent past screening

	curl -X GET http://localhost:8080/api/movie
*/
func GetMovie(c *gin.Context) {
	var nextMovie schema.Movie
	db := schema.GetDBConn()
	ctx := context.Background()

	// Try to find the closest upcoming movie
	err := db.NewSelect().
		Model(&nextMovie).
		Where("date > ?", time.Now()).
		Order("date ASC"). // closest future date
		Limit(1).
		Scan(ctx)

	if err != nil {
		// Try to find the most recent past movie
		err = db.NewSelect().
			Model(&nextMovie).
			Where("date <= ?", time.Now()).
			Order("date DESC"). // most recent past date
			Limit(1).
			Scan(ctx)

		// No movies in the database
		if err != nil {
			fmt.Printf("Error fetching movie: %v", err)
			c.AbortWithError(http.StatusNotFound, internal.ErrNotFound)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": nextMovie})
}

/*
Gets all movies in the database

	curl -X GET http://localhost:8080/api/movie/all -H "Authorization: Bearer YOUR API KEY"
*/
func GetAllMovies(c *gin.Context) {
	if !internal.CheckAuthorization(c) {
		c.AbortWithError(http.StatusUnauthorized, internal.ErrUnauthorized)
		return
	}

	var movies []schema.Movie
	db := schema.GetDBConn()
	ctx := context.Background()

	// Fetch all movies from the database
	err := db.NewSelect().
		Model(&movies).
		Order("date DESC").
		Scan(ctx)
	if err != nil {
		fmt.Printf("Error fetching movies: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": movies})
}

/*
Gets all past movies screened

	curl -X GET http://localhost:8080/api/movie/archive
*/
func GetMovieArchive(c *gin.Context) {
	var movieArchive []schema.Movie
	db := schema.GetDBConn()
	ctx := context.Background()

	// Select all movies whose screening date is strictly in the past
	err := db.NewSelect().
		Model(&movieArchive).
		Where("date < ?", time.Now()).
		Order("date DESC").
		Scan(ctx)

	if err != nil {
		fmt.Printf("Error fetching movie archive: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": movieArchive})
}

/*
Deletes movie from database

	curl -X DELETE http://localhost:8080/api/movie/00000000-0000-0000-0000-000000000000 \
	-H "Authorization: Bearer YOUR API KEY"
*/
func DeleteMovie(c *gin.Context) {
	if !internal.CheckAuthorization(c) {
		c.AbortWithError(http.StatusUnauthorized, internal.ErrUnauthorized)
		return
	}

	// Ensure movie_id is provided and is a valid UUID
	param := c.Param("movie_id")
	if param == "" {
		fmt.Println("movie_id path parameter is required")
		c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}
	movieID, err := uuid.Parse(param)
	if err != nil {
		fmt.Println("movie_id must be a valid UUID")
		c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}

	db := schema.GetDBConn()
	ctx := context.Background()

	// Delete the movie from the database
	result, err := db.NewDelete().
		Model((*schema.Movie)(nil)).
		Where("id = ?", movieID).
		Exec(ctx)

	if err != nil {
		fmt.Printf("Error deleting movie: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		fmt.Println("Movie not found")
		c.AbortWithError(http.StatusNotFound, internal.ErrNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Movie deleted successfully"})
}
