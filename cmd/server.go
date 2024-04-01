package main

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Mux http.Handler
	Db  *pgxpool.Pool
}

func CreateServer(
	ctx context.Context,
	db *pgxpool.Pool,
	logger *Logger) (*Server, error) {

	mux := http.NewServeMux()
	addRoutes(mux, logger)

	var handler http.Handler = mux

	return &Server{
		Mux: handler,
		Db:  db,
	}, nil
}

func addRoutes(mux *http.ServeMux, logger *Logger) {
	mux.Handle("/", http.NotFoundHandler())
}
