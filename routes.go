package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /create-patient", app.createPatientHandler)
	mux.HandleFunc("GET /patients", app.getAllPatientsHandler)
	mux.HandleFunc("GET /patients/{id}", app.getPatientHandler)
	mux.HandleFunc("PATCH /patients/{id}", app.updatePatientHandler)
	mux.HandleFunc("DELETE /patients/{id}", app.deletePatientHandler)
	mux.HandleFunc("GET /patients/count", app.countPatientsHandler)
	mux.HandleFunc("GET /patients/find-patient", app.findPatientHandler)

	// STAFF routes
	mux.HandleFunc("POST /staff/create-account", app.createStaffAccountHandler)
	mux.HandleFunc("POST /login", app.loginHandler)
	mux.HandleFunc("GET /staff", app.getAllStaffHandler)

	return mux
}
