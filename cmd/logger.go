package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

type Logger struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (l *Logger) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// To show where the error originated from
	l.ErrorLog.Output(2, trace)
	l.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	l.errorResponse(w, r, http.StatusInternalServerError, message)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (l *Logger) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	l.logError(r, err)

	message := "The request received is not valid"
	l.errorResponse(w, r, http.StatusBadRequest, message)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
func (l *Logger) NotFound(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	l.errorResponse(w, r, http.StatusNotFound, message)
}

func (l *Logger) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	l.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// Wrap the handler with a logging handler
func (l *Logger) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// Note that the errors parameter here has the type map[string]string, which is exactly
// the same as the errors map contained in our Validator type.
func (l *Logger) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	l.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *Logger) EditConflictError(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (l *Logger) errorResponse(
	w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := Encode(w, r, status, env)
	if err != nil {
		l.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (l *Logger) logError(r *http.Request, err error) {
	l.ErrorLog.Println(err)
}
