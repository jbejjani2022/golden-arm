package schema

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID      uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	Title   string    `bun:"title,notnull"`
	Date    time.Time `bun:"date,notnull,unique"` // Date of movie screening
	Runtime int       `bun:"runtime,notnull"`     // Movie runtime in minutes
	// Public URLs to images stored in AWS S3
	PosterURL string `bun:"poster_url"`
	MenuURL   string `bun:"menu_url"`
}

type Reservation struct {
	ID         uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	MovieID    uuid.UUID `bun:"type:uuid,notnull"`
	SeatNumber string    `bun:"seat_number,notnull"` // e.g. A1, A2, ...
	Date       time.Time `bun:"date,notnull"`        // When the reservation was made
	// Movie-goer information
	Name  string `bun:"name,notnull"`
	Email string `bun:"email,notnull"`

	// Foreign key relation to Movie
	Movie Movie `bun:"rel:belongs-to,join:movie_id=id"`
}

// Feedback from movie-goers; e.g. suggestion for future screening
type Comment struct {
	ID      uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	Name    string    `bun:"name,notnull"`
	Email   string    `bun:"email,notnull"`
	Comment string    `bun:"comment,notnull"`
	Date    time.Time `bun:"date,notnull"`
}

// Calendar with upcoming screenings; an image with an associated date range
type Calendar struct {
	ID        uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	StartDate time.Time `bun:"start_date,notnull"` // Start date of the calendar
	EndDate   time.Time `bun:"end_date,notnull"`   // End date of the calendar
	ImageURL  string    `bun:"image_url,notnull"`  // Public URL to calendar image stored in AWS S3
	Date      time.Time `bun:"date,notnull"`       // Date the calendar was added
}

// A merchandise item available for purchase (e.g. t-shirts)
type Merchandise struct {
	ID          uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	Name        string    `bun:"name,notnull"`
	Description string    `bun:"description"`
	Price       float64   `bun:"price,notnull"`
	ImageURL    string    `bun:"image_url"`
}

// An available size for a merchandise item
type MerchandiseSize struct {
	ID            uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	MerchandiseID uuid.UUID `bun:"type:uuid,notnull"`
	Size          string    `bun:"size,notnull"`               // e.g. "S", "M", "L", "XL"
	Quantity      int       `bun:"quantity,notnull,default:0"` // To track inventory count

	// Foreign key relation
	Merchandise Merchandise `bun:"rel:belongs-to,join:merchandise_id=id"`
}

// An individual item in a customer's order
type OrderItem struct {
	ID            uuid.UUID  `bun:"type:uuid,pk,default:gen_random_uuid()"`
	OrderID       uuid.UUID  `bun:"type:uuid,notnull"`
	MerchandiseID *uuid.UUID `bun:"type:uuid"` // Can be null for movie posters
	MovieID       *uuid.UUID `bun:"type:uuid"` // Can be null for regular merchandise
	Quantity      int        `bun:"quantity,notnull"`
	Size          string     `bun:"size"`          // Can be null for non-apparel items
	Price         float64    `bun:"price,notnull"` // Price at time of purchase

	// Foreign key relations
	Order       Order        `bun:"rel:belongs-to,join:order_id=id"`
	Merchandise *Merchandise `bun:"rel:belongs-to,join:merchandise_id=id"`
	Movie       *Movie       `bun:"rel:belongs-to,join:movie_id=id"`
}

// A customer's order
type Order struct {
	ID    uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	Name  string    `bun:"name,notnull"`
	Email string    `bun:"email,notnull"`
	Date  time.Time `bun:"date,notnull"`
	Total float64   `bun:"total,notnull"`              // Total cost of the order
	Paid  bool      `bun:"paid,notnull,default:false"` // True if order is complete and payment has been received
}
