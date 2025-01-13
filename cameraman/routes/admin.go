package routes

import (
	"context"
	"fmt"
	"golden-arm/internal"
	"golden-arm/schema"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AdminLoginRequest struct {
	Passkey string `json:"passkey"`
}

// Handles admin passkey validation
func AdminLogin(c *gin.Context) {
	var request AdminLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}

	if request.Passkey != os.Getenv("ADMIN_PASSKEY") {
		fmt.Println("Invalid passkey")
		c.AbortWithError(http.StatusUnauthorized, internal.ErrUnauthorized)
		return
	}

	// Set a cookie for session validation with lifetime 3600s = 1 hr
	// For production
	// change localhost to site domain
	// change `secure` from false to true to only send cookie over https
	c.SetCookie("isAdmin", "true", 3600, "/", "localhost", false, true)
	c.SetCookie("apiKey", os.Getenv("API_KEY"), 3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Login successful"})
}

func AdminLogout(c *gin.Context) {
	c.SetCookie("isAdmin", "", -1, "/", "localhost", false, true)
	c.SetCookie("apiKey", "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Logout successful"})
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
