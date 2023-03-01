package cache

import (
	"net/http"
	"strings"
	"sync"
)

type response struct {
	header http.Header
	code   int
	body   []byte
}

type memCache struct {
	lock sync.RWMutex
	data map[string]response
}

var (
	cache = memCache{data: map[string]response{}}
)

func set(resource string, response *response) {
	cache.lock.Lock()
	if response == nil {
		delete(cache.data, resource)
	} else {
		cache.data[resource] = *response
	}
	cache.lock.Unlock()

}

func get(resource string) *response {
	cache.lock.RLock()
	resp, ok := cache.data[resource]
	cache.lock.Unlock()
	if ok {
		return &resp
	}
	return nil
}

func copyHeader(src, dst http.Header) {
	for key, list := range src {
		for _, value := range list {
			dst.Add(key, value)
		}
	}
}

func MakeResource(r *http.Request) string {
	if r == nil {
		return ""
	}
	return strings.TrimSuffix(r.URL.RequestURI(), "/")
}

func Clean() {
	cache.lock.Lock()
	cache.data = map[string]response{}
	cache.lock.Unlock()
}

func Drop(res string) {
	set(res, nil)
}

// Serve Determine if cached data can or should be served
func Serve(w http.ResponseWriter, r *http.Request) bool {
	// Check for bad requests
	if w == nil || r == nil {
		return false
	}
	// Check for no-cache
	if r.Header.Get("Cache-Control") == "no-cache" {
		return false
	}
	// Check the resources
	resp := get(MakeResource(r))
	if resp == nil {
		return false
	}

	copyHeader(resp.header, w.Header())
	w.WriteHeader(resp.code)

	if r.Method != http.MethodHead {
		_, err := w.Write(resp.body)
		if err != nil {
			return false
		}
	}
	return true
}
