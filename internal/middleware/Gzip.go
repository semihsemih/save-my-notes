package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func GzipMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		wrw := NewWrappedResponseWriter(w)
		wrw.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(wrw, r)
		defer wrw.Flush()
		return
	}

	next.ServeHTTP(w, r)
}

type WrappedResponseWriter struct {
	w  http.ResponseWriter
	gw *gzip.Writer
}

func NewWrappedResponseWriter(w http.ResponseWriter) *WrappedResponseWriter {
	gw := gzip.NewWriter(w)
	return &WrappedResponseWriter{w: w, gw: gw}
}

func (wr *WrappedResponseWriter) Header() http.Header {
	return wr.w.Header()
}

func (wr *WrappedResponseWriter) Write(d []byte) (int, error) {
	return wr.gw.Write(d)
}

func (wr *WrappedResponseWriter) WriteHeader(statuscode int)  {
	wr.w.WriteHeader(statuscode)
}

func (wr *WrappedResponseWriter) Flush() {
	wr.gw.Flush()
	wr.gw.Close()
}