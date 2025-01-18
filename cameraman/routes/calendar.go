package routes

import (
	"context"
	"fmt"
	"golden-arm/internal"
	"golden-arm/schema"
	"golden-arm/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CalendarRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	ImageURL  string    `json:"image_url"`
}

/*
Gets all calendars in the database

	curl -X GET http://localhost:8080/api/calendar/all -H "Authorization: Bearer YOUR API KEY"
*/
func GetAllCalendars(c *gin.Context) {
	if !internal.CheckAuthorization(c) {
		c.AbortWithError(http.StatusUnauthorized, internal.ErrUnauthorized)
		return
	}

	var calendars []schema.Calendar
	db := schema.GetDBConn()
	ctx := context.Background()

	err := db.NewSelect().
		Model(&calendars).
		Order("end_date DESC").
		Scan(ctx)
	if err != nil {
		fmt.Printf("Error fetching calendars: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": calendars})
}

/*
Gets calendar whose date range contains the current date
If such a calendar doesn't exist, get calendar whose start date is closest in the future
If no future calendars, get most recent calendar

	curl -X GET http://localhost:8080/api/calendar
*/
func GetCalendar(c *gin.Context) {
	var calendar schema.Calendar
	db := schema.GetDBConn()
	ctx := context.Background()

	// Try to find calendar whose date range contains the current date
	err := db.NewSelect().
		Model(&calendar).
		Where("start_date <= ? AND end_date >= ?", time.Now(), time.Now()).
		Limit(1).
		Scan(ctx)

	if err != nil {
		// Try to find the calendar whose start date is closest in the future
		err = db.NewSelect().
			Model(&calendar).
			Where("start_date > ?", time.Now()).
			Order("start_date ASC"). // closest future start date
			Limit(1).
			Scan(ctx)

		if err != nil {
			// Get most recent calendar
			err = db.NewSelect().
				Model(&calendar).
				Where("end_date < ?", time.Now()).
				Order("end_date DESC"). // most recent end date
				Limit(1).
				Scan(ctx)

			if err != nil {
				// No calendars in the database
				fmt.Printf("Error fetching calendar: %v", err)
				c.AbortWithError(http.StatusNotFound, internal.ErrNotFound)
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": calendar})
}

/*
Adds new calendar to database; supports file upload and JSON-based submissions

For JSON-based submissions:

	curl -X POST http://localhost:8080/api/calendar -H "Authorization: Bearer YOUR API KEY" \
	-H "Content-Type: application/json" -d
	'{
		"start_date": "2025-01-01T00:00:00Z",
		"end_date": "2025-02-01T00:00:00Z",
		"image_url": "https://example.com/calendar.jpg",
	}'

For file upload submissions:

	curl -X POST http://localhost:8080/api/calendar -H "Authorization: Bearer YOUR_API_KEY" \
		-F "start_date=2025-01-01T00:00:00Z" \
		-F "end_date=2025-02-01T00:00:00Z" \
		-F "image=@/path/to/image.jpg"
*/
func AddCalendar(c *gin.Context) {
	if !internal.CheckAuthorization(c) {
		c.AbortWithError(http.StatusUnauthorized, internal.ErrUnauthorized)
		return
	}

	// Check if the request is multipart/form-data for file uploads
	contentType := c.Request.Header.Get("Content-Type")
	isMultipart := strings.HasPrefix(contentType, "multipart/form-data")

	var newCalendar CalendarRequest
	if isMultipart {
		// Handle file uploads
		var err error

		newCalendar.StartDate, err = time.Parse(time.RFC3339, c.PostForm("start_date"))
		if err != nil {
			fmt.Println("Error parsing date:", err)
			c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
			return
		}
		newCalendar.EndDate, err = time.Parse(time.RFC3339, c.PostForm("end_date"))
		if err != nil {
			fmt.Println("Error parsing date:", err)
			c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
			return
		}

		// Calendar image file
		imageFile, _ := c.FormFile("image")
		if imageFile != nil {
			newCalendar.ImageURL, err = utils.UploadToS3(imageFile, "Calendars")
			if err != nil {
				fmt.Println("Error uploading calendar image file:", err)
				c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
				return
			}
		} else {
			newCalendar.ImageURL = c.PostForm("poster_url")
		}

	} else {
		// Handle JSON requests (for URL-based submissions)
		if err := c.ShouldBindJSON(&newCalendar); err != nil {
			fmt.Println(err)
			c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
			return
		}
	}

	// Create calendar object
	calendar := schema.Calendar{
		ID:        uuid.New(),
		StartDate: newCalendar.StartDate,
		EndDate:   newCalendar.EndDate,
		ImageURL:  newCalendar.ImageURL,
		Date:      time.Now(),
	}

	// Database connection
	db := schema.GetDBConn()
	ctx := context.Background()

	// Check if new calendar overlaps with any existing calendars
	var count int
	err := db.NewSelect().
		Model(&schema.Calendar{}).
		Where("start_date <= ? AND end_date >= ?", calendar.EndDate, calendar.StartDate).
		ColumnExpr("COUNT(*)").
		Scan(ctx, &count)
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}
	if count > 0 {
		fmt.Println("Calendar overlaps with existing calendars")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "error": "Calendar overlaps with existing calendar."})
		return
	}

	// Insert the new calendar into the database
	_, err = db.NewInsert().
		Model(&calendar).
		Exec(ctx)
	if err != nil {
		fmt.Printf("Error adding calendar to database: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Calendar added successfully"})
}

/*
Deletes calendar from database

	curl -X DELETE http://localhost:8080/api/calendar/00000000-0000-0000-0000-000000000000 \
	-H "Authorization: Bearer YOUR API KEY"
*/
func DeleteCalendar(c *gin.Context) {
	if !internal.CheckAuthorization(c) {
		c.AbortWithError(http.StatusUnauthorized, internal.ErrUnauthorized)
		return
	}

	// Ensure calendar_id is provided and is a valid UUID
	param := c.Param("calendar_id")
	if param == "" {
		fmt.Println("calendar_id path parameter is required")
		c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}
	calendarID, err := uuid.Parse(param)
	if err != nil {
		fmt.Println("calendar_id must be a valid UUID")
		c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}

	db := schema.GetDBConn()
	ctx := context.Background()

	// Delete the calendar from the database
	result, err := db.NewDelete().
		Model((*schema.Calendar)(nil)).
		Where("id = ?", calendarID).
		Exec(ctx)

	if err != nil {
		fmt.Printf("Error deleting calendar: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		fmt.Println("Calendar not found")
		c.AbortWithError(http.StatusNotFound, internal.ErrNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Calendar deleted successfully"})
}
