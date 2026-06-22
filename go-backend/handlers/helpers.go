package handlers

import (
	"cars-viewer/analytics"
	"cars-viewer/cookies"
	"cars-viewer/models"
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	API_BASE_URL             string = "http://localhost:3000"
	MODELS_ROUTE             string = "/api/models/"
	MODELS__BY_BRAND_ROUTE   string = "/api/models/brand/"
	MODELS__BY_CHASSIS_ROUTE string = "/api/models/chassis/"
	MANUFACTURERS_ROUTE      string = "/api/manufacturers/"
	CATEGORIES_ROUTE         string = "/api/categories/"
	IMG_PATH_PREFIX          string = "/api/images/" // Used for the reverse proxy endpoint and prefixing images
	CAR_ENDPOINT             string = "/api/car/"

	// ROUTES FOR THE :8080/
	LOCAL_CARS_ROUTE string = "/car/"
)

// Add tracking data to the file
func AddTrackingItem(cookie_input *cookies.CookieCtx, car_obj *models.Car) error {

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

	analytics.LiveCookieData.Data[cookie_input.ShortCookie.Value].AddEntry(cookie_input.LongCookie.Value, cookie_input.ShortCookie.Value, car_obj.ManufactDetails.Name, car_obj.Category.Name)
	return nil

}

var store models.DataStore

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
func FetchCarByID(car_id string, errChan chan<- error, carpointer *models.Car) {

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

// Used for enriching Car object concurrently with FetchCarCategory
// No need to use sync.Mutex
func FetchCarManufacturer(errChan chan<- error, carpointer *models.Car) {

	if carpointer.DataPerID.ManufactrurerID < 1 {
		errChan <- errors.New("Manufactrurer ID has not been assigned or is invalid. ")
		return
	}

	err := FetchDataFromAPIByRouteAndID(MANUFACTURERS_ROUTE, carpointer.DataPerID.ManufactrurerID, &carpointer.ManufactDetails)

	if err != nil {
		errChan <- err
		return
	}

	errChan <- nil

}

// Used for enriching Car object concurrently with FetchCarManufacturer
// No need to use sync.Mutex
func FetchCarCategory(errChan chan<- error, carpointer *models.Car) {

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
