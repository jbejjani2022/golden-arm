package routes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golden-arm/internal"
	"golden-arm/schema"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReservationRequest struct {
	MovieID    uuid.UUID `json:"movie_id" binding:"required"`
	SeatNumber int       `json:"seat_number" binding:"required"`
	Name       string    `json:"name" binding:"required"`
	Email      string    `json:"email" binding:"required,email"`
}

// Reserves a seat
// Raises error for invalid seat or conflicting reservation
func Reserve(c *gin.Context) {
	var newRes ReservationRequest
	if err := c.ShouldBindJSON(&newRes); err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}

	// Validate that seat number is in range [1, 10] (how many seats does golden arm actually have?)
	if newRes.SeatNumber < 1 || newRes.SeatNumber > 10 {
		fmt.Println("Seat number must be between 1 and 10")
		c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}

	db := schema.GetDBConn()
	ctx := context.Background()

	// Check for conflicting reservation (same seat in same movie)
	var conflictingRes schema.Reservation
	err := db.NewSelect().
		Model(&conflictingRes).
		Where("movie_id = ? AND seat_number = ?", newRes.MovieID, newRes.SeatNumber).
		Scan(ctx)

	if err == nil {
		// Conflicting reservation found
		fmt.Printf("Seat %d already reserved", newRes.SeatNumber)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		// Handle unexpected errors
		fmt.Printf("Error checking seat availability: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	// Save new reservation
	res := schema.Reservation{
		ID:         uuid.New(),
		MovieID:    newRes.MovieID,
		SeatNumber: newRes.SeatNumber,
		Date:       time.Now(),
		Name:       newRes.Name,
		Email:      newRes.Email,
	}

	_, err = db.NewInsert().
		Model(&res).
		Exec(ctx)

	if err != nil {
		fmt.Printf("Error saving reservation: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

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
