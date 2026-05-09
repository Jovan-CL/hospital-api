package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

type application struct {
	db *sql.DB
}

func main() {
	db, err := sql.Open("sqlite", "./hospital.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the table automatically
	query := `CREATE TABLE IF NOT EXISTS patients (
        id TEXT PRIMARY KEY,
        name TEXT,
        age INTEGER,
        condition TEXT
    );`
	_, _ = db.Exec(query)

	app := &application{db: db}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.routes(),
	}

	log.Println("Hospital API starting on :8080...")
	log.Fatal(srv.ListenAndServe())
}
