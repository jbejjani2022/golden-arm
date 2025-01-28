package internal

import (
	"os"

	"github.com/gin-gonic/gin"
)

func CheckAuthorization(c *gin.Context) bool {
	// check for isAdmin cookie
	// if cookie is not present, check the request Authorization header
	// if API_KEY is not present in the request header or not equal to the value in .env, return 401

	isAdmin, err := c.Cookie("isAdmin")
	if err == nil && isAdmin == "true" {
		return true
	}

	authHeader := c.GetHeader("Authorization")
	return authHeader == "Bearer "+os.Getenv("API_KEY")
}
