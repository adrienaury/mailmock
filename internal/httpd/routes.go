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
	"fmt"
	"net/http"
	"strconv"

	"github.com/adrienaury/mailmock/internal/repository"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// Routes returns a configured router for REST API serving.
func (srv *Server) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.RequestID,                          // Creates a unique request ID
		chilogger{srv.logger}.middleware,              // Log API request calls
		middleware.DefaultCompress,                    // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes,                    // Redirect slashes to no slash URL versions
		middleware.Recoverer,                          // Recover from panics without crashing server
	)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/mailmock", myRoutes())
	})

	return router
}

func myRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{ID}", getOne)
	router.Get("/", getAll)
	return router
}

func getOne(w http.ResponseWriter, r *http.Request) {
	trID := chi.URLParam(r, "ID")
	i, err := strconv.ParseInt(trID, 10, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	obj := repository.Use(int(i))
	if obj == nil {
		http.NotFound(w, r)
		return
	}
	render.JSON(w, r, obj) // A chi router helper for serializing and returning json
}

const maxLimit = 50

func getAll(w http.ResponseWriter, r *http.Request) {
	var from, limit int64
	var err error

	if froms, ok := r.URL.Query()["from"]; !ok || len(froms) < 1 {
		from = 0
	} else {
		from, err = strconv.ParseInt(froms[0], 10, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if limits, ok := r.URL.Query()["limit"]; !ok || len(limits) < 1 {
		limit = 20
	} else {
		limit, err = strconv.ParseInt(limits[0], 10, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if limit > maxLimit {
		http.Error(w, fmt.Sprintf("Maximum allowed limit is %v", maxLimit), http.StatusBadRequest)
		return
	}

	objs, all := repository.All(int(from), int(limit))
	if objs == nil {
		http.NotFound(w, r)
		return
	}

	if !all {
		render.Status(r, http.StatusPartialContent)
	}
	w.Header().Set("Content-Range", fmt.Sprintf("%v-%v/%v", from, from+limit, repository.Len()))
	w.Header().Set("Accept-Range", fmt.Sprintf("%v %v", "mailmock", maxLimit))

	render.JSON(w, r, objs) // A chi router helper for serializing and returning json
}
