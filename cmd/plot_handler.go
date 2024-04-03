package main

import "github.com/jackc/pgx/v5/pgxpool"

type PlotHandler struct {
	Db      *pgxpool.Pool
	service PlotService
	logger  Logger
}
