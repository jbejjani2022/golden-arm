package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReservationRequest struct {
	MovieID    int    `json:"movie_id" binding:"required"`
	SeatNumber int    `json:"seat_number" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
}

// Reserves a seat and emails confirmation
func Reserve(c *gin.Context) {
	var req ReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Save reservation to the database
	// TODO: Send email confirmation

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Reservation confirmed"})
}

// Gets all available seats
func GetSeats(c *gin.Context) {
	// TODO: Fetch available seats from the database
	var availableSeats = []gin.H{
		{
			"id": 1,
		},
		{
			"id": 2,
		},
		{
			"id": 3,
		},
		{
			"id": 4,
		},
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": availableSeats})
}

// Resets all seats to available
func ResetSeats(c *gin.Context) {
	// TODO: Reset all seats to available
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Seats reset successfully"})
}
