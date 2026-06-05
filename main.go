package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Create the table automatically for patients
	query := `CREATE TABLE IF NOT EXISTS patients (
        id TEXT PRIMARY KEY,
        name TEXT,
        age INTEGER,
        condition TEXT
		);`
	_, _ = db.ExecContext(ctx, query)
	if err != nil {
		logger.Fatal("Failed to deploy sessions database schema on patients table: ", err)
	}

	// Create the table automatically for staff
	staffquery := `CREATE TABLE IF NOT EXISTS staff (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE, -- UNIQUE prevents duplicate usernames
    password_hash TEXT,
    role TEXT
	);`
	_, err = db.ExecContext(ctx, staffquery)
	if err != nil {
		logger.Fatal("Failed to deploy sessions database schema on staff table: ", err)
	}

	sessionsTableQuery := `CREATE TABLE IF NOT EXISTS sessions (
	token TEXT PRIMARY KEY,
	staff_id TEXT,
	expiry DATETIME,
	FOREIGN KEY (staff_id) REFERENCES staff(id)
	);`

	_, err = db.ExecContext(ctx, sessionsTableQuery)
	if err != nil {
		logger.Fatal("Failed to deploy sessions database schema on sessions table: ", err)
	}

	app := &application{db: db, logger: logger}

	if err != nil {
		logger.Fatal("Error initializing Logger: ", err)
	}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      app.routes(),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdownChan := make(chan error)

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
		s := <-exit

		logger.Printf("RECEIVED SHUTDOWN SIGNAL: signal='%s'. CLEANING UP CONNECTIONS...", s.String())
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownChan <- srv.Shutdown(ctx)
	}()

	log.Println("Hospital API starting on :8080...")

	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("Server error: ", err)
	}

	err = <-shutdownChan
	if err != nil {
		logger.Fatal("Error occurred while shutting down server: ", err)
	}
}
