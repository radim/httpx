package httpx

import (
	"fmt"
	"net/http"
	"reflect"
)

var defaultMessages = map[string]string{
	"required": "is required",
	"email":    "must be a valid email",
	"url":      "must be a valid URL",
	"min":      "is too short",
	"max":      "is too long",
	"oneof":    "must be one of the allowed values",
}

// RegisterMessages adds or overrides error messages
func RegisterMessages(msgs map[string]string) {
	for tag, msg := range msgs {
		defaultMessages[tag] = msg
	}
}

// TranslateError converts validator tag to human message
func TranslateError(tag string) string {
	if msg, ok := defaultMessages[tag]; ok {
		return msg
	}
	return "is invalid"
}

// Form wraps struct with validation errors
type Form[T any] struct {
	Data   T
	Errors ValidationErrors
}

// NewForm creates form wrapper
func NewForm[T any]() *Form[T] {
	return &Form[T]{Errors: make(ValidationErrors)}
}

// Bind binds and validates request data
func (f *Form[T]) Bind(r *http.Request) error {
	if err := Bind(r, &f.Data); err != nil {
		return err
	}
	f.Errors = Validate(f.Data)
	return nil
}

// Valid returns true if no errors
func (f *Form[T]) Valid() bool {
	return f.Errors.Empty()
}

// HasError checks field error
func (f *Form[T]) HasError(field string) bool {
	return f.Errors.Has(field)
}

// Error returns field error message
func (f *Form[T]) Error(field string) string {
	return f.Errors.Get(field)
}

// AddError adds custom error
func (f *Form[T]) AddError(field, message string) {
	f.Errors[field] = message
}

// Value returns field value as string
func (f *Form[T]) Value(field string) string {
	v := reflect.ValueOf(f.Data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return ""
	}
	fv := v.FieldByName(field)
	if !fv.IsValid() {
		return ""
	}
	return fmt.Sprintf("%v", fv.Interface())
}

// ErrorClass returns class if field has error
func (f *Form[T]) ErrorClass(field, class string) string {
	if f.HasError(field) {
		return class
	}
	return ""
}
