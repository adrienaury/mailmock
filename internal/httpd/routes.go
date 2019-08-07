package httpd

import (
	"net/http"

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
	render.JSON(w, r, "TODO") // A chi router helper for serializing and returning json
}

func getAll(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "TODO") // A chi router helper for serializing and returning json
}
