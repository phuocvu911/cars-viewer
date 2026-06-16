package handlers

import (
	"context"
	"cars-viewer/analytics"
	"cars-viewer/cookies"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	API_BASE_URL        string = "http://localhost:3000"
	MODELS_ROUTE        string = "/api/models/"
	MANUFACTURERS_ROUTE string = "/api/manufacturers/"
	CATEGORIES_ROUTE    string = "/api/categories/"
	IMG_PATH_PREFIX     string = "/api/images/" // Used for the reverse proxy endpoint and prefixing images
)

// Add tracking data to the file
func AddTrackingItem(cookie_input *cookies.CookieCtx, car_obj *Car) error {

	// Check if any cookie is missing
	if cookie_input.AllowTracking == nil || cookie_input.LongCookie == nil || cookie_input.ShortCookie == nil {
		return errors.New("nil cookie found.")
	}

	// Check if any cookie is missing Value field
	if cookie_input.AllowTracking.Value == "" || cookie_input.LongCookie.Value == "" || cookie_input.ShortCookie.Value == "" {
		return errors.New("Value without value found.")
	}

	// Avoid memory errors by verifying the struct exists
	if analytics.LiveCookieData.Data[cookie_input.ShortCookie.Value] == nil {
		analytics.LiveCookieData.Data[cookie_input.ShortCookie.Value] = &analytics.CookieData{}
	}

	analytics.LiveCookieData.Data[cookie_input.ShortCookie.Value].AddEntry(cookie_input.LongCookie.Value, cookie_input.ShortCookie.Value, car_obj.ManufactDetails.Make, car_obj.Category.Name)
	log.Println(analytics.LiveCookieData.Data[cookie_input.ShortCookie.Value].UsualBrand)
	log.Println(analytics.LiveCookieData.Data[cookie_input.ShortCookie.Value].UsualChassis)
	return nil

}

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
	Category        Category
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
	CarID            int              `json:"id"`             // The cars own individual unique id
	ManufactrurerID  int              `json:"manufacturerId"` // Holds data to the car manufacturer
	CategoryID       int              `json:"categoryId"`
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

// Use to fetch Car by id.
// Basically fetch CarSpecs struct.
func FetchCarByID(car_id string, errChan chan<- error, carpointer *Car) {

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
	// Add nil to the channel to confirm all went smoothly
	errChan <- nil
}

	err := FetchDataFromAPIByRouteAndID(MANUFACTURERS_ROUTE, carpointer.DataPerID.ManufactrurerID, &carpointer.ManufactDetails)
func FetchCarManufacturer(errChan chan<- error, carpointer *Car) {

	if carpointer.DataPerID.ManufactrurerID < 1 {
		errChan <- errors.New("Manufactrurer ID has not been assigned or is invalid. ")
		return
	}

	err := FetchDataFromAPIByRouteAndID(MANUFACTURER_ROUTE, carpointer.DataPerID.ManufactrurerID, &carpointer.ManufactDetails)

	if err != nil {
		errChan <- err
		return
	}

	errChan <- nil

}

// Used for enriching Car object concurrently with FetchCarManufacturer
// No need to use sync.Mutex
func FetchCarCategory(errChan chan<- error, carpointer *Car) {

	if carpointer.DataPerID.CategoryID < 1 {
		errChan <- errors.New("Category ID has not been assigned or is invalid. ")
	}

	err := FetchDataFromAPIByRouteAndID(CATEGORIES_ROUTE, carpointer.DataPerID.CategoryID, &carpointer.Category)

	if err != nil {
		errChan <- err
		return
	}

	errChan <- nil

}

// RWMutex: multiple readers can read simultaneously, but a writer gets exclusive access for lock and unlock
var mu sync.RWMutex

// Initialize the store by fetching all models and categories
func InitStore() error {
	var (
		models        []CarModel
		categories    []Category
		manufacturers []Manufacturer
	)

	errChan := make(chan error, 3)
	go func() {
		errChan <- fetchDataFromAPI(MODELS_ROUTE, &models)
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

	//lock and update the store with new data when refreshing to prevent data race
	mu.Lock()
	store.CarModels = models
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

// Enrich a single car model with manufacturer and category details, like joining tables in SQL, we mathching by ID.

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
