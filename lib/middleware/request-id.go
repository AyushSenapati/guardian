package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type reqIDCtxKey string

const (
	reqIDKey    reqIDCtxKey = "guardianCTX"
	reqIDHeader string      = "X-Request-ID"
)

// WithReqID adds the given request ID to the provided context
func WithReqID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, reqIDKey, reqID)
}

// ReqIDFromCtx retrieves request ID from the givent context
func ReqIDFromCtx(ctx context.Context) string {
	reqID := ctx.Value(reqIDKey)
	if reqID == nil {
		return ""
	}
	return reqID.(string)
}

// RequestID middleware adds X-Reques-ID header to every incoming request,
// so that the request can be traced in the distributed system
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(reqIDHeader)

		// If request ID already exist append new request ID to make it unique
		if reqID != "" {
			reqID += ("," + uuid.New().String())
		} else {
			reqID = uuid.New().String()
		}

		r.Header.Set(reqIDHeader, reqID)
		w.Header().Set(reqIDHeader, reqID)

		// Add request ID to context, so that system can log the error with request ID
		r = r.WithContext(context.WithValue(r.Context(), reqIDKey, reqID))

		next.ServeHTTP(w, r)
	})
}
