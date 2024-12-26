package main

import (
	"golden-arm/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// TODO: load environment variables

	r := gin.Default()

	// TODO: error handling middleware

	// Routes
	r.GET("/api/movie", routes.GetMovie)
	r.GET("/api/movie/archive", routes.GetMovieArchive)
	r.GET("/api/menu", routes.GetMenu)
	r.GET("/api/menu/archive", routes.GetMenuArchive)
	r.GET("/api/seats", routes.GetSeats)

	r.POST("/api/movie", routes.SetMovie)
	r.POST("/api/menu", routes.SetMenu)
	r.POST("/api/reserve", routes.Reserve)
	r.POST("/api/reset-seats", routes.ResetSeats)

	r.Run(":8080")
}
