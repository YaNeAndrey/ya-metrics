package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"slices"
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
			// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
			// который будем передавать следующей функции
			ow := w

			//ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
			allAcceptEncodingHeaders := r.Header.Values("Accept-Encoding")
			if slices.Contains(allAcceptEncodingHeaders, "gzip") {
				//if content type json or html
				//responseContentType := ww.Header().Values("Content-Type")
				//if slices.Contains(responseContentType, "application/json") || slices.Contains(responseContentType, "text/html") {
				// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
				cw := gzip.NewCompressWriter(w)
				cw.Header().Set("Content-Encoding", "gzip")
				// меняем оригинальный http.ResponseWriter на новый
				ow = cw
				// не забываем отправить клиенту все сжатые данные после завершения middleware
				defer cw.Close()
				//}
			}

			// проверяем, что клиент отправил серверу сжатые данные в формате gzip
			contentEncodings := r.Header.Values("Content-Encoding")
			if slices.Contains(contentEncodings, "gzip") {
				// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
				cr, err := gzip.NewCompressReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				// меняем тело запроса на новое
				r.Body = cr
				defer cr.Close()
			}

			// передаём управление хендлеру
			h.ServeHTTP(ow, r)
		}
		return http.HandlerFunc(fn)
	}
}

/*
func GzipDecompressRequestBodyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentEncodings := r.Header.Values("Content-Encoding")
		if !slices.Contains(contentEncodings, "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer gz.Close()

		body, err := io.ReadAll(gz)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))

		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func GzipCompressResponseBodyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allAcceptEncodingHeaders := r.Header.Values("Accept-Encoding")
		if !slices.Contains(allAcceptEncodingHeaders, "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(w, r)
	})
}
*/
