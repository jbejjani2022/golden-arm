package main

import (
	"golden-arm/internal"
	"golden-arm/routes"
	"golden-arm/schema"
	"os"

	// Add this line
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

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

	schema.CreateTables()

	// Routes
	router.GET("/api/movie/:movie_id", routes.GetMovie)
	router.GET("/api/movie/next", routes.GetNextMovie)
	router.GET("/api/movie/all", routes.GetAllMovies)
	router.GET("/api/movie/archive", routes.GetMovieArchive)
	router.GET("/api/reserved/:movie_id", routes.GetReservedSeats)
	router.GET("/api/reservations/:movie_id", routes.GetReservations)
	router.GET("/api/comments", routes.GetComments)
	router.GET("/api/emails", routes.GetEmails)
	router.GET("/api/calendar", routes.GetCalendar)
	router.GET("/api/calendar/all", routes.GetAllCalendars)
	router.GET("/api/merch/all", routes.GetAllMerchandise)
	router.GET("/api/order/all", routes.GetAllOrders)

	router.POST("/api/reserve", routes.Reserve)
	router.POST("/api/movie", routes.AddMovie)
	router.POST("/api/comment", routes.SubmitComment)
	router.POST("/api/calendar", routes.AddCalendar)
	router.POST("/api/admin/login", routes.AdminLogin)
	router.POST("/api/admin/logout", routes.AdminLogout)
	router.POST("/api/admin/validate-session", routes.ValidateSession)
	router.POST("/api/merch", routes.AddMerchandise)
	router.POST("/api/order", routes.AddOrder)

	router.PUT("/api/merch/:merch_id", routes.UpdateMerchandise)
	router.PUT("/api/order/status/:order_id", routes.UpdateOrderStatus)
	router.PUT("/api/movie/:movie_id", routes.UpdateMovie)

	router.DELETE("/api/movie/:movie_id", routes.DeleteMovie)
	router.DELETE("/api/reservation/:reservation_id", routes.DeleteReservation)
	router.DELETE("/api/comment/:comment_id", routes.DeleteComment)
	router.DELETE("/api/calendar/:calendar_id", routes.DeleteCalendar)
	router.DELETE("/api/merch/:merch_id", routes.DeleteMerchandise)
	router.DELETE("/api/order/:order_id", routes.DeleteOrder)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
