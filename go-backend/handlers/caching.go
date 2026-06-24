package handlers

import (
	"cars-viewer/models"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"
	"text/template"
	"time"
)

var store models.DataStore

// RWMutex: multiple readers can read simultaneously, but a writer gets exclusive access for lock and unlock
var mu sync.RWMutex

// Initialize the store by fetching all models and categories
func InitStore() error {
	var (
		carModels     []models.CarModel
		categories    []models.Category
		manufacturers []models.Manufacturer
	)

	errChan := make(chan error, 3)
	go func() {
		errChan <- fetchDataFromAPI(MODELS_ROUTE, &carModels)
	}()
	go func() {
		errChan <- fetchDataFromAPI(CATEGORIES_ROUTE, &categories)
	}()
	go func() {
		errChan <- fetchDataFromAPI(MANUFACTURERS_ROUTE, &manufacturers)
	}()
	for range 3 {
		if err := <-errChan; err != nil {
			return err
		}
	}
	close(errChan)
	//lock and update the store with new data when refreshing to prevent data race
	mu.Lock()
	store.CarModels = carModels
	store.Categories = categories
	store.Manufacturers = manufacturers
	buildDerived() //rebuild the derived data after updating the store
	mu.Unlock()
	return nil
}

func fetchDataFromAPI(path string, v any) error {
	resp, err := http.Get(API_BASE_URL + path)
	if err != nil {
		return errors.New("failed to fetch data from " + path + ": " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected status: " + strconv.Itoa(resp.StatusCode) + " for endpoint: " + path)
	}
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return errors.New("failed to decode JSON from " + path + ": " + err.Error())
	}
	return nil
}

// Start a background goroutine to refresh the store every 10 minutes.
func StoreRefresh(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := InitStore(); err != nil {
				log.Println("store refresh failed:", err)
			} else {
				log.Println("store refreshed successfully")
			}
		case <-ctx.Done():
			return
		}
	}
}

var templates = make(map[string]*template.Template)

// Initialize templates and cache them in a map for efficient rendering. This is called once at main.
func InitTemplates() {
	funcMap := template.FuncMap{
		"itoa": strconv.Itoa,
	}

	pages := []string{"home.html", "gallery.html", "car.html", "compare.html", "stats.html", "404.html"}
	for _, page := range pages {
		templates[page] = template.Must(template.New(page).Funcs(funcMap).ParseFiles(
			"templates/index.html",
			"templates/navfooter.html",
			"templates/"+page,
		))
	}
}

// Render a template with the given data
func render(w http.ResponseWriter, page string, data any) {
	t, ok := templates[page]
	if !ok {
		log.Printf("Template %s not found in cache", page)
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	err := t.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Printf("Error rendering template %s: %v", page, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
