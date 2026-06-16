package handlers

import (
	"net/http"
	"strconv"
)

type EnrichedCarModel struct {
	CarModel
	ManufacturerName    string
	ManufacturerCountry string
	FoundingYear        int
	CategoryName        string
}

type CompareData struct {
	Page                       string
	AllModels                  []CarModel
	Cars                       []EnrichedCarModel
	MaxHP, MaxYear             int
	HasResults, FilterReceived bool
}

func CompareHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	d := CompareData{Page: "compare", AllModels: store.CarModels}
	query := r.URL.Query()
	ids := query["ids"] // multiple values for ?ids=1&ids=2
	maxHP, maxYear := 0, 0

	for _, idStr := range ids {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		for _, em := range derived.EnrichedModels {
			if em.ID == id {
				d.Cars = append(d.Cars, em)
				if em.Specifications.Horsepower > maxHP {
					maxHP = em.Specifications.Horsepower
				}
				if em.Year > maxYear {
					maxYear = em.Year
				}
				break
			}
		}
	}
	d.MaxHP = maxHP
	d.MaxYear = maxYear
	d.HasResults = (len(d.Cars) >= 2 && len(d.Cars) <= 4)
	if !d.HasResults && query.Get("submitted") == "1" {
		d.FilterReceived = true
	}
	render(w, "compare.html", d)
}
