package handlers

import (
	"net/http"
)

type HomeData struct {
	Page, Tittle string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are supported.", http.StatusMethodNotAllowed)
		return
	}
	data := HomeData{Page: "home", Tittle: "Home"}
	if err := render(w, "home.html", data); err != nil {
		http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
	}
}
