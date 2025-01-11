package internal

import (
	"os"

	"github.com/gin-gonic/gin"
)

func CheckAuthorization(c *gin.Context) bool {
	// read API_KEY from .env
	// if API_KEY is not present as a cookie, check the request Authorization header
	// if API_KEY is not present in the request header or not equal to the value in .env, return 401
	// otherwise, continue

	apiKey, err := c.Cookie("apiKey")
	if err == nil && apiKey == os.Getenv("API_KEY") {
		return true
	}

	authHeader := c.GetHeader("Authorization")
	return authHeader == "Bearer "+os.Getenv("API_KEY")
}
