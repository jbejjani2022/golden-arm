// movie, menu, seating, reservation

package schema

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID        uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	Title     string    `bun:"title"`
	PosterURL string    `bun:"poster_url"`
	Date      time.Time `bun:"date"`
}

type Menu struct {
	ID      uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	MovieID uuid.UUID `bun:"type:uuid,notnull"`
	MenuURL string    `bun:"menu_url"`

	// Foreign key relation to Movie
	Movie Movie `bun:"rel:belongs-to,join:movie_id=id"`
}

type Reservation struct {
	ID         uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	MovieID    uuid.UUID `bun:"type:uuid,notnull"`
	SeatNumber string    `bun:"seat_number"`
	Name       string    `bun:"name"`
	Email      string    `bun:"email"`

	// Foreign key relation to Movie
	Movie Movie `bun:"rel:belongs-to,join:movie_id=id"`
}
