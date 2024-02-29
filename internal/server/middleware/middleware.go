package middleware

import (
	"github.com/YaNeAndrey/ya-metrics/internal/server/utils"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/YaNeAndrey/ya-metrics/internal/server/config"
	"github.com/YaNeAndrey/ya-metrics/internal/server/gzip"
)

func MyLoggerMiddleware(logger logrus.FieldLogger) func(h http.Handler) http.Handler {
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

func GzipMiddleware() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ow := w

			//allAcceptEncodingHeaders := strings.Split(r.Header.Values("Accept-Encoding")[0], ", ")
			var allAcceptEncodingSlice []string
			allAcceptEncodingHeaders := r.Header.Values("Accept-Encoding")
			if len(allAcceptEncodingHeaders) > 0 {
				allAcceptEncodingSlice = strings.Split(allAcceptEncodingHeaders[0], ", ")
			}
			if slices.Contains(allAcceptEncodingSlice, "gzip") {
				cw := gzip.NewCompressWriter(w)
				ow = cw
				defer cw.Close()
			}

			contentEncodings := r.Header.Values("Content-Encoding")
			if slices.Contains(contentEncodings, "gzip") {
				cr, err := gzip.NewCompressReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = cr
				defer cr.Close()
			}
			h.ServeHTTP(ow, r)
		}
		return http.HandlerFunc(fn)
	}
}

func SyncUpdateAndFileStorageMiddleware(c config.Config, st *storage.StorageRepo) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)

			err := utils.SaveAllMetricsToFile(c, st)
			if err != nil {
				return
			}
		}
		return http.HandlerFunc(fn)
	}
}
