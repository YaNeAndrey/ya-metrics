package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

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

			allAcceptEncodingHeaders := strings.Split(r.Header.Values("Accept-Encoding")[0], ", ")
			log.Println(allAcceptEncodingHeaders)
			if slices.Contains(allAcceptEncodingHeaders, "gzip") {
				log.Println("gzip - ok")
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
