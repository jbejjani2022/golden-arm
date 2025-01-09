package schema

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID        uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	Title     string    `bun:"title,notnull"`
	Date      time.Time `bun:"date,notnull"` // Date of movie screening
	PosterURL string    `bun:"poster_url"`   // Public URL to movie poster image stored in AWS S3
}

type Menu struct {
	ID      uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	MovieID uuid.UUID `bun:"type:uuid,notnull"`
	MenuURL string    `bun:"menu_url"` // Public URL to menu image stored in AWS S3

	// Foreign key relation to Movie
	Movie Movie `bun:"rel:belongs-to,join:movie_id=id"`
}

type Reservation struct {
	ID         uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	MovieID    uuid.UUID `bun:"type:uuid,notnull"`
	SeatNumber string    `bun:"seat_number,notnull"`
	// Movie-goer information
	Name  string `bun:"name,notnull"`
	Email string `bun:"email,notnull"`

	// Foreign key relation to Movie
	Movie Movie `bun:"rel:belongs-to,join:movie_id=id"`
}
