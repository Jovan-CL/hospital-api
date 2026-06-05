package main

import (
	"context"
	"database/sql"
	"fmt"
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

		var user AuthUser
		var expiry time.Time

		fmt.Println("🔍 MIDDLEWARE LOOKING FOR TOKEN:", token)

		query := `
			SELECT s.staff_id, u.role, s.expiry 
			FROM sessions s
			JOIN staff u ON s.staff_id = u.id
			WHERE s.token = ?`

		err := app.db.QueryRow(query, token).Scan(&user.ID, &user.Role, &expiry)
		if err == sql.ErrNoRows {
			app.logger.Printf("AUTH ERROR: invalid session token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		} else if err != nil {
			app.logger.Printf("AUTH ERROR: database error details: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if time.Now().After(expiry) {
			app.logger.Printf("AUTH ERROR: session token expired")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), authUserKey, user)
		r = r.WithContext(ctx)

		app.logger.Printf("AUTH SUCCESS: User ID %s authorized", user.ID)
		next.ServeHTTP(w, r)

	}
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
