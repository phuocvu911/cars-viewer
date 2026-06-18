package handlers

import (
	"cars-viewer/cookies"
	"log"
	"net/http"
	"sync"
)

func CarDetailsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are supported.", http.StatusMethodNotAllowed)
		return
	}

	car_id := r.URL.Path[len(LOCAL_CARS_ROUTE):]

	if len(car_id) == 0 || len(car_id) > 10 {
		http.Error(w, "Bad request.", http.StatusBadRequest)
		return
	}

	var car Car

	errChannel := make(chan error, 1)

	FetchCarByID(car_id, errChannel, &car)

	if err := <-errChannel; err != nil {
		http.Error(w, "Failed to fetch backend data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	car.DataPerID.ImgSrc = IMG_PATH_PREFIX + car.DataPerID.ImgSrc
	car.Page = "gallery"

	cookieCtx, problem := r.Context().Value(cookies.CookieCtxKey{}).(cookies.CookieCtx)

	if !problem {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var wg sync.WaitGroup
	errChannel = make(chan error, 2)

	wg.Go(func() {
		FetchCarCategory(errChannel, &car)
	})
	wg.Go(func() {
		FetchCarManufacturer(errChannel, &car)
	})

	// Wait for both to return
	wg.Wait()

	// Check for errors
	if err := <-errChannel; err != nil {
		http.Error(w, "Failed to fetch car related data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(car)
	AddTrackingItem(&cookieCtx, &car)

	render(w, "car.html", car)
}
