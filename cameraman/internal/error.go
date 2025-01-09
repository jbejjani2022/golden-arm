package internal

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrUnauthorized     = fmt.Errorf("Unauthorized")
	ErrBadRequest       = fmt.Errorf("Bad Request")
	ErrNotFound         = fmt.Errorf("Not Found")
	ErrMethodNotAllowed = fmt.Errorf("Method Not Allowed")
	ErrInternalServer   = fmt.Errorf("Internal Server Error")
	ErrNotImplemented   = fmt.Errorf("Not Implemented")
)

func Handle400(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Bad Request"})
}

func Handle401(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
}

func Handle404(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Not found"})
}

func Handle405(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"success": false, "error": "Method Not Allowed"})
}

func Handle500(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Internal Server Error"})
}

func Handle501(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"success": false, "error": "Not Implemented"})
}
