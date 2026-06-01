package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

type application struct {
	db     *sql.DB
	logger *log.Logger
}

func main() {
	db, err := sql.Open("sqlite", "./hospital.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the table automatically for patients
	query := `CREATE TABLE IF NOT EXISTS patients (
        id TEXT PRIMARY KEY,
        name TEXT,
        age INTEGER,
        condition TEXT
		);`
	_, _ = db.Exec(query)

	// Create the table automatically for staff
	staffquery := `CREATE TABLE IF NOT EXISTS staff (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE, -- UNIQUE prevents duplicate usernames
    password_hash TEXT,
    role TEXT
	);`
	_, _ = db.Exec(staffquery)

	sessionsTableQuery := `CREATE TABLE IF NOT EXISTS sessions (
	token TEXT PRIMARY KEY,
	staff_id TEXT,
	expiry DATETIME,
	FOREIGN KEY (staff_id) REFERENCES staff(id)
	);`

	_, _ = db.Exec(sessionsTableQuery)

	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{db: db, logger: logger}

	if err != nil {
		logger.Fatal("Error initializing Logger: ", err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.routes(),
	}

	log.Println("Hospital API starting on :8080...")
	log.Fatal(srv.ListenAndServe())
}
