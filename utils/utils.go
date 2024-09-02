package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	_, err = w.Write(jsonBytes) // Write the JSON bytes to the response
	return err                  // Return any potential write error
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

var Validate = validator.New()

func ValidateDecimalPlaces(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(float64)
	if !ok {
		return false
	}

	truncated := math.Trunc(value * 100)
	return value*100 == truncated
}
