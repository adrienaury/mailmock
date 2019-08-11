// Package httpd exposes the REST API of Mailmock
//
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
	"strconv"

	"github.com/adrienaury/mailmock/internal/repository"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// Routes returns a configures router for REST API serving
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,          // Log API request calls
		middleware.DefaultCompress, // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing server
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

func getAll(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, repository.All()) // A chi router helper for serializing and returning json
}
