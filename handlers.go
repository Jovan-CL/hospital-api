package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func (app *application) createPatientHandler(w http.ResponseWriter, r *http.Request) {
	var p Patient

	fmt.Println(p)

	au, ok := r.Context().Value(authUserKey).(AuthUser)
	if !ok {
		app.logger.Printf("AUTH ERROR: failed to retrieve authenticated user from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	app.logger.Printf("STAFF LOG: User %s (%s) is attempting to create a patient", au.ID, au.Role)

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		fmt.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	app.logger.Printf("Received new patient data: %+v", p)

	// Generate a random ID and convert to string
	p.ID = strconv.Itoa(rand.Intn(9000) + 1000)

	fmt.Println(p)

	if p.Name == "" || p.Age == 0 || p.Condition == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO patients (id, name, age, condition) VALUES (?, ?, ?, ?)"
	_, err := app.db.ExecContext(ctx, query, p.ID, p.Name, p.Age, p.Condition)
	if err != nil {
		http.Error(w, "DB Error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (app *application) getAllPatientsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := app.db.QueryContext(ctx, "SELECT * FROM patients")
	if err != nil {
		http.Error(w, "DB Error", 500)
		return
	}
	defer rows.Close()

	app.logger.Printf("Fetched all patients from database")

	patientsSlice := make([]Patient, 0)
	for rows.Next() {
		var p Patient
		if err := rows.Scan(&p.ID, &p.Name, &p.Age, &p.Condition); err != nil {
			http.Error(w, "DB Error", 500)
			return
		}
		patientsSlice = append(patientsSlice, p)
		// fmt.Println("Patients:", patientsSlice)
		// fmt.Println("f:", f)
		// fmt.Println("p:", p)
		// fmt.Println("app.patients:", app.patients)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patientsSlice)
}

func (app *application) getPatientHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var p Patient

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := app.db.QueryRowContext(ctx, "SELECT id, name, age, condition FROM patients WHERE id = ?", id).
		Scan(&p.ID, &p.Name, &p.Age, &p.Condition)

	if err == sql.ErrNoRows {
		http.Error(w, "Patient not found", 404)
		return
	}
	app.logger.Printf("Fetched a patient from database: %+v", p)

	// w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (app *application) deletePatientHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get the ID from the URL
	id := r.PathValue("id")

	app.logger.Printf("Attempting to delete patient with ID: %s", id)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, _ = app.db.ExecContext(ctx, "DELETE FROM patients WHERE id = ?", id)
	w.WriteHeader(http.StatusNoContent)

	// 4. Send a success response
	w.WriteHeader(http.StatusNoContent) // 204 No Content is standard for DELETE
}

func (app *application) updatePatientHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	// 1. Fetch current data from DB
	var p Patient

	query := "SELECT id, name, age, condition FROM patients WHERE id = ?"
	err := app.db.QueryRowContext(ctx, query, id).
		Scan(&p.ID, &p.Name, &p.Age, &p.Condition)

	app.logger.Printf("Fetched patient for update: %+v", p)

	if err == sql.ErrNoRows {
		http.Error(w, "Patient not found", 404)
		return
	}

	// 2. Decode the new data into a temporary map or struct
	// We use a map here to see EXACTLY which fields the user sent
	var input struct {
		Name      *string `json:"name"`
		Age       *int    `json:"age"`
		Condition *string `json:"condition"`
	}

	json.NewDecoder(r.Body).Decode(&input)

	// 3. Only update fields that were actually in the JSON
	if input.Name != nil {
		p.Name = *input.Name
	}
	if input.Age != nil {
		p.Age = *input.Age
	}
	if input.Condition != nil {
		p.Condition = *input.Condition
	}

	// 4. Save the updated version back to DB
	_, err = app.db.ExecContext(ctx, "UPDATE patients SET name = ?, age = ?, condition = ? WHERE id = ?",
		p.Name, p.Age, p.Condition, id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
	w.WriteHeader(http.StatusOK)
}

func (app *application) countPatientsHandler(w http.ResponseWriter, r *http.Request) {
	var totalPatientCount int
	err := app.db.QueryRow("SELECT COUNT(*) FROM patients").Scan(&totalPatientCount)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	app.logger.Printf("Total patient count: %d", totalPatientCount)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"Total patients": totalPatientCount})
}
