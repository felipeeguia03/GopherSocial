package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	max_bytes := 1_048_528
	http.MaxBytesReader(w, r.Body, int64(max_bytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func WriteJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func JsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}
	return WriteJson(w, status, &envelope{Data: data})
}

func WriteJsonError(w http.ResponseWriter, status int, message string) {
	type envelope struct {
		Error string `json:"error"`
	}
	WriteJson(w, status, &envelope{Error: message})
}
