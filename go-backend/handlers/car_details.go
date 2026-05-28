package handlers

import (
	"encoding/json"
	"net/http"
)

func CarDetailsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are supported.", http.StatusMethodNotAllowed)
		return
	}

	car_id := r.URL.Path[len("/car/"):]

	if len(car_id) == 0 || len(car_id) > 10 {
		http.Error(w, "Bad request.", http.StatusBadRequest)
		return
	}

	var car Car

	errChannel := make(chan error, 1)

	go FetchCar(car_id, errChannel, &car)

	if err := <-errChannel; err != nil {
		http.Error(w, "Failed to fetch backend data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(car)

	if err != nil {
		http.Error(w, "Failed to encode Car object into JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

}
