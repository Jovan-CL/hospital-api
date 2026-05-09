package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /patients", app.createPatientHandler)
	mux.HandleFunc("GET /patients", app.getAllPatientsHandler)
	mux.HandleFunc("GET /patients/{id}", app.getPatientHandler)
	mux.HandleFunc("PATCH /patients/{id}", app.updatePatientHandler)
	mux.HandleFunc("DELETE /patients/{id}", app.deletePatientHandler)

	return mux
}
