package handlers

import "net/http"

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("asset not found\n"))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Running API V1\n"))
}
