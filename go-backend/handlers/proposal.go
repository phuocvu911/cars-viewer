package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

// =============================================================================
// PROPOSAL — fixes for inefficiencies #1 and #2
//
// This file is a self-contained, COMPILING proposal. New symbols are suffixed
// with "Proposed" so the package still builds next to the current code. When
// you adopt this, fold these into the real functions (mapping noted on each
// block) and delete this file.
//
// Core idea: the store only changes every 10 minutes (StoreRefresh), yet the
// gallery currently recomputes enrichment (O(N*M)) and the filter dropdown
// sets (O(N)) on EVERY request. Compute all of that ONCE per refresh, under the
// existing write lock, and let request handlers just read the cached results.
// =============================================================================

// ---- Cached, precomputed view of the store -------------------------------
//
// Populated once per refresh (inside InitStore). Guarded by the existing `mu`.
// Replaces the per-request work in enrichAll() and GalleryHandler.
type derivedData struct {
	EnrichedModels []EnrichedCarModel // fix #1: full enriched slice, built once
	Years          []string           // fix #2: unique years for the dropdown
	Drivetrains    []string           // fix #2: unique drivetrains for the dropdown
}

var derived derivedData

// rebuildDerived recomputes every cached projection from the raw store slices.
//
// PROPOSAL: call this from InitStore() while still holding mu.Lock(), right
// after assigning store.CarModels / Categories / Manufacturers:
//
//	mu.Lock()
//	store.CarModels = models
//	store.Categories = categories
//	store.Manufacturers = manufacturers
//	rebuildDerived()          // <-- add this
//	mu.Unlock()
//
// It assumes the caller already holds the write lock.
func rebuildDerived() {
	// fix #1: O(1) lookup maps instead of scanning the slices for every model.
	mfgByID := make(map[int]Manufacturer, len(store.Manufacturers))
	for _, m := range store.Manufacturers {
		mfgByID[m.ID] = m
	}
	catByID := make(map[int]Category, len(store.Categories))
	for _, c := range store.Categories {
		catByID[c.ID] = c
	}

	enriched := make([]EnrichedCarModel, 0, len(store.CarModels))

	// fix #2: collect the dropdown sets in the SAME pass we enrich, so we touch
	// each model once instead of looping the models twice per request.
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

// enrichWithMaps is the proposed replacement for enrich() (helpers.go:206).
//
// Two wins over the original:
//   - O(1) map lookups instead of two linear scans over the slices.
//   - the manufacturer lookup no longer keeps scanning after a match (the
//     original loop at helpers.go:212-218 has no break).
func enrichWithMaps(m CarModel, mfgByID map[int]Manufacturer, catByID map[int]Category) EnrichedCarModel {
	e := EnrichedCarModel{CarModel: m}
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

// ---- Proposed GalleryHandler body -----------------------------------------
//
// PROPOSAL: replace the body of GalleryHandler (gallery.go:24). It now reads the
// cached projections instead of recomputing them:
//
//	before (per request): models := enrichAll(); ... rebuild yearSet/driveSet ...
//	after  (per request): read derived.EnrichedModels / .Years / .Drivetrains
//
// The filtering loop is unchanged EXCEPT the filter ints are parsed once up
// front (also addresses finding #3) instead of once per model.
func GalleryHandlerProposed(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	q := r.URL.Query()
	catF := q.Get("category")
	mfgF := q.Get("manufacturer")
	yearF := q.Get("year")
	driveF := q.Get("drivetrain")
	search := q.Get("q")

	// fix #3 (bonus): parse filter ints once, not once per model.
	catID, hasCat := atoiOK(catF)
	mfgID, hasMfg := atoiOK(mfgF)
	yearV, hasYear := atoiOK(yearF)
	searchLower := strings.ToLower(search)

	models := derived.EnrichedModels // fix #1: cached, no per-request enrichment
	filtered := make([]EnrichedCarModel, 0, len(models))
	for _, m := range models {
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
		if search != "" && !matchesSearch(m, searchLower) {
			continue
		}
		filtered = append(filtered, m)
	}

	data := GalleryData{
		Page:          "gallery",
		Models:        filtered,
		Categories:    store.Categories,
		Manufacturers: store.Manufacturers,
		Drivetrains:   derived.Drivetrains, // fix #2: cached
		Years:         derived.Years,       // fix #2: cached
		Query:         search,
		CatF:          catF,
		MfgF:          mfgF,
		YearF:         yearF,
		DriveF:        driveF,
		ResultCount:   len(filtered),
	}
	render(w, "gallery.html", data)
}

// atoiOK reports whether s parses as an int; ok is false for "" or junk, so the
// caller can skip the filter entirely (matching the current "" guard behaviour).
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

// matchesSearch keeps the original 5-field substring search (gallery.go:57-61);
// the needle is pre-lowered once by the caller.
func matchesSearch(m EnrichedCarModel, needle string) bool {
	return strings.Contains(strings.ToLower(m.Name), needle) ||
		strings.Contains(strings.ToLower(m.ManufacturerName), needle) ||
		strings.Contains(strings.ToLower(m.CategoryName), needle) ||
		strings.Contains(strings.ToLower(m.ManufacturerCountry), needle) ||
		strings.Contains(strings.ToLower(m.Specifications.Engine), needle)
}

// enrichAll (helpers.go:232) collapses to a one-line cache read once #1 lands:
//
//	func enrichAll() []EnrichedCarModel { return derived.EnrichedModels }
//
// CompareHandler / StatsHandler can keep enrich() per-id, or read
// derived.EnrichedModels too — either way no full re-enrichment per request.
