package internal

import (
	"os"

	"github.com/gin-gonic/gin"
)

func CheckAuthorization(c *gin.Context) bool {
	// check for and validate sessionToken cookie
	// if not present or invalid, check the request Authorization header
	// if API_KEY is not present in the request header or not equal to the value in .env, return 401

	sessionToken, err := c.Cookie("sessionToken")
	if err == nil && ValidateSession(sessionToken) {
		return true
	}

	authHeader := c.GetHeader("Authorization")
	return authHeader == "Bearer "+os.Getenv("API_KEY")
}
