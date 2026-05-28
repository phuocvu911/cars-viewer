package handlers

import (
	"html/template"
	"net/http"
)

var tmpl, _ = template.ParseFiles("./templates/index.html", "./templates/home.html", "./templates/navfooter.html", "./templates/car.html")

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

	car.DataPerID.ImgSrc = IMG_PATH_PREFIX + car.DataPerID.ImgSrc

	err := tmpl.ExecuteTemplate(w, "index.html", car)

	if err != nil {
		http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
		return
	}

}
