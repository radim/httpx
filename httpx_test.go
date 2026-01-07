package httpx

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOnError_HandlerError(t *testing.T) {
	var capturedErr error
	adapter := &HandlerAdapter{
		OnError: func(r *http.Request, err error) {
			capturedErr = err
		},
		InternalErrs: func(w http.ResponseWriter, r *http.Request, err error) {
			w.WriteHeader(http.StatusInternalServerError)
		},
	}

	expectedErr := errors.New("handler error")
	handler := adapter.Handle(func(w http.ResponseWriter, r *http.Request) error {
		return expectedErr
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if capturedErr != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, capturedErr)
	}
}

func TestOnError_Panic(t *testing.T) {
	var capturedErr error
	adapter := &HandlerAdapter{
		OnError: func(r *http.Request, err error) {
			capturedErr = err
		},
		InternalErrs: func(w http.ResponseWriter, r *http.Request, err error) {
			w.WriteHeader(http.StatusInternalServerError)
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("panic error")
	})

	wrapped := RecoverMiddleware(adapter, handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	wrapped.ServeHTTP(w, req)

	if capturedErr == nil {
		t.Fatal("expected captured error, got nil")
	}

	if capturedErr.Error() != "panic error" {
		t.Errorf("expected error string 'panic error', got %v", capturedErr)
	}
}
