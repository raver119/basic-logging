package sample

import (
	"log"
	"net/http"

	"github.com/raver119/logging/logging"
)

func Middleware(r http.Request, w http.ResponseWriter, next http.Handler) {
	logger := logging.NewBasicLogger()
	emitter := logger.Emitter()
	defer func() {
		// this emitter will actually pump out the logs we've accumulated during this request
		if err := emitter(); err != nil {
			log.Printf("failed to emit logs: %v", err)
		}
	}()

	// add some shared fields
	logger.Add("request_id", r.Header.Get("X-Request-Id"))
	logger.Add("user_agent", r.Header.Get("User-Agent"))
	logger.Add("remote_addr", r.RemoteAddr)

	// package it to new context
	rc := r.WithContext(logging.ToContext(r.Context(), logger))

	// pass it to the next handler/mandleware.
	next.ServeHTTP(w, rc)
}
