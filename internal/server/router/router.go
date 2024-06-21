package router

import (
	"github.com/YaNeAndrey/ya-metrics/internal/server/config"
	"github.com/YaNeAndrey/ya-metrics/internal/server/handlers"
	"github.com/YaNeAndrey/ya-metrics/internal/server/middleware"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

func InitRouter(c config.Config, st *storage.StorageRepo) http.Handler {
	r := chi.NewRouter()
	r.NotFound(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
	})

	logger := log.New()
	logger.SetLevel(log.InfoLevel)
	r.Use(middleware.MyLoggerMiddleware(logger))
	if c.ServerPrivKey() != nil {
		r.Use(middleware.DecryptMiddleware(c.ServerPrivKey()))
	}
	r.Use(middleware.GzipMiddleware())

	if c.EncryptionKey() != nil {
		r.Use(middleware.SignatureVerificationMiddleware(c.EncryptionKey()))
		r.Use(middleware.SignatureDateMiddleware(c.EncryptionKey()))
	}

	r.Route("/", func(r chi.Router) {
		r.Mount("/debug", chimiddleware.Profiler())
		r.Post("/updates/", func(rw http.ResponseWriter, req *http.Request) {
			handlers.HandlePostUpdateMultipleMetricsJSON(rw, req, st)
		})

		r.Get("/", func(rw http.ResponseWriter, req *http.Request) {
			handlers.HandleGetReadMetrics(rw, req, st)
		})

		r.Get("/ping", func(rw http.ResponseWriter, req *http.Request) {
			handlers.HandleGetPing(rw, req, c)
		})

		r.Route("/value", func(r chi.Router) {
			r.Post("/", func(rw http.ResponseWriter, req *http.Request) {
				handlers.HandlePostReadOneMetricJSON(rw, req, st)
			})
			r.Get("/{type}/{name}", func(rw http.ResponseWriter, req *http.Request) {
				handlers.HandleGetReadOneMetric(rw, req, st)
			})
		})

		r.Route("/update", func(r chi.Router) {
			if c.StoreInterval() == 0 {
				r.Use(middleware.SyncUpdateAndFileStorageMiddleware(c, st))
			}
			r.Post("/", func(rw http.ResponseWriter, req *http.Request) {
				handlers.HandlePostUpdateOneMetricJSON(rw, req, st)
			})
			r.Post("/{type}/{name}/{value}", func(rw http.ResponseWriter, req *http.Request) {
				handlers.HandlePostUpdateOneMetric(rw, req, st)
			})
		})
	})
	return r
}
