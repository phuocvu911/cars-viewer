package handlers

import (
	"net/http"
)

type HomeData struct {
	Page string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := HomeData{Page: "home"}
	render(w, "home.html", data)
}
