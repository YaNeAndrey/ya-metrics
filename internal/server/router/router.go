package router

import (
	"net/http"

	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"github.com/YaNeAndrey/ya-metrics/internal/server/handlers"

	"github.com/go-chi/chi/v5"
)

func InitRouter(ms *storage.MemStorage) http.Handler {
	r := chi.NewRouter()
	r.NotFound(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
			handlers.HandleGetRoot(rw,r,ms)
		})
		
		r.Route("/value", func(r chi.Router) {
			r.Get("/{type}/{name}",func(rw http.ResponseWriter, r *http.Request) {
				handlers.HandleGetMetricValue(rw,r,ms)
			})
		})

		r.Route("/update", func(r chi.Router) {
			r.Post("/{type}/{name}/{value}", func(rw http.ResponseWriter, r *http.Request) {
				handlers.HandlePostUpdateMetricValue(rw,r,ms)
			})
		})
	}) 
	return r
}