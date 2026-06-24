package handlers

import (
	"cars-viewer/cookies"
	"cars-viewer/models"
	"net/http"
	"strconv"
	"sync"
)

func CarDetailsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are supported.", http.StatusMethodNotAllowed)
		return
	}

	car_id_str := r.URL.Path[len(LOCAL_CARS_ROUTE):]
	car_id, err := strconv.Atoi(car_id_str)
	if err != nil || car_id < 0 {
		NotFoundHandler(w, r)
		return
	}

	var car models.Car

	errChannel := make(chan error, 1)

	FetchCarByID(car_id_str, errChannel, &car)

	if err := <-errChannel; err != nil {
		http.Error(w, "Failed to fetch backend data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	car.DataPerID.ImgSrc = IMG_PATH_PREFIX + car.DataPerID.ImgSrc
	car.Page = "gallery"

	// -----------------------
	// Use for reading context the previous handler has passed.
	cookieCtx, no_problem := r.Context().Value(cookies.CookieCtxKey{}).(cookies.CookieCtx)

	if !no_problem {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// -----------------------

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

	close(errChannel)

	// Check for errors
	for err := range errChannel {
		if err != nil {
			NotFoundHandler(w, r)
			// http.Error(w, "Failed to fetch car related data: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if cookieCtx.AllowTracking == nil || cookieCtx.AllowTracking.Value == "false" {

		// Check if all cookies needs to be revoked
		if cookieCtx.DeleteAllCookies {
			cookies.SetDeleteAllCookiesHeader(w)
		}

		render(w, "car.html", car)
		return

	}

	if cookieCtx.AllowTracking.Value == "true" {
		AddTrackingItem(&cookieCtx, &car)
		http.SetCookie(w, cookieCtx.ShortCookie)
	}

	if cookieCtx.ReturnAllCookies {
		cookies.WriteLongCookieHeader(w, cookieCtx.LongCookie)
	}

	render(w, "car.html", car)
}
