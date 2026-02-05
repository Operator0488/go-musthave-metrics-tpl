package handler

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func compressGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if !supportsGzip {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			errorBadRequest(w, "gzip не инициализировался")
			return
		}
		defer gz.Close()

		rw := compressWriter{
			ResponseWriter: w,
			Writer:         gz,
		}

		w.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(rw, r)
	})
}

func (rw compressWriter) Write(b []byte) (int, error) {
	n, err := rw.Writer.Write(b)
	return n, err
}

func decompressGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentEncoding := r.Header.Get("Content-Encoding")
		supportsGzip := strings.Contains(contentEncoding, "gzip")
		if !supportsGzip {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			errorBadRequest(w, "gzip не распакован")
			return
		}

		rw := compressReader{
			r:  r.Body,
			zr: gz,
		}

		r.Body = &rw
		defer rw.Close()

		next.ServeHTTP(w, r)
	})
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
