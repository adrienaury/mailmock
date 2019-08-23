package httpd

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/goph/logur"
)

type chilogger struct {
	logger logur.Logger
}

func (c chilogger) middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var requestID string
		if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
			requestID = reqID.(string)
		}

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		latency := time.Since(start)
		status := ww.Status()
		remote := r.RemoteAddr
		request := r.RequestURI
		method := r.Method
		bytes := ww.BytesWritten()

		fields := logur.Fields{"took": latency, "status": status, "remote": remote, "request": request, "method": method, "bytes": bytes}
		if requestID != "" {
			fields["request-id"] = requestID
		}

		c.logger.Info("Request completed", fields)
	}
	return http.HandlerFunc(fn)
}
