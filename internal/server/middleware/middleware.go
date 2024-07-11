package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"github.com/YaNeAndrey/ya-metrics/internal/server/signer"
	"github.com/YaNeAndrey/ya-metrics/internal/server/utils"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/YaNeAndrey/ya-metrics/internal/server/config"
	"github.com/YaNeAndrey/ya-metrics/internal/server/gzip"
)

func MyLoggerMiddleware(logger log.FieldLogger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			timeStart := time.Now()
			defer func() {
				fields := log.Fields{
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

func DecryptMiddleware(key *rsa.PrivateKey) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if r.Method == http.MethodPost {
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				} else {
					decryptBody, decErr := rsa.DecryptPKCS1v15(rand.Reader, key, body)
					if decErr != nil {
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					r.Body = io.NopCloser(bytes.NewReader(decryptBody))
				}
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func SignatureDateMiddleware(key []byte) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sw := signer.NewSignerWriter(w, key)
			h.ServeHTTP(sw, r)
		}
		return http.HandlerFunc(fn)
	}
}

func SignatureVerificationMiddleware(key []byte) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ow := w
			base64HashSHA256 := r.Header.Get("HashSHA256")
			if base64HashSHA256 == "" {
				h.ServeHTTP(ow, r)
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			hash, err := base64.URLEncoding.DecodeString(base64HashSHA256)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if !hmac.Equal(hash, generateSignature(key, body)) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			h.ServeHTTP(ow, r)
		}
		return http.HandlerFunc(fn)
	}
}

func CheckSubnetMiddleware(subnet *net.IPNet) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			agentIPstr := r.Header.Get("X-Real-IP")
			agentIP := net.ParseIP(agentIPstr)
			if agentIP == nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if !subnet.Contains(agentIP) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func generateSignature(key []byte, date []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(date)
	return h.Sum(nil)
}

func SyncUpdateAndFileStorageMiddleware(c config.Config, st *storage.StorageRepo) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)

			err := utils.SaveAllMetricsToFile(c.FileStoragePath(), st)
			if err != nil {
				return
			}
		}
		return http.HandlerFunc(fn)
	}
}
