// Package api contains HTTP handlers for server API routes.
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/waterfountain1996/cardvalidate"
)

// Application error code.
type apiErrorCode int

const (
	errGeneralError apiErrorCode = iota
	errMalformedNumber
	errUnknownIssuer
	errInvalidAccountNumber
	errMalformedDate
	errCardExpired
)

// apiError represents an HTTP API error returned from handlers.
type apiError struct {
	Code       apiErrorCode `json:"code"`
	Message    string       `json:"message,omitempty"`
	StatusCode int          `json:"-"`
	OrigError  error        `json:"-"`
}

// Error implements error.
func (e *apiError) Error() string {
	return e.Message
}

// ValidationHandler returns a handler that validates credit card information.
func ValidationHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := handleCardValidation(w, r)
		if err == nil {
			return
		}

		var res *apiError
		if !errors.As(err, &res) {
			// If handler returns an unwrapped error we can assume it's an internal server error.
			log.Printf("unhandled error: %s\n", err)
			res = &apiError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal server error",
			}
		}

		renderJSON(w, res.StatusCode, validationResponse{Error: res})
	})
}

// handleCardValidation handles credit card validation for given request.
func handleCardValidation(w http.ResponseWriter, r *http.Request) error {
	ccInfo, err := decodeJSON[creditCardInfo](r)
	if err != nil {
		return &apiError{
			StatusCode: http.StatusBadRequest,
			Code:       errGeneralError,
			Message:    "Invalid JSON request",
		}
	}

	if err := cardvalidate.Validate(ccInfo.CardNumber, ccInfo.ExpirationDate); err != nil {
		e := &apiError{
			StatusCode: http.StatusUnprocessableEntity,
			OrigError:  err,
		}
		switch {
		case errors.Is(err, cardvalidate.ErrMalformedNumber):
			e.Code, e.Message = errMalformedNumber, "Malformed credit card number"
		case errors.Is(err, cardvalidate.ErrUnknownIssuer):
			e.Code, e.Message = errUnknownIssuer, "Unknown IIN"
		case errors.Is(err, cardvalidate.ErrInvalidAccountNumber):
			e.Code, e.Message = errInvalidAccountNumber, "Invalid account number"
		case errors.Is(err, cardvalidate.ErrMalformedDate):
			e.Code, e.Message = errMalformedDate, "Malformed expiration date"
		case errors.Is(err, cardvalidate.ErrCardExpired):
			e.Code, e.Message = errCardExpired, "Credit card has expired"
		}
		return e
	}

	return renderJSON(w, http.StatusOK, validationResponse{Valid: true})
}

// creditCardInfo is a request payload for validation handler.
type creditCardInfo struct {
	CardNumber     string `json:"number"`
	ExpirationDate string `json:"exp_date"`
}

// validationResponse is a response structure for validation handler.
type validationResponse struct {
	Valid bool      `json:"valid"`
	Error *apiError `json:"error,omitempty"`
}

// decodeJSON unmarshals JSON request body into T.
func decodeJSON[T any](r *http.Request) (*T, error) {
	v := new(T)
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return nil, fmt.Errorf("error decoding a JSON request: %w", err)
	}
	return v, nil
}

// renderJSON writes a JSON response with given status code.
func renderJSON(w http.ResponseWriter, statusCode int, o any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(o); err != nil {
		return fmt.Errorf("error encoding a JSON response: %w", err)
	}
	return nil
}
