package cache

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type mockWriter response

func newMockWriter() *mockWriter {
	return &mockWriter{
		body:   []byte{},
		header: http.Header{},
	}
}
func (mw *mockWriter) Write(b []byte) (int, error) {
	mw.body = make([]byte, len(b))
	for k, v := range b {
		mw.body[k] = v
	}
	return len(b), nil
}

func (mw *mockWriter) WriteHeader(statusCode int) { mw.code = statusCode }
func (mw *mockWriter) Header() http.Header        { return mw.header }

func TestWriter(t *testing.T) {
	mw := newMockWriter()

	res := "/test/url?with=params"
	u, err := url.Parse(res)
	if err != nil {
		t.Fatal("Invalid url")
	}
	req := &http.Request{
		URL: u,
	}
	t.Log("test NewWriter")
	w := NewWriter(mw, req)
	if w.resource != res {
		t.Errorf("Resources are different. Expected %s / Actual: %s", res, w.resource)
	}
	if w.writer != mw {
		t.Fatal("Writer not assigned")
	}

	t.Log("test Header")
	h := w.Header()
	h.Add("test", "value")
	h2 := w.response.header
	if h2.Get("test") != "value" {
		t.Error("Value not stored in header")
	}

	t.Log("test WriteHeader")
	c := 200
	w.WriteHeader(c)
	if mw.code != c {
		t.Errorf("Status code not written. Expected %d / Actual %d", c, mw.code)
	}
	if w.response.code != c {
		t.Errorf("Status code not stored. Expected %d / Actual %d", c, w.response.code)
	}
	h2 = mw.header
	if h2.Get("test") != "value" {
		t.Error("Header not written")
	}

	t.Log("test Write")
	bd := []byte{1, 2, 3, 4, 5}
	n, err := w.Write(bd)
	if err != nil {
		t.Fatalf("Unexpected error while writing: %s", err)
	}
	if n != len(bd) {
		t.Errorf("Unexpected number of bytes written. Expected %d / Actual: %d", len(bd), n)
	}
	if &w.response.body == &bd {
		t.Error("Body assigned, not copied")
	}
	if !reflect.DeepEqual(w.response.body, bd) {
		t.Error("Body no copied")
		t.Error(w.response.body)
		t.Error(bd)
	}
	if &mw.body == &bd {
		t.Error("Body assigned, not copied")
	}
	if !reflect.DeepEqual(mw.body, bd) {
		t.Error("Body not passed copied")
		t.Error(w.response.body)
		t.Error(bd)
	}

}
