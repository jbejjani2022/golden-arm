package routes

import (
	"fmt"
	"golden-arm/internal"
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
