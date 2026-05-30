package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func (app *application) createStaffAccountHandler(w http.ResponseWriter, r *http.Request) {
	var staff Staff

	if err := json.NewDecoder(r.Body).Decode(&staff); err != nil {
		http.Error(w, "invalid data", http.StatusBadRequest)
		return
	}

	staff.ID = "STF-" + strconv.Itoa(rand.Intn(9000)+1000)

	if staff.Username == "" || staff.Password == "" || staff.Role == "" || staff.Name == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(staff.Password), bcrypt.DefaultCost)
	if err != nil {
		app.logger.Printf("BCRYPT ERROR: encryption failure: %v", err)
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	staff.Password = string(hashedPassword)

	_, err = app.db.Exec("INSERT INTO staff (id, username, password_hash, role) VALUES (?, ?, ?, ?)",
		staff.ID, staff.Username, staff.Password, staff.Role)
	if err != nil {
		app.logger.Printf("DATABASE ERROR: failed to insert staff account: %v", err)
		http.Error(w, "DB Error", 500)
		return
	}

	app.logger.Printf("STAFF REGISTERED: ID=%s, Username=%s, Role=%s", staff.ID, staff.Username, staff.Role)

	staff.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(staff)
}

func (app *application) getAllStaffHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := app.db.Query("SELECT id, username, role FROM staff")
	if err != nil {
		http.Error(w, "DB Error", 500)
		return
	}
	defer rows.Close()
	app.logger.Println("FETCH: retrieving complete staff directory")

	var staffList []Staff
	for rows.Next() {
		var staff Staff
		if err := rows.Scan(&staff.ID, &staff.Username, &staff.Role); err != nil {
			http.Error(w, "DB Error", 500)
			return
		}
		staffList = append(staffList, staff)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(staffList)
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Staff
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "invalid data", http.StatusBadRequest)
		return
	}

	var storedHash string
	err := app.db.QueryRow("SELECT password_hash FROM staff WHERE username = ?", creds.Username).Scan(&storedHash)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(creds.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	fmt.Printf("LOGIN SUCCESS: Username=%s\n", creds.Username)

	creds.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(creds)
}
