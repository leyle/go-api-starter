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

func ZeroLogMiddleware(logger zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GenerateReqId()

		// first save id into current ctx
		ctx := r.Context()
		ctx = context.WithValue(ctx, ReqIdContextName, id)
		r = r.WithContext(ctx)

		// setup logger req id field
		l := logger.With().Str(ReqIdContextName, id).Logger()
		lctx := l.WithContext(r.Context())
		r = r.WithContext(lctx)

		// setup response headers req id
		w.Header().Set(ReqIdHeaderName, id)

		next.ServeHTTP(w, r)
	})
}
