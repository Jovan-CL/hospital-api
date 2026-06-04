package main

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"time"
)

type AuthUser struct {
	ID   string
	Role Role
}

type contextKey string

const authUserKey = contextKey("authUser")

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

		var userID AuthUser
		var expiry time.Time

		query := `
			SELECT s.staff_id, u.role, s.expiry 
			FROM sessions s
			JOIN users u ON s.staff_id = u.id
			WHERE s.token = ?`

		err := app.db.QueryRow(query, token).Scan(&userID.ID, &userID.Role, &expiry)
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

		ctx := context.WithValue(r.Context(), authUserKey, userID)
		r = r.WithContext(ctx)

		app.logger.Printf("AUTH SUCCESS: User ID %s authorized", userID.ID)
		next.ServeHTTP(w, r)

	}
}
