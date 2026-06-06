package main

import (
	"cars-viewer/handlers"
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// Innitialize templates
	handlers.InitTemplates()

	// Initialize the store with all car models and categories
	if err := handlers.InitStore(); err != nil {
		log.Fatal("Failed to fetch data: " + err.Error())
	}

	// Run the background goroutine to refresh the store every 10 minutes
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handlers.StoreRefresh(ctx)

	mux := http.NewServeMux()

	// Serving css file and hooking up the handlers
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("GET /gallery", handlers.GalleryHandler)
	mux.HandleFunc("GET /car/", handlers.CarDetailsHandler)
	mux.HandleFunc("GET /compare", handlers.CompareHandler)
	mux.HandleFunc("GET /stats", handlers.StatsHandler)

	// Proxy image requests to localhost:3000
	remoteURL, _ := url.Parse(handlers.API_BASE_URL)
	proxy := httputil.NewSingleHostReverseProxy(remoteURL)
	mux.Handle("/api/images/", proxy)

	log.Println("AutoVault ready to see at http://localhost:8080/")
	http.ListenAndServe(":8080", mux)
}
