package routes

import (
	"context"
	"fmt"
	"golden-arm/internal"
	"golden-arm/schema"
	"golden-arm/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MerchandiseRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	ImageURL    string     `json:"image_url"`
	Sizes       []SizeInfo `json:"sizes"`
}

type SizeInfo struct {
	Size     string `json:"size"`
	Quantity int    `json:"quantity"`
}

/*
Adds new merchandise item to the database, including its available sizes; supports file upload and JSON-based submissions

	For JSON-based submissions:

	curl -X POST http://localhost:8080/api/merch -H "Authorization: Bearer YOUR API KEY" \
		-H "Content-Type: application/json" -d
		'{
			"name": "Movie T-Shirt",
			"description": "Put a description here",
			"price": 15.00,
			"image_url": "https://example.com/images/movie-tshirt.jpg",
			"sizes": [
				{
				"size": "S",
				"quantity": 10
				},
				{
				"size": "M",
				"quantity": 15
				},
				{
				"size": "L",
				"quantity": 20
				},
				{
				"size": "XL",
				"quantity": 0
				}
			]
		}'

	For file upload submissions:

	curl -X POST http://localhost:8080/api/merch -H "Authorization: Bearer YOUR API KEY" \
		-F "name=Movie T-Shirt" \
		-F "description=Put a description here" \
		-F "price=15.00" \
		-F "image=@/path/to/image.jpg"
		-F "sizes=S,10" \
		-F "sizes=M,15" \
		-F "sizes=L,20" \
		-F "sizes=XL,0"
*/
func AddMerchandise(c *gin.Context) {
	if !internal.CheckAuthorization(c) {
		c.AbortWithError(http.StatusUnauthorized, internal.ErrUnauthorized)
		return
	}

	// Check if the request is multipart/form-data for file uploads
	contentType := c.Request.Header.Get("Content-Type")
	isMultipart := strings.HasPrefix(contentType, "multipart/form-data")

	var newMerch MerchandiseRequest
	if isMultipart {
		// Handle file uploads
		var err error

		// Merch item details
		newMerch.Name = c.PostForm("name")
		newMerch.Description = c.PostForm("description")
		newMerch.Price, err = strconv.ParseFloat(c.PostForm("price"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format"})
			return
		}

		// Merch item image file
		imageFile, _ := c.FormFile("image")
		if imageFile != nil {
			newMerch.ImageURL, err = utils.UploadToS3(imageFile, "Merchandise", imageFile.Filename)
			if err != nil {
				fmt.Println("Error uploading merch image file:", err)
				c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
				return
			}
		} else {
			newMerch.ImageURL = c.PostForm("image_url")
		}

		// Process sizes from form data
		// Each size is submitted as "sizes=S,10"
		sizesValues := c.PostFormArray("sizes")

		newMerch.Sizes = []SizeInfo{}
		for _, sizeValue := range sizesValues {
			parts := strings.Split(sizeValue, ",")
			if len(parts) != 2 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid size format. Expected: Size,Quantity"})
				return
			}

			size := parts[0]
			quantity, err := strconv.Atoi(parts[1])
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity value. Must be a number"})
				return
			}

			newMerch.Sizes = append(newMerch.Sizes, SizeInfo{
				Size:     size,
				Quantity: quantity,
			})
		}

	} else {
		// Handle JSON requests
		if err := c.ShouldBindJSON(&newMerch); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Validate request
		if newMerch.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
			return
		}

		if newMerch.Price <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be greater than 0"})
			return
		}
	}

	// Create merchandise item
	merch := schema.Merchandise{
		ID:          uuid.New(),
		Name:        newMerch.Name,
		Description: newMerch.Description,
		Price:       newMerch.Price,
		ImageURL:    newMerch.ImageURL,
	}

	// Begin transaction
	ctx := context.Background()
	tx, err := schema.GetDBConn().BeginTx(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	// Ensure rollback - undoes entire db transaction if an error occurs
	defer tx.Rollback()

	// Insert merchandise
	_, err = tx.NewInsert().Model(&merch).Exec(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert merchandise"})
		return
	}

	// Insert sizes
	for _, sizeInfo := range newMerch.Sizes {
		size := schema.MerchandiseSize{
			ID:            uuid.New(),
			MerchandiseID: merch.ID,
			Size:          sizeInfo.Size,
			Quantity:      sizeInfo.Quantity,
		}

		_, err = tx.NewInsert().Model(&size).Exec(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert size"})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Merchandise added successfully",
		"id":      merch.ID,
	})
}

/*
Gets all merchandise items in the database with their associated sizes

	curl -X GET http://localhost:8080/api/merch/all
*/
func GetAllMerchandise(c *gin.Context) {
	type MerchandiseWithSizes struct {
		schema.Merchandise
		Sizes []schema.MerchandiseSize `json:"sizes"`
	}

	var merchandise []schema.Merchandise
	db := schema.GetDBConn()
	ctx := context.Background()

	// Get all merchandise items
	err := db.NewSelect().
		Model(&merchandise).
		Scan(ctx)

	if err != nil {
		fmt.Printf("Error fetching merchandise: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	result := make([]MerchandiseWithSizes, 0, len(merchandise))

	// For each merchandise item, get its sizes
	for _, merch := range merchandise {
		var sizes []schema.MerchandiseSize

		// Query sizes for this specific merchandise item
		err := db.NewSelect().
			Model(&sizes).
			Where("merchandise_id = ?", merch.ID).
			Scan(ctx)

		if err != nil {
			fmt.Printf("Error fetching sizes for merchandise ID %s: %v", merch.ID, err)
			c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
			return
		}

		// Add merchandise with its sizes to result
		result = append(result, MerchandiseWithSizes{
			Merchandise: merch,
			Sizes:       sizes,
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

/*
Deletes merch item from database along with all associated sizes

	curl -X DELETE http://localhost:8080/api/merch/00000000-0000-0000-0000-000000000000 \
	-H "Authorization: Bearer YOUR API KEY"
*/
func DeleteMerchandise(c *gin.Context) {
	if !internal.CheckAuthorization(c) {
		c.AbortWithError(http.StatusUnauthorized, internal.ErrUnauthorized)
		return
	}

	// Ensure merch_id is provided and is a valid UUID
	param := c.Param("merch_id")
	if param == "" {
		fmt.Println("merch_id path parameter is required")
		c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}

	merchID, err := uuid.Parse(param)
	if err != nil {
		fmt.Println("merch_id must be a valid UUID")
		c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}

	// Start a transaction
	ctx := context.Background()
	tx, err := schema.GetDBConn().BeginTx(ctx, nil)
	if err != nil {
		fmt.Printf("Error starting transaction: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}
	defer tx.Rollback()

	// Check if merchandise exists
	var merch schema.Merchandise
	exists, err := tx.NewSelect().
		Model(&merch).
		Where("id = ?", merchID).
		Exists(ctx)

	if err != nil {
		fmt.Printf("Error checking if merchandise exists: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	if !exists {
		fmt.Println("Merch not found")
		c.AbortWithError(http.StatusNotFound, internal.ErrNotFound)
		return
	}

	// Delete the associated sizes first
	_, err = tx.NewDelete().
		Model((*schema.MerchandiseSize)(nil)).
		Where("merchandise_id = ?", merchID).
		Exec(ctx)

	if err != nil {
		fmt.Printf("Error deleting associated sizes: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	// Delete the merchandise item
	_, err = tx.NewDelete().
		Model((*schema.Merchandise)(nil)).
		Where("id = ?", merchID).
		Exec(ctx)

	if err != nil {
		fmt.Printf("Error deleting merchandise: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		fmt.Printf("Error committing transaction: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Merchandise item and associated sizes deleted successfully",
	})
}

type SizeUpdateInfo struct {
	Size     string `json:"size"`
	Quantity *int   `json:"quantity"` // Pointer to allow nil (no update)
}

type MerchandiseUpdateRequest struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       float64          `json:"price"`
	ImageURL    string           `json:"image_url"`
	SizeUpdates []SizeUpdateInfo `json:"sizes"` // Array of size updates
}

/*
Updates an existing merchandise item and its associated sizes

	curl -X PUT http://localhost:8080/api/merch/00000000-0000-0000-0000-000000000000 \
		-H "Authorization: Bearer YOUR API KEY" \
		-H "Content-Type: application/json" \
		-d '{
			"name": "Updated Movie T-Shirt",
			"description": "Updated description",
			"price": 19.99,
			"image_url": "https://example.com/images/updated-tshirt.jpg",
			"sizes": [
				{
					"size": "S",
					"quantity": 0
				},
				{
					"size": "M",
					"quantity": 25
				},
				{
					"size": "L",
					"quantity": 15
				}
			]
		}'

	For file upload submissions:

	curl -X PUT http://localhost:8080/api/merch/00000000-0000-0000-0000-000000000000 \
		-H "Authorization: Bearer YOUR API KEY" \
		-F "name=Updated Movie T-Shirt" \
		-F "description=Updated description" \
		-F "price=19.99" \
		-F "image=@/path/to/updated-image.jpg" \
		-F "sizes=S,0" \
		-F "sizes=M,25" \
		-F "sizes=L,15"
*/
func UpdateMerchandise(c *gin.Context) {
	if !internal.CheckAuthorization(c) {
		c.AbortWithError(http.StatusUnauthorized, internal.ErrUnauthorized)
		return
	}

	// Ensure merch_id is provided and is a valid UUID
	param := c.Param("merch_id")
	if param == "" {
		fmt.Println("merch_id path parameter is required")
		c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}

	merchID, err := uuid.Parse(param)
	if err != nil {
		fmt.Println("merch_id must be a valid UUID")
		c.AbortWithError(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}

	// Check if the request is multipart/form-data for file uploads
	contentType := c.Request.Header.Get("Content-Type")
	isMultipart := strings.HasPrefix(contentType, "multipart/form-data")

	var updateReq MerchandiseUpdateRequest
	if isMultipart {
		// Handle file uploads
		var err error

		// Merch item details
		updateReq.Name = c.PostForm("name")
		updateReq.Description = c.PostForm("description")

		if priceStr := c.PostForm("price"); priceStr != "" {
			updateReq.Price, err = strconv.ParseFloat(priceStr, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format"})
				return
			}
		}

		// Merch item image file
		imageFile, _ := c.FormFile("image")
		if imageFile != nil {
			updateReq.ImageURL, err = utils.UploadToS3(imageFile, "Merchandise", c.PostForm("name"))
			if err != nil {
				fmt.Println("Error uploading merch image file:", err)
				c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
				return
			}
		} else if imageURL := c.PostForm("image_url"); imageURL != "" {
			updateReq.ImageURL = imageURL
		}

		// Process sizes from form data
		sizesValues := c.PostFormArray("sizes")
		for _, sizeValue := range sizesValues {
			parts := strings.Split(sizeValue, ",")
			if len(parts) != 2 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid size format. Expected: Size,Quantity"})
				return
			}

			size := parts[0]
			quantity, err := strconv.Atoi(parts[1])
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity value. Must be a number"})
				return
			}

			updateReq.SizeUpdates = append(updateReq.SizeUpdates, SizeUpdateInfo{
				Size:     size,
				Quantity: &quantity,
			})
		}
	} else {
		// Handle JSON requests
		if err := c.ShouldBindJSON(&updateReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}
	}

	// Start a transaction
	ctx := context.Background()
	tx, err := schema.GetDBConn().BeginTx(ctx, nil)
	if err != nil {
		fmt.Printf("Error starting transaction: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}
	defer tx.Rollback()

	// First check if merchandise exists
	var existingMerch schema.Merchandise
	err = tx.NewSelect().
		Model(&existingMerch).
		Where("id = ?", merchID).
		Scan(ctx)

	if err != nil {
		fmt.Printf("Error checking if merchandise exists: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	// Update merchandise fields (only update non-empty fields)
	updates := make(map[string]any)

	if updateReq.Name != "" {
		updates["name"] = updateReq.Name
	}

	if updateReq.Description != "" {
		updates["description"] = updateReq.Description
	}

	if updateReq.Price > 0 {
		updates["price"] = updateReq.Price
	}

	if updateReq.ImageURL != "" {
		updates["image_url"] = updateReq.ImageURL
	}

	// Only update if there are changes
	if len(updates) > 0 {
		// First get the existing merchandise
		merch := new(schema.Merchandise)
		err = tx.NewSelect().
			Model(merch).
			Where("id = ?", merchID).
			Scan(ctx)
		if err != nil {
			fmt.Printf("Error finding merchandise: %v", err)
			c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
			return
		}

		// Apply updates
		if price, ok := updates["price"].(float64); ok {
			merch.Price = price
		}
		if name, ok := updates["name"].(string); ok {
			merch.Name = name
		}
		if desc, ok := updates["description"].(string); ok {
			merch.Description = desc
		}
		if imgURL, ok := updates["image_url"].(string); ok {
			merch.ImageURL = imgURL
		}

		// Save the updates
		_, err = tx.NewUpdate().
			Model(merch).
			WherePK().
			Exec(ctx)

		if err != nil {
			fmt.Printf("Error updating merchandise: %v", err)
			c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
			return
		}
	}

	// Handle size updates
	for _, sizeUpdate := range updateReq.SizeUpdates {
		if sizeUpdate.Quantity == nil {
			continue // Skip if no quantity update
		}

		// Try to update existing size
		result, err := tx.NewUpdate().
			Model((*schema.MerchandiseSize)(nil)).
			Set("quantity = ?", *sizeUpdate.Quantity).
			Where("merchandise_id = ? AND size = ?", merchID, sizeUpdate.Size).
			Exec(ctx)

		if err != nil {
			fmt.Printf("Error updating size: %v", err)
			c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
			return
		}

		// If no rows were updated, insert new size
		if affected, _ := result.RowsAffected(); affected == 0 {
			size := schema.MerchandiseSize{
				ID:            uuid.New(),
				MerchandiseID: merchID,
				Size:          sizeUpdate.Size,
				Quantity:      *sizeUpdate.Quantity,
			}

			_, err = tx.NewInsert().
				Model(&size).
				Exec(ctx)

			if err != nil {
				fmt.Printf("Error inserting new size: %v", err)
				c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
				return
			}
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		fmt.Printf("Error committing transaction: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Merchandise item updated successfully",
	})
}
