package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /create-patient", app.createPatientHandler)
	mux.HandleFunc("GET /patients", app.requireAuthentication(app.getAllPatientsHandler))           // Protected
	mux.HandleFunc("GET /patients/{id}", app.requireAuthentication(app.getPatientHandler))          // Protected
	mux.HandleFunc("PATCH /patients/{id}", app.requireAuthentication(app.updatePatientHandler))     // Protected
	mux.HandleFunc("DELETE /patients/{id}", app.requireAuthentication(app.deletePatientHandler))    // Protected
	mux.HandleFunc("GET /patients/count", app.requireAuthentication(app.countPatientsHandler))      // Protected
	mux.HandleFunc("GET /patients/find-patient", app.requireAuthentication(app.findPatientHandler)) // Protected

	// STAFF routes
	mux.HandleFunc("POST /staff/create-account", app.createStaffAccountHandler)
	mux.HandleFunc("POST /login", app.loginHandler)
	mux.HandleFunc("GET /staff", app.requireAuthentication(app.getAllStaffHandler)) // Protected

	return mux
}
