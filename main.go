package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"nimblestack/database"
	"nimblestack/router"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func initializeSchema(db *sql.DB) error {
	schemaSQL, err := os.ReadFile("sqlc/schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema file: %v", err)
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	queries := strings.SplitSeq(string(schemaSQL), ";")
	for query := range queries {
		if _, err := tx.Exec(query); err != nil {
			log.Printf("Error executing schema query: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit schema: %v", err)
	}

	log.Println("Schema initialized successfully")
	return nil
}

func main() {
	db, err := sql.Open("sqlite3", "sqlite.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	if err := initializeSchema(db); err != nil {
		log.Fatalf("Error initializing schema: %v", err)
	}

	queries := database.New(db)
	jwtSecret := os.Getenv("API_TOKEN")
	route := router.NewRouter(queries, []byte(jwtSecret))

	log.Println("NimbleStack server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", route.Handler()))
}
