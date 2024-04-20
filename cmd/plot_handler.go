package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PlotHandler struct {
	Db      *pgxpool.Pool
	service PlotService
	logger  *Logger
}

func CreatePlotHandler(db *pgxpool.Pool, logger *Logger) PlotHandler {
	plotService := CreatePlotService(db)
	return PlotHandler{
		Db:      db,
		service: plotService,
		logger:  logger,
	}
}

func (ph *PlotHandler) handleAddPlot() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ph.logger.InfoLog.Printf("Handling add plot request")

			if r.Method != http.MethodPost {
				ph.logger.ErrorLog.Printf("using wrong restful method in add plot %s", r.Method)

				w.Header().Set("Allow", http.MethodPost)
				ph.logger.MethodNotAllowed(w, r)
				return
			}

			if err := r.ParseForm(); err != nil {
				ph.logger.BadRequestResponse(w, r, err)
			}

			plotId := r.PostForm.Get("plot_id")
			info := r.PostForm.Get("info")

			ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
			defer cancel()

			plot, err := ph.service.AddPlot(ctx, plotId, info)
			if err != nil {
				ph.logger.BadRequestResponse(w, r, err)
			}

			jsonResponse, err := json.Marshal(plot)
			if err != nil {
				ph.logger.BadRequestResponse(w, r, err)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)
			ph.logger.InfoLog.Printf("handling adding plot is successful")
		},
	)
}
