package routes

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"golden-arm/internal"
	"golden-arm/schema"
	"html/template"
	"net/http"
	"net/smtp"
	"os"
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

// Confirmation email
type EmailData struct {
	To         string
	Name       string
	MovieTitle string
	MovieDate  string
	SeatNumber int
	PosterURL  string
}

/*
Reserves a seat and sends email confirmation
Raises error for invalid seat or conflicting reservation

	curl -X POST http://localhost:8080/api/reserve -H "Content-Type: application/json" -d
	'{
		"movie_id": "00000000-0000-0000-0000-000000000000",
		"seat_number": 4,
		"name": "Joey B",
		"email": "jb@example.com"
	}'
*/
func Reserve(c *gin.Context) {
	var newRes ReservationRequest
	if err := c.ShouldBindJSON(&newRes); err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}

	// Validate that seat number is in range [1, 18] (how many seats does golden arm actually have?)
	if newRes.SeatNumber < 1 || newRes.SeatNumber > 18 {
		fmt.Println("Seat number must be between 1 and 18")
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

	fmt.Printf("Reservation saved: %v\n", res)

	// Send confirmation email
	// Load movie details
	var movie schema.Movie
	err = db.NewSelect().
		Model(&movie).
		Where("id = ?", res.MovieID).
		Scan(ctx)
	if err != nil {
		fmt.Println("Error loading movie details: ", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}
	// Prepare email data
	var data EmailData
	data.To = res.Email
	data.Name = res.Name
	data.MovieTitle = movie.Title
	data.MovieDate = movie.Date.Format("Monday, January 2 3:04 PM")
	data.SeatNumber = res.SeatNumber
	data.PosterURL = movie.PosterURL

	if err := sendConfirmationEmail(data); err != nil {
		fmt.Printf("Error sending confirmation email: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Reservation confirmed"})
}

//go:embed templates/*
var emailTemplate embed.FS

func sendConfirmationEmail(data EmailData) error {
	// Parse and fill the HTML email template
	tmpl, err := template.ParseFS(emailTemplate, "templates/email.html")
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	// Configure SMTP settings
	smtpHost := "smtp.gmail.com"
	smtpPort := "587" // Common port for TLS
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// Prepare email message
	from := smtpUsername
	subject := fmt.Sprintf("You're set to watch \"%s\" @ The Golden Arm: %s", data.MovieTitle, data.MovieDate)
	message := fmt.Sprintf("Subject: %s\r\nFrom: %s\r\nTo: %s\r\nContent-Type: text/html\r\n\r\n%s",
		subject, from, data.To, body.String())

	// Connect to the SMTP server
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{data.To}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Printf("Confirmation email sent to %s\n", data.To)
	return nil
}

/*
Gets the seats that have been reserved for a movie

	curl -X GET http://localhost:8080/api/reserved/00000000-0000-0000-0000-000000000000
*/
func GetReservedSeats(c *gin.Context) {
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

	reservations, err := getReservations(movieID)
	if err != nil {
		fmt.Printf("Error fetching reservations: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	var reservedSeats []int
	for _, reservation := range reservations {
		reservedSeats = append(reservedSeats, reservation.SeatNumber)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"movie_id":       movieID,
			"reserved_seats": reservedSeats,
		},
	})
}

/*
Gets full reservation data for a movie including names and emails

	curl -X GET http://localhost:8080/api/reservations/00000000-0000-0000-0000-000000000000 \
	-H "Authorization: Bearer YOUR API KEY"
*/
func GetReservations(c *gin.Context) {
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

	reservations, err := getReservations(movieID)
	if err != nil {
		fmt.Printf("Error fetching reservations: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}
	if reservations == nil {
		reservations = []schema.Reservation{}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": reservations})
}

// Helper function returning all reservation data for a movie
// Returns error if movie does not exist
func getReservations(movieID uuid.UUID) ([]schema.Reservation, error) {
	db := schema.GetDBConn()
	ctx := context.Background()

	// Validate if movie exists in the database
	var movieExists uuid.UUID
	err := db.NewSelect().
		Model((*schema.Movie)(nil)).
		Where("id = ?", movieID).
		Column("id").
		Scan(ctx, &movieExists)

	if err != nil {
		fmt.Printf("Error checking movie existence: %v", err)
		return nil, internal.ErrInternalServer
	}
	if movieExists == uuid.Nil {
		fmt.Printf("Movie not found: %v", err)
		return nil, internal.ErrNotFound
	}

	// Fetch reservations for the movie
	var reservations []schema.Reservation
	err = db.NewSelect().
		Model(&reservations).
		Relation("Movie").
		Where("movie_id = ?", movieID).
		Scan(ctx)

	if err != nil {
		fmt.Printf("Error fetching reservations: %v", err)
		return nil, err
	}

	return reservations, nil
}

/*
Deletes reservation from database

	curl -X DELETE http://localhost:8080/api/reservation/00000000-0000-0000-0000-000000000000 \
	-H "Authorization: Bearer YOUR API KEY"
*/
func DeleteReservation(c *gin.Context) {
	if !internal.CheckAuthorization(c) {
		c.AbortWithError(http.StatusUnauthorized, internal.ErrUnauthorized)
		return
	}

	// Ensure reservation_id is provided and is a valid UUID
	param := c.Param("reservation_id")
	if param == "" {
		fmt.Println("reservation_id path parameter is required")
		c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}
	resID, err := uuid.Parse(param)
	if err != nil {
		fmt.Println("reservation_id must be a valid UUID")
		c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}

	db := schema.GetDBConn()
	ctx := context.Background()

	// Delete the reservation from the database
	result, err := db.NewDelete().
		Model((*schema.Reservation)(nil)).
		Where("id = ?", resID).
		Exec(ctx)

	if err != nil {
		fmt.Printf("Error deleting reservation: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		fmt.Println("Reservation not found")
		c.AbortWithError(http.StatusNotFound, internal.ErrNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Reservation deleted successfully"})
}
