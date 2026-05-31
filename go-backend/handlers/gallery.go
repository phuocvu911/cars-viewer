package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

type GalleryData struct {
	Page, Title                      string
	Models                           []EnrichedCarModel
	Categories                       []Category
	Manufacturers                    []Manufacturer
	Drivetrains, Years               []string
	Query, CatF, MfgF, YearF, DriveF string
	ResultCount                      int
}

func GalleryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are supported.", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()
	models := enrichAll()

	catF := q.Get("category")
	mfgF := q.Get("manufacturer")
	yearF := q.Get("year")
	driveF := q.Get("drivetrain")
	search := q.Get("q")

	filtered := make([]EnrichedCarModel, 0, len(models))
	for _, m := range models {
		if catF != "" {
			id, _ := strconv.Atoi(catF)
			if m.CategoryID != id {
				continue
			}
		}
		if mfgF != "" {
			id, _ := strconv.Atoi(mfgF)
			if m.ManufacturerID != id {
				continue
			}
		}
		if yearF != "" {
			y, _ := strconv.Atoi(yearF)
			if m.Year != y {
				continue
			}
		}
		if driveF != "" && !strings.EqualFold(m.Specifications.Drivetrain, driveF) {
			continue
		}
		if search != "" {
			s := strings.ToLower(search)
			if !strings.Contains(strings.ToLower(m.Name), s) &&
				!strings.Contains(strings.ToLower(m.ManufacturerName), s) &&
				!strings.Contains(strings.ToLower(m.CategoryName), s) &&
				!strings.Contains(strings.ToLower(m.ManufacturerCountry), s) &&
				!strings.Contains(strings.ToLower(m.Specifications.Engine), s) {
				continue
			}
		}
		filtered = append(filtered, m)
	}

	// Collect unique values for filter dropdowns
	yearSet := map[string]bool{}
	driveSet := map[string]bool{}
	for _, m := range models {
		yearSet[strconv.Itoa(m.Year)] = true
		driveSet[m.Specifications.Drivetrain] = true
	}
	years := make([]string, 0)
	for y := range yearSet {
		years = append(years, y)
	}
	drives := make([]string, 0)
	for d := range driveSet {
		drives = append(drives, d)
	}

	data := GalleryData{
		Page: "gallery", Title: "Gallery",
		Models: filtered, Categories: store.Categories, Manufacturers: store.Manufacturers,
		Drivetrains: drives, Years: years,
		Query: search, CatF: catF, MfgF: mfgF, YearF: yearF, DriveF: driveF,
		ResultCount: len(filtered),
	}
	if err := render(w, "gallery.html", data); err != nil {
		http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
	}
}
