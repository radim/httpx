package httpx

import (
	"net/http"

	"github.com/a-h/templ"
)

// HTMX headers
const (
	// Request headers
	HXRequest    = "HX-Request"
	HXBoosted    = "HX-Boosted"
	HXTrigger    = "HX-Trigger"
	HXTarget     = "HX-Target"
	HXCurrentURL = "HX-Current-URL"

	// Response headers
	HXRedirect      = "HX-Redirect"
	HXRefresh       = "HX-Refresh"
	HXRetarget      = "HX-Retarget"
	HXReswap        = "HX-Reswap"
	HXTriggerHeader = "HX-Trigger"
	HXPushURL       = "HX-Push-Url"
	HXReplaceURL    = "HX-Replace-Url"
	HXReselect      = "HX-Reselect"
	HXTriggerAfter  = "HX-Trigger-After-Swap"
	HXTriggerSettle = "HX-Trigger-After-Settle"
)

// IsHTMX returns true if request came from HTMX
func IsHTMX(r *http.Request) bool {
	return r.Header.Get(HXRequest) == "true"
}

// IsBoosted returns true if request is from boosted link
func IsBoosted(r *http.Request) bool {
	return r.Header.Get(HXBoosted) == "true"
}

// GetHXTarget returns the target element ID from HTMX request
func GetHXTarget(r *http.Request) string {
	return r.Header.Get(HXTarget)
}

// GetHXTrigger returns the triggering element ID from HTMX request
func GetHXTrigger(r *http.Request) string {
	return r.Header.Get(HXTrigger)
}

// GetHXCurrentURL returns the current URL from HTMX request
func GetHXCurrentURL(r *http.Request) string {
	return r.Header.Get(HXCurrentURL)
}

// Render renders full page for normal requests, partial for HTMX
func Render(w http.ResponseWriter, r *http.Request, full, partial templ.Component) error {
	if IsHTMX(r) && !IsBoosted(r) {
		return partial.Render(r.Context(), w)
	}
	return full.Render(r.Context(), w)
}

// Partial renders a component for HTMX responses
func Partial(w http.ResponseWriter, r *http.Request, c templ.Component) error {
	return c.Render(r.Context(), w)
}

// HTMXRedirect sends HX-Redirect header for client-side redirect
func HTMXRedirect(w http.ResponseWriter, url string) {
	w.Header().Set(HXRedirect, url)
}

// HTMXRefresh triggers full page refresh
func HTMXRefresh(w http.ResponseWriter) {
	w.Header().Set(HXRefresh, "true")
}

// HTMXRetarget changes target element for response
func HTMXRetarget(w http.ResponseWriter, selector string) {
	w.Header().Set(HXRetarget, selector)
}

// HTMXReswap changes swap method (innerHTML, outerHTML, beforebegin, etc.)
func HTMXReswap(w http.ResponseWriter, method string) {
	w.Header().Set(HXReswap, method)
}

// HTMXTrigger triggers client-side event
func HTMXTrigger(w http.ResponseWriter, event string) {
	w.Header().Set(HXTriggerHeader, event)
}

// HTMXTriggerAfterSwap triggers event after swap completes
func HTMXTriggerAfterSwap(w http.ResponseWriter, event string) {
	w.Header().Set(HXTriggerAfter, event)
}

// HTMXTriggerAfterSettle triggers event after settle completes
func HTMXTriggerAfterSettle(w http.ResponseWriter, event string) {
	w.Header().Set(HXTriggerSettle, event)
}

// HTMXPushURL pushes URL to browser history
func HTMXPushURL(w http.ResponseWriter, url string) {
	w.Header().Set(HXPushURL, url)
}

// HTMXReplaceURL replaces current URL in browser history
func HTMXReplaceURL(w http.ResponseWriter, url string) {
	w.Header().Set(HXReplaceURL, url)
}

// HTMXReselect changes CSS selector for content to swap
func HTMXReselect(w http.ResponseWriter, selector string) {
	w.Header().Set(HXReselect, selector)
}
