package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MenuRequest struct {
	MovieID int    `json:"movie_id" binding:"required"`
	MenuUrl string `json:"menu_url"`
}

// Gets most recent menu added to database
// i.e. the current week's menu
func GetMenu(c *gin.Context) {
	// TODO: fetch current menu from database
	// Could do this by fetching current week's movie and using ID to fetch menu
	var currentMenu = gin.H{
		"movie_id": 2,
		"menu_url": "/assets/interstellar_menu.jpg",
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": currentMenu})
}

// Adds new menu to database
// i.e. sets the upcoming screening's menu
func SetMenu(c *gin.Context) {
	var newMenu MenuRequest
	if err := c.ShouldBindJSON(&newMenu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: enter new menu into database

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Menu updated successfully"})
}

// Gets all past menus
func GetMenuArchive(c *gin.Context) {
	// TODO: fetch menu archive from database
	var menuArchive = []gin.H{
		{
			"movie_id": 2,
			"menu_url": "/assets/interstellar_menu.jpg",
		},
		{
			"movie_id": 1,
			"menu_url": "/assets/dark_knight_menu.jpg",
		},
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": menuArchive})
}
