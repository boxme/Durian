package main

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func addRoutes(mux *http.ServeMux, db *pgxpool.Pool, logger *Logger) {
	mux.Handle("/", http.NotFoundHandler())

	authHandler := CreateAuthHandler(db, logger)
	mux.Handle("/login", authHandler.handleLogin())

	plotHandler := CreatePlotHandler(db, logger)
	mux.Handle("/addPlot/", plotHandler.handleAddPlot())
}
