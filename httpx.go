package httpx

import (
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
		StatusCode() int
	}

	AppError struct {
		Err        error
		StatusCode int
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

func defaultAppError(w http.ResponseWriter, req *http.Request, err error) {
	if e, ok := err.(AppError); ok {
		http.Error(w, err.Error(), e.StatusCode)
		return
	}

	http.Error(w, "Bad Request", http.StatusBadRequest)

}

func defaultInternalError(w http.ResponseWriter, req *http.Request, err error) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

				http.Error(w, e.Error(), e.StatusCode)
			default:
				if a.InternalErrs != nil {
					a.InternalErrs(w, req, err)

					return
				}

				defaultInternalError(w, req, err)
			}
		}
	}
}
