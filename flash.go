package httpx

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type (
	FlashType string

	Flash struct {
		Type    FlashType
		Message string
	}
)

const (
	FlashTypeSuccess FlashType = "success"
	FlashTypeError   FlashType = "error"
	FlashTypeInfo    FlashType = "info"
	FlashTypeWarning FlashType = "warning"
)

var (
	FlashStore *sessions.CookieStore
	flashTypes = []FlashType{FlashTypeSuccess, FlashTypeError, FlashTypeInfo, FlashTypeWarning}
)

// InitFlashStore initializes the flash message store with secret key
func InitFlashStore(secret []byte) {
	FlashStore = sessions.NewCookieStore(secret)
}

// SetFlash adds a flash message
func SetFlash(w http.ResponseWriter, r *http.Request, ftype FlashType, msg string) error {
	session, err := FlashStore.Get(r, "flash")
	if err != nil {
		return err
	}
	session.AddFlash(msg, string(ftype))
	return session.Save(r, w)
}

// GetFlashes retrieves and clears all flash messages
func GetFlashes(w http.ResponseWriter, r *http.Request) ([]Flash, error) {
	session, err := FlashStore.Get(r, "flash")
	if err != nil {
		return nil, err
	}

	var flashes []Flash
	for _, ftype := range flashTypes {
		for _, msg := range session.Flashes(string(ftype)) {
			if s, ok := msg.(string); ok {
				flashes = append(flashes, Flash{Type: ftype, Message: s})
			}
		}
	}

	if err := session.Save(r, w); err != nil {
		return flashes, err
	}
	return flashes, nil
}

// FlashSuccess adds a success flash message
func FlashSuccess(w http.ResponseWriter, r *http.Request, msg string) error {
	return SetFlash(w, r, FlashTypeSuccess, msg)
}

// FlashError adds an error flash message
func FlashError(w http.ResponseWriter, r *http.Request, msg string) error {
	return SetFlash(w, r, FlashTypeError, msg)
}

// FlashInfo adds an info flash message
func FlashInfo(w http.ResponseWriter, r *http.Request, msg string) error {
	return SetFlash(w, r, FlashTypeInfo, msg)
}

// FlashWarning adds a warning flash message
func FlashWarning(w http.ResponseWriter, r *http.Request, msg string) error {
	return SetFlash(w, r, FlashTypeWarning, msg)
}
