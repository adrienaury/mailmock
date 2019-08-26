// Copyright (C) 2019  Adrien Aury
//
// This file is part of Mailmock.
//
// Mailmock is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Mailmock is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Mailmock.  If not, see <https://www.gnu.org/licenses/>.

package httpd

import (
	"net/http"
	"time"

	"github.com/adrienaury/mailmock/internal/log"
	"github.com/go-chi/chi/middleware"
)

type chilogger struct {
	logger log.Logger
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

		fields := log.Fields{"took": latency, "status": status, "remote": remote, "request": request, "method": method, "bytes": bytes}
		if requestID != "" {
			fields["request-id"] = requestID
		}

		switch {
		case status < 400:
			c.logger.Info("Request completed", fields)
		case status >= 400 && status < 500:
			c.logger.Warn("Request completed", fields)
		default:
			c.logger.Error("Request completed", fields)
		}

	}
	return http.HandlerFunc(fn)
}
