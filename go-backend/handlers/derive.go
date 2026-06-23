package handlers

import (
	"cars-viewer/models"
	"strconv"
)

type derivedData struct {
	EnrichedModels []models.EnrichedCarModel
	Years          []string
	Drivetrains    []string
}

var derived derivedData

//buildDerived computes every derivedData field for used later in html pages.
func buildDerived() {
	//build lookup maps instead of scanning the slices for every model. O(1) complexity instead of O(n).
	mfgByID := make(map[int]models.Manufacturer, len(store.Manufacturers))
	for _, m := range store.Manufacturers {
		mfgByID[m.ID] = m
	}
	catByID := make(map[int]models.Category, len(store.Categories))
	for _, c := range store.Categories {
		catByID[c.ID] = c
	}

	//pre allocate the enriched slice to avoid resizing during append.
	enriched := make([]models.EnrichedCarModel, 0, len(store.CarModels))

	//use map set to catch only unique years and drivetrains. values will be discarded, so struct{}{} is used to save memory.
	yearSet := make(map[string]struct{})
	driveSet := make(map[string]struct{})

	for _, m := range store.CarModels {
		enriched = append(enriched, enrichWithMaps(m, mfgByID, catByID))
		yearSet[strconv.Itoa(m.Year)] = struct{}{}
		driveSet[m.Specifications.Drivetrain] = struct{}{}
	}

	years := make([]string, 0, len(yearSet))
	for y := range yearSet {
		years = append(years, y)
	}
	drives := make([]string, 0, len(driveSet))
	for d := range driveSet {
		drives = append(drives, d)
	}

	derived = derivedData{
		EnrichedModels: enriched,
		Years:          years,
		Drivetrains:    drives,
	}
}

//enrichWithMaps builds an EnrichedCarModel from a CarModel using the provided lookup maps. O(1) complexity.
func enrichWithMaps(m models.CarModel, mfgByID map[int]models.Manufacturer, catByID map[int]models.Category) models.EnrichedCarModel {
	e := models.EnrichedCarModel{CarModel: m}
	if mfg, ok := mfgByID[m.ManufacturerID]; ok {
		e.ManufacturerName = mfg.Name
		e.ManufacturerCountry = mfg.CountryOfOrigin
		e.FoundingYear = mfg.FoundingYear
	}
	if cat, ok := catByID[m.CategoryID]; ok {
		e.CategoryName = cat.Name
	}
	return e
}
