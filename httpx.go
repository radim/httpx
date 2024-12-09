package httpx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type (
	HTTPHandlerExt func(http.ResponseWriter, *http.Request) error
	AdapterFunc    func(http.ResponseWriter, *http.Request, error)

	HandlerAdapter struct {
		InternalErrs AdapterFunc
		ClientErrs   AdapterFunc

		UnauthorizedErr AdapterFunc
	}

	Error interface {
		error
		GetStatusCode() int
	}

	AppError struct {
		Err        error
		StatusCode int
	}

	ErrorReporter interface {
		ReportError(ctx context.Context, err error)
	}

	Renderer interface {
		Render500(ctx context.Context, w http.ResponseWriter, errInfo *ErrorInfo)
		RenderAppError(ctx context.Context, w http.ResponseWriter, appErr AppError)
	}

	AppConfig interface {
		IsDevelopment() bool
		ErrorReporter() ErrorReporter
		Renderer() Renderer
	}

	ErrorInfo struct {
		Message string `json:"message,omitempty"`
		Cause   string `json:"cause,omitempty"`
		Stack   string `json:"stack,omitempty"`
	}
)

func BadRequestError(content string, params ...interface{}) AppError {
	return StatusError(http.StatusBadRequest, content, params...)
}

func UnauthorizedError(content string, params ...interface{}) AppError {
	return StatusError(http.StatusUnauthorized, content, params...)
}

func StatusError(statusCode int, content string, params ...interface{}) AppError {
	return AppError{
		Err:        fmt.Errorf(content, params...),
		StatusCode: statusCode,
	}
}

func (e AppError) Error() string {
	return e.Err.Error()
}

func (e AppError) GetStatusCode() int {
	return e.StatusCode
}

func defaultAppError(w http.ResponseWriter, req *http.Request, err error) {
	if e, ok := err.(AppError); ok {
		http.Error(w, err.Error(), e.GetStatusCode())
		return
	}

	http.Error(w, "Bad Request", http.StatusBadRequest)
}

func defaultInternalError(w http.ResponseWriter, req *http.Request, err error) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func InternalErrorsHandler(config AppConfig) func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		var errInfo *ErrorInfo

		w.WriteHeader(http.StatusInternalServerError)

		config.ErrorReporter().ReportError(req.Context(), err)

		if config.IsDevelopment() {
			errInfo = &ErrorInfo{
				Message: fmt.Sprintf("%s", err),
			}

			// Unwrap the error to get the root cause, if any
			if cause := errors.Unwrap(err); cause != nil {
				errInfo.Cause = cause.Error()
			}

			// Check if the error has a stack trace
			type stackTracer interface {
				StackTrace() string
			}

			if stackErr, ok := err.(stackTracer); ok {
				errInfo.Stack = stackErr.StackTrace()
			}
		}

		// Use the Renderer to render the 500 error response
		config.Renderer().Render500(context.Background(), w, errInfo)
	}
}

func NewDefaultHandlerAdapter(config AppConfig) *HandlerAdapter {
	return &HandlerAdapter{
		InternalErrs:    InternalErrorsHandler(config),
		ClientErrs:      defaultAppError,
		UnauthorizedErr: nil,
	}
}

func RecoverMiddleware(adapter *HandlerAdapter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				var err error
				switch x := rec.(type) {
				case string:
					err = fmt.Errorf(x)
				case error:
					err = x
				default:
					err = fmt.Errorf("unknown panic")
				}
				adapter.InternalErrs(w, r, err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (a *HandlerAdapter) Handle(h HTTPHandlerExt) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := h(w, req); err != nil {
			switch e := err.(type) {
			case AppError:
				if e.StatusCode == http.StatusUnauthorized && a.UnauthorizedErr != nil {
					a.UnauthorizedErr(w, req, e)
					return
				}

				if a.ClientErrs != nil {
					a.ClientErrs(w, req, err)
					return
				}

				// Use default AppError handler if no ClientErrs adapter is provided
				defaultAppError(w, req, err)

			default:
				if a.InternalErrs != nil {
					a.InternalErrs(w, req, err)
					return
				}

				// Use default internal error handler if no InternalErrs adapter is provided
				defaultInternalError(w, req, err)
			}
		}
	}
}
