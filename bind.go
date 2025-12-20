package httpx

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

var (
	formDecoder = form.NewDecoder()
	validate    = validator.New()
)

// ValidationErrors maps field names to error messages
type ValidationErrors map[string]string

// Has returns true if field has an error
func (e ValidationErrors) Has(field string) bool {
	_, ok := e[field]
	return ok
}

// Get returns error message for field
func (e ValidationErrors) Get(field string) string {
	return e[field]
}

// Empty returns true if no validation errors
func (e ValidationErrors) Empty() bool {
	return len(e) == 0
}

// BindForm binds form data to struct
func BindForm(r *http.Request, v any) error {
	if err := r.ParseForm(); err != nil {
		return BadRequestError("invalid form: %v", err)
	}
	if err := formDecoder.Decode(v, r.PostForm); err != nil {
		return BadRequestError("form decode: %v", err)
	}
	return nil
}

// BindJSON binds JSON body to struct
func BindJSON(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return BadRequestError("invalid json: %v", err)
	}
	return nil
}

// Bind auto-detects content type and binds request body
func Bind(r *http.Request, v any) error {
	ct := r.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "application/json") {
		return BindJSON(r, v)
	}
	return BindForm(r, v)
}

// Validate validates struct and returns error map with human-readable messages
func Validate(v any) ValidationErrors {
	errs := make(ValidationErrors)
	if err := validate.Struct(v); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				errs[e.Field()] = TranslateError(e.Tag())
			}
		}
	}
	return errs
}

// Param gets URL parameter (Go 1.22+ native routing)
func Param(r *http.Request, key string) string {
	return r.PathValue(key)
}

// Query gets query parameter
func Query(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
