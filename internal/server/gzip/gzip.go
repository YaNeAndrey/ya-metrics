package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"slices"
)

type CompressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *CompressWriter) Header() http.Header {
	return c.w.Header()
}

/*
func (c *CompressWriter) Write(p []byte) (int, error) {
	len, err := c.w.Write(p)
	//c.zw = nil
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return len, err
}*/

func (c *CompressWriter) Write(p []byte) (int, error) {
	//c.body.Write(p)
	contentTypes := c.w.Header().Values("Content-Type")
	if slices.Contains(contentTypes, "application/json") || slices.Contains(contentTypes, "text/html") {
		c.w.Header().Set("Content-Encoding", "gzip")
		return c.zw.Write(p)
	} else {
		c.zw = nil
		return c.w.Write(p)
	}
	//return c.zw.Write(p)
}

func (c *CompressWriter) WriteHeader(statusCode int) {
	c.w.WriteHeader(statusCode)
}

/*
	func (c *CompressWriter) WriteHeader(statusCode int) {
		contentType := c.w.Header().Get("Content-Type")
		log.Println(contentType)
		if statusCode < 300 { //&& (contentType == "application/json" || contentType == "text/html") {
			log.Println("c.w.Header().Set(\"Content-Encoding\", \"gzip\")")
			c.w.Header().Set("Content-Encoding", "gzip")
		}
		c.w.WriteHeader(statusCode)
	}
*/
/*
func (c *CompressWriter) WriteHeader(statusCode int) {
	contentType := c.w.Header().Get("Content-Type")
	if statusCode < 300 && (contentType == "application/json" || contentType == "text/html") {

		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}*/

func (c *CompressWriter) Close() error {
	if c.zw != nil {
		return c.zw.Close()
	}
	return nil
}

type CompressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &CompressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c CompressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *CompressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
