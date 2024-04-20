package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthHandler struct {
	Db      *pgxpool.Pool
	logger  *Logger
	service AuthService
}

func CreateAuthHandler(db *pgxpool.Pool, logger *Logger) AuthHandler {
	authService := CreateAuthService(db)
	return AuthHandler{
		Db:      db,
		logger:  logger,
		service: authService,
	}
}

func (ah *AuthHandler) handleLogin() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ah.logger.InfoLog.Printf("Handling login")

			if r.Method != http.MethodPost {
				ah.logger.ErrorLog.Printf("using wrong restful method in user login %s", r.Method)

				w.Header().Set("Allow", http.MethodPost)
				ah.logger.MethodNotAllowed(w, r)
				return
			}

			if err := r.ParseForm(); err != nil {
				ah.logger.BadRequestResponse(w, r, err)
				return
			}

			mobileNumber := r.PostForm.Get("mobile_number")
			countryCode, err := strconv.Atoi(r.PostForm.Get("country_code"))
			if err != nil {
				ah.logger.BadRequestResponse(w, r, err)
				return
			}

			password := r.PostForm.Get("password")

			ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
			defer cancel()

			jwtToken, err := ah.service.Login(ctx, countryCode, mobileNumber, password)
			if err != nil {
				ah.logger.BadRequestResponse(w, r, err)
				return
			}

			jsonResponse, err := json.Marshal(jwtToken)
			if err != nil {
				ah.logger.BadRequestResponse(w, r, err)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)
			ah.logger.InfoLog.Printf("handling user login successful")
		},
	)
}
