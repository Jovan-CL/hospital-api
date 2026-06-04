package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /patients", app.createPatientHandler)
	mux.HandleFunc("GET /patients", app.requireAuthentication(app.getAllPatientsHandler))        // Protected
	mux.HandleFunc("GET /patients/{id}", app.requireAuthentication(app.getPatientHandler))       // Protected
	mux.HandleFunc("PATCH /patients/{id}", app.requireAuthentication(app.updatePatientHandler))  // Protected
	mux.HandleFunc("DELETE /patients/{id}", app.requireAuthentication(app.deletePatientHandler)) // Protected

	mux.HandleFunc("GET /stats", app.requireAuthentication(app.countPatientsHandler)) // Protected

	// STAFF routes
	mux.HandleFunc("POST /staff/create-account", app.createStaffAccountHandler)
	mux.HandleFunc("POST /login", app.loginHandler)
	mux.HandleFunc("GET /staff", app.requireAuthentication(app.getAllStaffHandler)) // Protected

	return mux
}
