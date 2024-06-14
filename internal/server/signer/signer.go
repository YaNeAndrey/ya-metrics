package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
)

type SignerWriter struct {
	w   http.ResponseWriter
	key []byte
}

func NewSignerWriter(w http.ResponseWriter, key []byte) *SignerWriter {
	return &SignerWriter{
		w:   w,
		key: key,
	}
}

func (c *SignerWriter) Header() http.Header {
	return c.w.Header()
}

func (c *SignerWriter) Write(p []byte) (int, error) {
	hashSHA256 := generateSignature(c.key, p)
	c.w.Header().Set("HashSHA256", base64.URLEncoding.EncodeToString(hashSHA256))
	lenBuf, err := c.w.Write(p)

	if err != nil {
		return 0, err
	}
	return lenBuf, nil
}

func (c *SignerWriter) WriteHeader(statusCode int) {
	c.w.WriteHeader(statusCode)
}

func generateSignature(key []byte, date []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(date)
	return h.Sum(nil)
}
