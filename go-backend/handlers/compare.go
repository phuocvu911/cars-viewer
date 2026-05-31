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
	Page, Title    string
	AllModels      []CarModel
	Cars           []EnrichedCarModel
	MaxHP, MaxYear int
	HasResults     bool
}

func CompareHandler(w http.ResponseWriter, r *http.Request) {
	d := CompareData{Page: "compare", Title: "Compare", AllModels: store.CarModels}

	if r.Method == http.MethodPost {
		r.ParseForm()
		ids := r.Form["ids"]
		maxHP, maxYear := 0, 0
		for _, idStr := range ids {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				continue
			}
			for _, m := range d.AllModels {
				if m.ID == id {
					em := enrich(m)
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
		d.HasResults = len(d.Cars) >= 2
	}
	if err := render(w, "compare.html", d); err != nil {
		http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
	}
}
