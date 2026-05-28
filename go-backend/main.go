package main

import (
	"cars-viewer/handlers"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /car/", handlers.CarDetailsHandler)

	// Proxy image requests to localhost:3000
	remoteURL, _ := url.Parse(handlers.API_BASE_URL)
	proxy := httputil.NewSingleHostReverseProxy(remoteURL)

	mux.HandleFunc("/api/images/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	http.ListenAndServe(":8080", mux)
}
