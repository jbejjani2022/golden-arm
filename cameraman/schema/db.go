package schema

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	// "golden-arm/schema"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var (
	db   *bun.DB
	once sync.Once
)

func GetDBConn() *bun.DB {
	once.Do(func() {
		// Build DSN string
		connstr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"))

		// Open the database connection
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connstr)))
		db = bun.NewDB(sqldb, pgdialect.New())

		// Add query hook for debugging
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))

		// Test the connection
		err := db.Ping()
		if err != nil {
			log.Fatalf("Failed to connect to the database: %v", err)
		}

		log.Println("Successfully connected to the database.")
	})

	return db
}

func CreateTables() {
	db := GetDBConn()

	// Create the Movie table
	_, err := db.NewCreateTable().Model(&Movie{}).IfNotExists().Exec(context.Background())
	if err != nil {
		log.Fatalf("Failed to create Movie table: %v", err)
	}

	// Create the Menu table
	_, err = db.NewCreateTable().Model(&Menu{}).IfNotExists().Exec(context.Background())
	if err != nil {
		log.Fatalf("Failed to create Menu table: %v", err)
	}

	// Create the Reservation table
	_, err = db.NewCreateTable().Model(&Reservation{}).IfNotExists().Exec(context.Background())
	if err != nil {
		log.Fatalf("Failed to create Reservation table: %v", err)
	}

	fmt.Println("✅ Tables created successfully")
}
