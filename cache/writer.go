package cache

import (
	"net/http"
)

type Writer struct {
	writer   http.ResponseWriter
	response response
	resource string
}

// interface implementation check
var (
	_ http.ResponseWriter = (*Writer)(nil)
)

// NewWriter returns the cache writer
func NewWriter(w http.ResponseWriter, r *http.Request) *Writer {
	return &Writer{
		writer:   w,
		resource: MakeResource(r),
		response: response{
			header: http.Header{},
		},
	}
}

// Header return the response header
func (w *Writer) Header() http.Header {
	return w.response.header
}

// WriteHeader writes headers to the response writer
func (w *Writer) WriteHeader(code int) {
	copyHeader(w.response.header, w.writer.Header())
	w.response.code = code
	w.writer.WriteHeader(code)
}

// Write and cache data
func (w *Writer) Write(bytes []byte) (int, error) {
	w.response.body = make([]byte, len(bytes))
	for k, v := range bytes {
		w.response.body[k] = v
	}
	copyHeader(w.Header(), w.writer.Header())
	set(w.resource, &w.response)
	return w.writer.Write(bytes)
}
