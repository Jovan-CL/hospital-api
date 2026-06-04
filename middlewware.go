package main

import (
	"database/sql"
	"net/http"
	"strings"
	"time"
)

func (app *application) requireAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.logger.Printf("AUTH ERROR: missing Authorization header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.logger.Printf("AUTH ERROR: invalid Authorization header format")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		var userID string
		var expiry time.Time

		query := "SELECT staff_id, expiry FROM sessions WHERE token = ?"

		err := app.db.QueryRow(query, token).Scan(&userID, &expiry)
		if err == sql.ErrNoRows {
			app.logger.Printf("AUTH ERROR: invalid session token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		} else if err != nil {
			app.logger.Printf("AUTH ERROR: database error")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if time.Now().After(expiry) {
			app.logger.Printf("AUTH ERROR: session token expired")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		app.logger.Printf("AUTH SUCCESS: User ID %s authorized", userID)
		next.ServeHTTP(w, r)

	}
}
