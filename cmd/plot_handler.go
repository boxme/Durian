package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
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

			if !ph.checkIsPostRequest(r.Method, w.Header()) {
				ph.logger.MethodNotAllowed(w, r)
				return
			}

			if err := r.ParseForm(); err != nil {
				ph.logger.BadRequestResponse(w, r, err)
			}

			plotId := getPlotId(r)
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

func (ph *PlotHandler) handleEditPlot() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ph.logger.InfoLog.Printf("Handling edit plot request")

			if !ph.checkIsPostRequest(r.Method, w.Header()) {
				ph.logger.MethodNotAllowed(w, r)
				return
			}

			if err := r.ParseForm(); err != nil {
				ph.logger.BadRequestResponse(w, r, err)
			}

			plotId := getPlotId(r)
			info := r.PostForm.Get("info")

			ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
			defer cancel()

			plot, err := ph.service.EditPlot(ctx, plotId, info)
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
			ph.logger.InfoLog.Printf("handling edit plot is successful")
		},
	)
}

func (ph *PlotHandler) checkIsPostRequest(method string, header http.Header) bool {
	if method != http.MethodPost {
		ph.logger.ErrorLog.Printf("using wrong restful method in plot handler %s", method)

		header.Set("Allow", http.MethodPost)
		return false
	}

	return true
}

func getPlotId(r *http.Request) string {
	plotId := r.PostForm.Get("plot_id")
	if len(plotId) > 0 {
		return plotId
	}

	guid := xid.New()
	return guid.String()
}
