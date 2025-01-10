package main

import (
	"golden-arm/internal"
	"golden-arm/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	// Error-handling middleware
	router.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors[len(c.Errors)-1].Err // Process the last error only

			switch err {
			case internal.ErrBadRequest:
				internal.Handle400(c)
			case internal.ErrUnauthorized:
				internal.Handle401(c)
			case internal.ErrNotFound:
				internal.Handle404(c)
			case internal.ErrMethodNotAllowed:
				internal.Handle405(c)
			case internal.ErrNotImplemented:
				internal.Handle501(c)
			default:
				internal.Handle500(c)
			}
			// Abort after handling to prevent further processing
			c.Abort()
		}
	})
	router.NoRoute(internal.Handle404)
	router.NoMethod(internal.Handle405)

	// schema.CreateTables()

	// Routes
	router.GET("/api/movie", routes.GetMovie)
	router.GET("/api/movie/archive", routes.GetMovieArchive)
	router.GET("/api/seats", routes.GetSeats)
	router.GET("/api/comments", routes.GetComments)

	router.POST("/api/reserve", routes.Reserve)
	router.POST("/api/movie", routes.SetMovie)
	router.POST("/api/menu", routes.SubmitComment)

	router.Run(":8080")
}
