package main

import (
	"cars-viewer/handlers"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// Initialize the store with all car models and categories
	if err := handlers.InitStore(); err != nil {
		log.Fatal("Failed to initialize store: " + err.Error())
	}

	//page routes
	mux := http.NewServeMux()

	//serve css file
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("GET /gallery", handlers.GalleryHandler)
	mux.HandleFunc("GET /car/", handlers.CarDetailsHandler)
	mux.HandleFunc("/compare", handlers.CompareHandler)

	// Proxy image requests to localhost:3000
	remoteURL, _ := url.Parse(handlers.API_BASE_URL)
	proxy := httputil.NewSingleHostReverseProxy(remoteURL)
	mux.HandleFunc("/api/images/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
	log.Println("AutoVault ready to see at http://localhost:8080/")
	http.ListenAndServe(":8080", mux)
}
