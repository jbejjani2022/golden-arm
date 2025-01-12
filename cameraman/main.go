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
	router.GET("/api/movie/:movie_id", routes.GetMovie)
	router.GET("/api/movie/next", routes.GetNextMovie)
	router.GET("/api/movie/all", routes.GetAllMovies)
	router.GET("/api/movie/archive", routes.GetMovieArchive)
	router.GET("/api/reserved/:movie_id", routes.GetReservedSeats)
	router.GET("/api/reservations/:movie_id", routes.GetReservations)
	router.GET("/api/comments", routes.GetComments)

	router.POST("/api/reserve", routes.Reserve)
	router.POST("/api/movie", routes.AddMovie)
	router.POST("/api/comment", routes.SubmitComment)
	router.POST("/api/admin/login", routes.AdminLogin)
	router.POST("/api/admin/logout", routes.AdminLogout)

	router.DELETE("/api/movie/:movie_id", routes.DeleteMovie)
	router.DELETE("/api/reservation/:reservation_id", routes.DeleteReservation)
	router.DELETE("/api/comment/:comment_id", routes.DeleteComment)

	router.Run(":8080")
}
