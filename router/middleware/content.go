package middleware

import (
	"net/http"
	"path/filepath"
)

// AddContentType is a middleware that adds appropriate Content-Type headers based on file extensions.
// It wraps an http.Handler and sets the Content-Type header before passing the request to the next handler.
//
// Parameters:
//   - next: The http.Handler to be wrapped
//
// Returns:
//   - http.Handler: A new handler that adds Content-Type headers
func AddContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := filepath.Ext(r.URL.Path)
		switch ext {
		case "js":
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case "svg":
			w.Header().Set("Content-Type", "image/svg+xml")
		}
		next.ServeHTTP(w, r)
	})
}
