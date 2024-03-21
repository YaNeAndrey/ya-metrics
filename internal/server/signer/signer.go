package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
)

type SigrerWriter struct {
	w   http.ResponseWriter
	key []byte
}

func NewSigrerWriter(w http.ResponseWriter, key []byte) *SigrerWriter {
	return &SigrerWriter{
		w:   w,
		key: key,
	}
}

func (c *SigrerWriter) Header() http.Header {
	return c.w.Header()
}

func (c *SigrerWriter) Write(p []byte) (int, error) {
	hashSHA256 := generateSignature(c.key, p)
	c.w.Header().Set("HashSHA256", base64.URLEncoding.EncodeToString(hashSHA256))
	lenBuf, err := c.Write(p)

	if err != nil {
		return 0, err
	}
	return lenBuf, err
}

func (c *SigrerWriter) WriteHeader(statusCode int) {
	c.w.WriteHeader(statusCode)
}

func generateSignature(key []byte, date []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(date)
	return h.Sum(nil)
}
