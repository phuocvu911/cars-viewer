package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	API_BASE_URL        string = "http://localhost:3000"
	MODELS_ROUTE        string = "/api/models/"
	MANUFACTURERS_ROUTE string = "/api/manufacturers/"
	CATEGORIES_ROUTE    string = "/api/categories/"

	IMG_PATH_PREFIX string = "/api/images/" // Used for the reverse proxy endpoint and prefixing images
)

// Global store for all models and categories
type DataStore struct {
	Manufacturers []Manufacturer `json:"manufacturers"`
	Categories    []Category     `json:"categories"`
	CarModels     []CarModel     `json:"carModels"`
}

var store DataStore

type CarModel struct {
	ID             int              `json:"id"`
	Name           string           `json:"name"`
	ManufacturerID int              `json:"manufacturerId"`
	CategoryID     int              `json:"categoryId"`
	Year           int              `json:"year"`
	Image          string           `json:"image"`
	Specifications TechnicalDetails `json:"specifications"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// DATA STRUCTS
type Car struct {
	DataPerID       CarSpecs
	ManufactDetails Manufacturer
	Page            string
}

// Access via /api/manufacturers/{id}
type Manufacturer struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	CountryOfOrigin string `json:"country"`
	FoundingYear    int    `json:"foundingYear"`
}

type TechnicalDetails struct {
	Engine       string `json:"engine"`
	Horsepower   int    `json:"horsepower"`
	Transmission string `json:"transmission"`
	Drivetrain   string `json:"drivetrain"`
}

// Access via /api/models/{id}
type CarSpecs struct {
	CarID           int `json:"id"`             // The cars own individual unique id
	ManufactrurerID int `json:"manufacturerId"` // Holds data to the car manufacturer
	CategoryID      int `json:"categoryId"`

	MakeModel        string           `json:"name"` // e.g. "Audi Q5"
	Year             int              `json:"year"`
	TechnicalDetails TechnicalDetails `json:"specifications"`
	ImgSrc           string           `json:"image"`
}

func FetchDataFromAPIByRouteAndID(route string, id int, DataModel any) error {

	endpoint := API_BASE_URL + route + strconv.Itoa(id)

	res, err := http.Get(endpoint)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("API error. Statuscode != 200. code: " + strconv.Itoa(res.StatusCode))
	}

	err = json.NewDecoder(res.Body).Decode(DataModel)

	if err != nil {
		return err
	}

	return nil
}

// Use to fetch all related information to Car page.
// Second request needs manufacturer id so it cannot be concurrent
func FetchCar(car_id string, errChan chan<- error, carpointer *Car) {

	int_id, err := strconv.Atoi(car_id)

	if err != nil {
		errChan <- err
		return
	}

	err = FetchDataFromAPIByRouteAndID(MODELS_ROUTE, int_id, &carpointer.DataPerID)

	if err != nil {
		errChan <- err
		return
	}

	err = FetchDataFromAPIByRouteAndID(MANUFACTURERS_ROUTE, carpointer.DataPerID.ManufactrurerID, &carpointer.ManufactDetails)

	if err != nil {
		errChan <- err
		return
	}

	errChan <- nil
}

// Initialize the store by fetching all models and categories
func InitStore() error {
	errChan := make(chan error, 3)
	go func() {
		errChan <- fetchDataFromAPI(MODELS_ROUTE, &store.CarModels)
	}()
	go func() {
		errChan <- fetchDataFromAPI(CATEGORIES_ROUTE, &store.Categories)
	}()
	go func() {
		errChan <- fetchDataFromAPI(MANUFACTURERS_ROUTE, &store.Manufacturers)
	}()
	for range 3 {
		if err := <-errChan; err != nil {
			return err
		}
	}
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

// Enrich a single car model with manufacturer and category details, like joining tables in SQL, we mathching by ID.
func enrich(m CarModel) EnrichedCarModel {
	enriched := EnrichedCarModel{
		CarModel: m,
	}

	// Get manufacturer details from store
	for _, mfg := range store.Manufacturers {
		if mfg.ID == m.ManufacturerID {
			enriched.ManufacturerName = mfg.Name
			enriched.ManufacturerCountry = mfg.CountryOfOrigin
			enriched.FoundingYear = mfg.FoundingYear
		}
	}

	// Get category name from store
	for _, cat := range store.Categories {
		if cat.ID == m.CategoryID {
			enriched.CategoryName = cat.Name
			break
		}
	}

	return enriched
}

// Enrich all car models
func enrichAll() []EnrichedCarModel {
	enriched := make([]EnrichedCarModel, 0, len(store.CarModels))
	for _, model := range store.CarModels {
		enriched = append(enriched, enrich(model))
	}
	return enriched
}

var templates = make(map[string]*template.Template)

// Initialize templates and cache them in a map for efficient rendering. This is called once at main.
func InitTemplates() {
	funcMap := template.FuncMap{
		"itoa": strconv.Itoa,
	}

	pages := []string{"home.html", "gallery.html", "car.html", "compare.html", "stats.html"}
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
