package routes

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golden-arm/internal"
	"golden-arm/schema"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminLoginRequest struct {
	Passkey string `json:"passkey"`
}

// Generates a secure random session token
func generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func AdminLogin(c *gin.Context) {
	var request AdminLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}

	// Validate password
	if request.Passkey != os.Getenv("ADMIN_PASSKEY") {
		fmt.Println("Invalid passkey")
		c.AbortWithError(http.StatusUnauthorized, internal.ErrUnauthorized)
		return
	}

	sessionToken, err := generateSessionToken()
	if err != nil {
		fmt.Println("Failed to generate session token:", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	// Store the session in an in-memory store
	internal.StoreSession(sessionToken, "admin", time.Now().Add(1*time.Hour))

	// Set a cookie for session validation with lifetime 3600s = 1 hr
	c.SetCookie(
		"sessionToken",
		sessionToken,
		3600,
		"/", // Path (root scope)
		"goldenarmtheater.com",
		true, // Secure (only send over HTTPS)
		true, // HttpOnly (not accessible via JavaScript)
	)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Login successful"})
}

func AdminLogout(c *gin.Context) {
	sessionToken, _ := c.Cookie("sessionToken")
	if sessionToken != "" {
		internal.DeleteSession(sessionToken)
	}

	c.SetCookie(
		"sessionToken",
		"",
		-1, // MaxAge (negative value indicates deletion)
		"/",
		"goldenarmtheater.com",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Logout successful"})
}

// Validates session token
func ValidateSession(c *gin.Context) {
	sessionToken, err := c.Cookie("sessionToken")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"valid": false, "message": "Invalid request"})
		return
	}

	isValid := internal.ValidateSession(sessionToken)
	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false, "message": "Invalid session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true, "message": "Session is valid"})
}

/*
Gets list of all unique emails of anyone who has made a reservation or commented on the site

	curl -X GET http://localhost:8080/api/emails -H "Authorization: Bearer DO NOT USE IN PRODUCTION"
*/
func GetEmails(c *gin.Context) {
	if !internal.CheckAuthorization(c) {
		c.AbortWithError(http.StatusUnauthorized, internal.ErrUnauthorized)
		return
	}

	db := schema.GetDBConn()
	ctx := context.Background()

	var reservationEmails []string
	var commentEmails []string

	// Query unique emails from the reservations table
	err := db.NewSelect().
		Model((*schema.Reservation)(nil)).
		ColumnExpr("DISTINCT email").
		Scan(ctx, &reservationEmails)

	if err != nil {
		fmt.Printf("Error fetching reservation emails: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	// Query unique emails from the comments table
	err = db.NewSelect().
		Model((*schema.Comment)(nil)).
		ColumnExpr("DISTINCT email").
		Scan(c, &commentEmails)

	if err != nil {
		fmt.Printf("Error fetching comment emails: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	// Combine both email lists and remove duplicates
	emailSet := make(map[string]struct{})
	for _, email := range reservationEmails {
		emailSet[email] = struct{}{}
	}
	for _, email := range commentEmails {
		emailSet[email] = struct{}{}
	}

	// Convert the set back to a slice
	var uniqueEmails []string
	for email := range emailSet {
		uniqueEmails = append(uniqueEmails, email)
	}
	if uniqueEmails == nil {
		uniqueEmails = []string{}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": uniqueEmails})
}
