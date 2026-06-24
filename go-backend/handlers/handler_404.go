package handlers

import "net/http"

type NotFoundData struct {
	Page string
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	render(w, "404.html", NotFoundData{Page: "404"})
}
