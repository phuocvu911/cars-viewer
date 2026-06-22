package handlers

import (
	"cars-viewer/cookies"
	"net/http"
)

// Used just to fetch new set of allowed cookies.
func AllowedCookiesHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are supported.", http.StatusMethodNotAllowed)
		return
	}

	cookies.SetTrackingAllowanceAndLongtermIDCookies(w, true)
}

// Used just to fetch new set of allowed cookies.
func NotAllowedCookiesHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are supported.", http.StatusMethodNotAllowed)
		return
	}

	cookies.SetTrackingAllowanceAndLongtermIDCookies(w, false)
}
