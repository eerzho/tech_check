package middlewares

import (
	"context"
	"net/http"
	"tech_check/internal/constants"

	"github.com/google/uuid"
)

type RequestID struct {
}

func NewRequestId() *RequestID {
	return &RequestID{}
}

func (req *RequestID) Mwr(next http.Handler) http.Handler {
	return req.MwrFunc(next.ServeHTTP)
}

func (req *RequestID) MwrFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		w.Header().Set(constants.HeaderRequestID.String(), reqID)
		ctx := context.WithValue(r.Context(), constants.HeaderRequestID, reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
