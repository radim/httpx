# httpx

HTTP helpers for Go web applications with HTMX support.

## Install

```bash
go get github.com/radim/httpx
```

## Setup

```go
import "github.com/radim/httpx"

func main() {
    httpx.InitFlashStore([]byte("your-32-byte-secret-key-here!!!"))
}
```

## Error Handling

```go
func handler(w http.ResponseWriter, r *http.Request) error {
    return httpx.NotFoundError("resource not found")
    return httpx.BadRequestError("invalid input: %v", err)
    return httpx.UnauthorizedError("not authorized")
    return httpx.ForbiddenError("access denied")
    return httpx.ConflictError("already exists")
    return httpx.UnprocessableEntityError("validation failed")
    return httpx.StatusError(418, "I'm a teapot")
}
```

## HTMX

```go
// Check request type
if httpx.IsHTMX(r) { ... }
if httpx.IsBoosted(r) { ... }

// Render full page or partial based on request
httpx.Render(w, r, fullPage, partial)

// Render partial only
httpx.Partial(w, r, component)

// Response headers
httpx.HTMXRedirect(w, "/new-location")
httpx.HTMXRefresh(w)
httpx.HTMXRetarget(w, "#element")
httpx.HTMXReswap(w, "outerHTML")
httpx.HTMXTrigger(w, "itemAdded")
httpx.HTMXPushURL(w, "/new-url")
```

## Flash Messages

```go
// Set
httpx.FlashSuccess(w, r, "Saved!")
httpx.FlashError(w, r, "Something went wrong")
httpx.FlashInfo(w, r, "FYI...")
httpx.FlashWarning(w, r, "Watch out!")

// Get (clears them)
flashes, _ := httpx.GetFlashes(w, r)
for _, f := range flashes {
    // f.Type: FlashTypeSuccess, FlashTypeError, etc.
    // f.Message: the message
}
```

## Request Binding

```go
type Request struct {
    Email string `form:"email" validate:"required,email"`
    Name  string `form:"name" validate:"required"`
}

func handler(w http.ResponseWriter, r *http.Request) error {
    var req Request
    if err := httpx.Bind(r, &req); err != nil {
        return err
    }

    errs := httpx.Validate(req)
    if !errs.Empty() {
        if errs.Has("Email") {
            // errs.Get("Email") == "email" (the failed tag)
        }
        return httpx.Partial(w, r, form(req, errs))
    }
}

// Specific binding
httpx.BindForm(r, &req)
httpx.BindJSON(r, &req)
```

## URL Parameters

```go
id := httpx.Param(r, "id")     // from path (Go 1.22+)
page := httpx.Query(r, "page") // from query string
```

## Example

```go
func LaunchLab(w http.ResponseWriter, r *http.Request) error {
    labID := httpx.Param(r, "lab_id")

    var req struct {
        Email string `form:"email" validate:"required,email"`
    }
    if err := httpx.Bind(r, &req); err != nil {
        return err
    }

    errs := httpx.Validate(req)
    if !errs.Empty() {
        httpx.FlashError(w, r, "Invalid email")
        return httpx.Partial(w, r, labs.Form(req, errs))
    }

    // ... business logic ...

    httpx.FlashSuccess(w, r, "Check your email!")
    httpx.HTMXRedirect(w, "/labs/"+labID)
    return nil
}
```
