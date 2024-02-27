package router

import (
	"net/http"
	"time"

	"github.com/YaNeAndrey/ya-metrics/internal/server/handlers"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

func InitRouter(st *storage.StorageRepo) http.Handler {
	r := chi.NewRouter()
	r.NotFound(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
	})

	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	r.Use(LoggerMiddleware(log))

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

func LoggerMiddleware(logger logrus.FieldLogger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			timeStart := time.Now()
			defer func() {
				fields := logrus.Fields{
					//request fields
					"URI":      r.RequestURI,
					"method":   r.Method,
					"duration": time.Since(timeStart),

					//response fields
					"status_code":   ww.Status(),
					"bytes_written": ww.BytesWritten(),
				}
				logger.WithFields(fields).Infoln("New request")
			}()
			h.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
