package log

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"net/http"
)

// based on github.com/rs/zerolog
// add correlation id in the request and response
// provide a http/net middleware

const (
	ReqIdHeaderName  = "X-REQ-ID"
	ReqIdContextName = "reqId"
)

func GenerateReqId() string {
	return uuid.New().String()
}

func ReqIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.Header.Get(ReqIdHeaderName)
		if id == "" {
			id = GenerateReqId()
		}
		ctx = context.WithValue(ctx, ReqIdContextName, id)
		r = r.WithContext(ctx)
		log := zerolog.Ctx(ctx)
		log.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str(ReqIdContextName, id)
		})
		w.Header().Set(ReqIdHeaderName, id)
		next.ServeHTTP(w, r)
	})
}

func LogMiddleware(log zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := log.With().Logger()
			r = r.WithContext(l.WithContext(r.Context()))
			next.ServeHTTP(w, r)
		})
	}
}
