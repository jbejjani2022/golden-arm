package internal

import (
	"os"

	"github.com/gin-gonic/gin"
)

func CheckAuthorization(c *gin.Context) bool {
	// read API_KEY Tfrom .env
	// if API_KEY is not present in the request header, return 401
	// if API_KEY is not equal to the value in .env, return 401
	// otherwise, continue

	var req_apikey = c.GetHeader("Authorization")
	var act_apikey = "Bearer " + os.Getenv("API_KEY")
	return req_apikey == act_apikey
}
