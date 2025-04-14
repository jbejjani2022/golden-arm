package routes

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"golden-arm/internal"
	"golden-arm/schema"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Fixed poster price
const PosterPrice = 10.00

type OrderRequest struct {
	Name  string      `json:"name"`
	Email string      `json:"email"`
	Items []OrderItem `json:"items"`
}

type OrderItem struct {
	MerchandiseID *uuid.UUID `json:"merchandise_id,omitempty"`
	MovieID       *uuid.UUID `json:"movie_id,omitempty"`
	Quantity      int        `json:"quantity"`
	Size          string     `json:"size,omitempty"`
}

type OrderResponse struct {
	OrderID uuid.UUID `json:"order_id"`
	Total   float64   `json:"total"`
}

// Order confirmation email
type OrderEmailData struct {
	Order struct {
		Name  string             `json:"name"`
		Email string             `json:"email"`
		Items []schema.OrderItem `json:"items"`
	}
	Response OrderResponse
}

/*
Adds new order

	curl -X POST http://localhost:8080/api/order \
		-H "Content-Type: application/json" -d \
		'{
			"name": "John Doe",
			"email": "johndoe@example.com",
			"items": [
				{
				"merchandise_id": "00000000-0000-0000-0000-000000000000",
				"quantity": 2,
				"size": "M"
				},
				{
				"merchandise_id": "00000000-0000-0000-0000-000000000000",
				"quantity": 1,
				"size": "L"
				},
				{
				"movie_id": "00000000-0000-0000-0000-000000000000",
				"quantity": 1
				},
				{
				"movie_id": "00000000-0000-0000-0000-000000000000",
				"quantity": 2
				}
			]
		}'
*/
func AddOrder(c *gin.Context) {
	var newOrder OrderRequest
	if err := c.ShouldBindJSON(&newOrder); err != nil {
		c.JSON(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}

	if newOrder.Name == "" || newOrder.Email == "" || len(newOrder.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
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

	// Process order items
	total, orderItems, err := processOrderItems(ctx, tx, newOrder.Items)
	if err != nil {
		c.JSON(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}

	// Create order
	orderID := uuid.New()
	order := schema.Order{
		ID:    orderID,
		Name:  newOrder.Name,
		Email: newOrder.Email,
		Date:  time.Now(),
		Total: total,
		Paid:  false,
	}

	_, err = tx.NewInsert().Model(&order).Exec(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	for i := range orderItems {
		orderItems[i].OrderID = orderID
		_, err = tx.NewInsert().Model(&orderItems[i]).Exec(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add order items"})
			return
		}
	}

	// Update inventory
	if err := updateInventory(ctx, tx, orderItems); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Prepare response
	response := OrderResponse{
		OrderID: orderID,
		Total:   total,
	}

	// Prepare email data with schema.OrderItem that includes relationships
	var orderItemsWithRelations []schema.OrderItem
	err = tx.NewSelect().
		Model(&orderItemsWithRelations).
		Relation("Merchandise").
		Relation("Movie").
		Where("order_id = ?", orderID).
		Scan(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items for email"})
		return
	}

	// Send confirmation email
	emailData := OrderEmailData{
		Order: struct {
			Name  string             `json:"name"`
			Email string             `json:"email"`
			Items []schema.OrderItem `json:"items"`
		}{
			Name:  newOrder.Name,
			Email: newOrder.Email,
			Items: orderItemsWithRelations,
		},
		Response: response,
	}

	if err := sendOrderConfirmationEmail(emailData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to send confirmation email: %v", err)})
		return
	}

	// Only commit if email was sent successfully
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, response)
}

// Validates order items and calculates the total cost
func processOrderItems(ctx context.Context, tx bun.Tx, items []OrderItem) (float64, []schema.OrderItem, error) {
	var total float64
	var orderItems []schema.OrderItem

	for _, item := range items {
		if item.Quantity <= 0 {
			return 0, nil, errors.New("quantity must be positive")
		}

		orderItem := schema.OrderItem{
			ID:       uuid.New(),
			Quantity: item.Quantity,
		}

		// Process merchandise item
		if item.MerchandiseID != nil {
			var merch schema.Merchandise
			err := tx.NewSelect().Model(&merch).Where("id = ?", item.MerchandiseID).Scan(ctx)
			if err != nil {
				return 0, nil, errors.New("merchandise not found")
			}

			// Set merchandise-specific fields
			orderItem.MerchandiseID = item.MerchandiseID
			orderItem.MovieID = nil
			orderItem.Price = merch.Price
			orderItem.Size = item.Size

			// Validate size and check inventory if size is specified
			if item.Size != "" {
				var merchSize schema.MerchandiseSize
				err := tx.NewSelect().Model(&merchSize).
					Where("merchandise_id = ? AND size = ?", item.MerchandiseID, item.Size).
					Scan(ctx)

				if err != nil {
					return 0, nil, errors.New("size not available for this merchandise")
				}

				if merchSize.Quantity < item.Quantity {
					return 0, nil, errors.New("insufficient inventory for " + merch.Name + " size " + item.Size)
				}
			}

			total += merch.Price * float64(item.Quantity)
		} else if item.MovieID != nil {
			// Process movie poster item
			var movie schema.Movie
			err := tx.NewSelect().Model(&movie).Where("id = ?", item.MovieID).Scan(ctx)
			if err != nil {
				return 0, nil, errors.New("movie not found")
			}

			// Set movie-specific fields
			orderItem.MovieID = item.MovieID
			orderItem.MerchandiseID = nil
			orderItem.Price = PosterPrice
			orderItem.Size = "" // No size for movie posters

			total += PosterPrice * float64(item.Quantity)
		} else {
			return 0, nil, errors.New("either merchandise_id or movie_id must be provided")
		}

		orderItems = append(orderItems, orderItem)
	}

	return total, orderItems, nil
}

// Reduces the quantity of merchandise in inventory according to the order
func updateInventory(ctx context.Context, tx bun.Tx, items []schema.OrderItem) error {
	for _, item := range items {
		// Only update inventory for merchandise items with a size
		if item.MerchandiseID != nil && item.Size != "" {
			result, err := tx.NewUpdate().Model((*schema.MerchandiseSize)(nil)).
				Set("quantity = quantity - ?", item.Quantity).
				Where("merchandise_id = ? AND size = ? AND quantity >= ?",
					item.MerchandiseID, item.Size, item.Quantity).
				Exec(ctx)

			if err != nil {
				return err
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil || rowsAffected == 0 {
				return errors.New("could not update inventory")
			}
		}
	}

	return nil
}

//go:embed templates/*
var orderEmailTemplate embed.FS

func sendOrderConfirmationEmail(data OrderEmailData) error {
	// Create base template and register functions
	base := template.New("order_email.html")
	base.Funcs(template.FuncMap{
		"mul": func(price float64, quantity int) float64 {
			return price * float64(quantity)
		},
	})

	// Parse the HTML email template
	tmpl, err := base.ParseFS(orderEmailTemplate, "templates/order_email.html")
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	// Configure SMTP settings
	smtpHost := "smtp.gmail.com"
	smtpPort := "587" // Common port for TLS
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// Prepare email message
	from := smtpUsername
	subject := "Confirming your order at The Golden Arm"
	message := fmt.Sprintf("Subject: %s\r\nFrom: %s\r\nTo: %s\r\nCC: %s\r\nContent-Type: text/html\r\n\r\n%s",
		subject, from, data.Order.Email, from, body.String())

	// Connect to the SMTP server
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, append([]string{data.Order.Email}, from), []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Printf("Confirmation email sent to %s\n", data.Order.Email)
	return nil
}

/*
Updates an order's payment status

	curl -X PUT http://localhost:8080/api/order/status/:order_id -H "Authorization: Bearer YOUR API KEY" \
		-H "Content-Type: application/json" \
		-d '{"paid": true}'
*/
func UpdateOrderStatus(c *gin.Context) {
	if !internal.CheckAuthorization(c) {
		c.AbortWithError(http.StatusUnauthorized, internal.ErrUnauthorized)
		return
	}

	// Get order ID from URL parameter
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Parse request body
	var request struct {
		Paid bool `json:"paid"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, internal.ErrBadRequest)
		return
	}

	// Update order in database
	db := schema.GetDBConn()
	ctx := context.Background()

	result, err := db.NewUpdate().
		Model((*schema.Order)(nil)).
		Set("paid = ?", request.Paid).
		Where("id = ?", orderID).
		Exec(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

/*
Deletes an order and restores inventory

	curl -X DELETE http://localhost:8080/api/order/:order_id -H "Authorization: Bearer YOUR API KEY"
*/
func DeleteOrder(c *gin.Context) {
	if !internal.CheckAuthorization(c) {
		c.AbortWithError(http.StatusUnauthorized, internal.ErrUnauthorized)
		return
	}

	// Get order ID from URL parameter
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Begin transaction
	ctx := context.Background()
	tx, err := schema.GetDBConn().BeginTx(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback()

	// Check if the order exists and get its payment status
	var order schema.Order
	err = tx.NewSelect().
		Model(&order).
		Where("id = ?", orderID).
		Scan(ctx)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Fetch the order items
	var orderItems []schema.OrderItem
	err = tx.NewSelect().
		Model(&orderItems).
		Where("order_id = ?", orderID).
		Scan(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items"})
		return
	}

	// Only restore inventory if the order was unpaid
	if !order.Paid {
		if err := restoreInventory(ctx, tx, orderItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Delete order items
	_, err = tx.NewDelete().
		Model((*schema.OrderItem)(nil)).
		Where("order_id = ?", orderID).
		Exec(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order items"})
		return
	}

	// Delete the order
	result, err := tx.NewDelete().
		Model((*schema.Order)(nil)).
		Where("id = ?", orderID).
		Exec(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Restores inventory quantities when an order is deleted
func restoreInventory(ctx context.Context, tx bun.Tx, items []schema.OrderItem) error {
	for _, item := range items {
		// Only update inventory for merchandise items with a size
		if item.MerchandiseID != nil && item.Size != "" {
			result, err := tx.NewUpdate().Model((*schema.MerchandiseSize)(nil)).
				Set("quantity = quantity + ?", item.Quantity).
				Where("merchandise_id = ? AND size = ?",
					item.MerchandiseID, item.Size).
				Exec(ctx)

			if err != nil {
				return err
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil || rowsAffected == 0 {
				return errors.New("could not restore inventory for merchandise ID: " + item.MerchandiseID.String())
			}
		}
	}

	return nil
}

/*
Gets all orders in the database with their items

	curl -X GET http://localhost:8080/api/order/all -H "Authorization: Bearer YOUR API KEY"
*/
func GetAllOrders(c *gin.Context) {
	if !internal.CheckAuthorization(c) {
		c.AbortWithError(http.StatusUnauthorized, internal.ErrUnauthorized)
		return
	}

	type OrderItem struct {
		ID            uuid.UUID       `json:"id"`
		MerchandiseID *uuid.UUID      `json:"merchandise_id,omitempty"`
		MovieID       *uuid.UUID      `json:"movie_id,omitempty"`
		Quantity      int             `json:"quantity"`
		Size          string          `json:"size,omitempty"`
		Price         float64         `json:"price"`
		Merchandise   *schema.Merchandise `json:"merchandise,omitempty"`
		Movie         *schema.Movie      `json:"movie,omitempty"`
	}

	type OrderWithItems struct {
		ID     uuid.UUID    `json:"id"`
		Name   string       `json:"name"`
		Email  string       `json:"email"`
		Date   time.Time    `json:"date"`
		Total  float64      `json:"total"`
		Paid   bool         `json:"paid"`
		Items  []OrderItem  `json:"items"`
	}

	var orders []schema.Order
	db := schema.GetDBConn()
	ctx := context.Background()

	// Fetch all orders from the database
	err := db.NewSelect().
		Model(&orders).
		Scan(ctx)

	if err != nil {
		fmt.Printf("Error fetching orders: %v", err)
		c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
		return
	}

	result := make([]OrderWithItems, 0, len(orders))

	for _, order := range orders {
		var schemaItems []schema.OrderItem
		err := db.NewSelect().
			Model(&schemaItems).
			Relation("Merchandise").
			Relation("Movie").
			Where("order_id = ?", order.ID).
			Scan(ctx)

		if err != nil {
			fmt.Printf("Error fetching order items: %v", err)
			c.AbortWithError(http.StatusInternalServerError, internal.ErrInternalServer)
			return
		}

		// Convert schema items to our simplified response items
		responseItems := make([]OrderItem, len(schemaItems))
		for i, item := range schemaItems {
			responseItems[i] = OrderItem{
				ID:            item.ID,
				MerchandiseID: item.MerchandiseID,
				MovieID:       item.MovieID,
				Quantity:      item.Quantity,
				Size:          item.Size,
				Price:         item.Price,
				Merchandise:   item.Merchandise,
				Movie:         item.Movie,
			}
		}

		result = append(result, OrderWithItems{
			ID:    order.ID,
			Name:  order.Name,
			Email: order.Email,
			Date:  order.Date,
			Total: order.Total,
			Paid:  order.Paid,
			Items: responseItems,
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}
