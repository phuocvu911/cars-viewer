package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

var tmpl, _ = template.ParseFiles("./templates/index.html", "./templates/home.html", "./templates/navfooter.html", "./templates/car.html")

func CarDetailsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are supported.", http.StatusMethodNotAllowed)
		return
	}

	c, _ := r.Cookie("id")

	fmt.Println(c)

	car_id := r.URL.Path[len(CARS_ENDPOINT):]

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

	w.Header().Add("Set-Cookie", "id="+GenerateCookie(5)+"; Max-Age=2592000")

	// Execute template
	err := tmpl.ExecuteTemplate(w, "index.html", car)

	if err != nil {
		http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
		return
	}

}
