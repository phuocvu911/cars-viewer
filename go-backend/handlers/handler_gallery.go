package handlers

import (
	"cars-viewer/models"
	"net/http"
	"strconv"
	"strings"
)

type GalleryData struct {
	Page                             string
	Models                           []models.EnrichedCarModel
	Categories                       []models.Category
	Manufacturers                    []models.Manufacturer
	Drivetrains, Years               []string
	Query, CatF, MfgF, YearF, DriveF string
	ResultCount                      int
	Recommendations                  []models.CarSpecs
}

func GalleryHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	enrichedModels := derived.EnrichedModels
	q := r.URL.Query()
	catF := q.Get("category")
	mfgF := q.Get("manufacturer")
	yearF := q.Get("year")
	driveF := q.Get("drivetrain")
	search := q.Get("q")

	catID, hasCat := atoiOK(catF)
	mfgID, hasMfg := atoiOK(mfgF)
	yearV, hasYear := atoiOK(yearF)

	filtered := make([]models.EnrichedCarModel, 0, len(enrichedModels))

	for _, m := range enrichedModels {
		if hasCat && m.CategoryID != catID {
			continue
		}
		if hasMfg && m.ManufacturerID != mfgID {
			continue
		}
		if hasYear && m.Year != yearV {
			continue
		}
		if driveF != "" && !strings.EqualFold(m.Specifications.Drivetrain, driveF) {
			continue
		}
		if search != "" {
			//free word search across multiple fields, case insensitive.
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

	data := GalleryData{
		Page:            "gallery",
		Models:          filtered,
		Categories:      store.Categories,
		Manufacturers:   store.Manufacturers,
		Drivetrains:     derived.Drivetrains,
		Years:           derived.Years,
		Query:           search,
		CatF:            catF,
		MfgF:            mfgF,
		YearF:           yearF,
		DriveF:          driveF,
		ResultCount:     len(filtered),
		Recommendations: nil,
	}

	// I had to do this in a funny way due to mutex issue with not passing a pointer.
	// -----------------
	recom, err := FetchRecommendations(w, r)

	if err == nil {
		data.Recommendations = recom
	}

	// -----------------
	render(w, "gallery.html", data)
	//this line is for debugging purposes, to see the query parameters in the console when the gallery page is accessed
	//fmt.Println(q)
}

func atoiOK(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return n, true
}
