package middleware

import (
	"mime"
	"net/http"
)

// EnforceJSON checks if content-type is application/json
func EnforceJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if contentType != "" {
			mediatype, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
				return
			}

			if mediatype != "application/json" {
				http.Error(
					w, "Content-Type header must be application/json",
					http.StatusUnsupportedMediaType,
				)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
