package router

import (
	"github.com/YaNeAndrey/ya-metrics/internal/server/config"
	"github.com/YaNeAndrey/ya-metrics/internal/server/handlers"
	"github.com/YaNeAndrey/ya-metrics/internal/server/middleware"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func InitRouter(c config.Config, st *storage.StorageRepo) http.Handler {
	r := chi.NewRouter()
	r.NotFound(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
	})

	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	r.Use(middleware.MyLoggerMiddleware(log))
	r.Use(middleware.GzipMiddleware())

	r.Route("/", func(r chi.Router) {
		r.Get("/", func(rw http.ResponseWriter, req *http.Request) {

			handlers.HandleGetRoot(rw, req, st)
		})

		r.Route("/value", func(r chi.Router) {
			r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
				handlers.HandlePostMetricValueJSON(rw, r, st)
			})
			r.Get("/{type}/{name}", func(rw http.ResponseWriter, r *http.Request) {
				handlers.HandleGetMetricValue(rw, r, st)
			})
		})

		r.Route("/update", func(r chi.Router) {
			if c.StoreInterval() == 0 {
				r.Use(middleware.SyncUpdateAndFileStorageMiddleware(c, st))
			}

			r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
				handlers.HandlePostUpdateMetricValueJSON(rw, r, st)
			})
			r.Post("/{type}/{name}/{value}", func(rw http.ResponseWriter, r *http.Request) {
				handlers.HandlePostUpdateMetricValue(rw, r, st)
			})
		})
	})
	return r
}
