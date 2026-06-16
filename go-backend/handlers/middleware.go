package handlers

import (
	"net/http"
)

func GetRecommendationsMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//----------------------------------------
		// Add logic for recommending files
		//----------------------------------------
		// 1. Get cookie if exists
		// 2. Get permission for tracking (mandatory cookie)
		// 3. fetch data for the user based on the data
		next.ServeHTTP(w, r) // serve the page normally, but inject recommendations if possible
	})
}
