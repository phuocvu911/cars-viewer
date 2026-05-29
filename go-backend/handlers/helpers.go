package handlers

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strconv"
)

const (
	// GO ENDPOINTS
	CARS_ENDPOINT string = "/cars/"

	// ROUTES FOR THE JS API
	API_BASE_URL         string = "http://localhost:3000"
	MODELS_ROUTE         string = "/api/models/"
	MANUFACTURER_ROUTE   string = "/api/manufacturers/"
	CATEGORIES_ROUTE     string = "/api/categories/"
	ALL_MODELS_ROUTE     string = "/api/models"
	ALL_CATEGORIES_ROUTE string = "/api/categories"

	IMG_PATH_PREFIX string = "/api/images/" // Used for the reverse proxy endpoint and prefixing images

)

// Global store for all models and categories
type Store struct {
	CarModels  []CarModel
	Categories map[int]string // categoryID -> categoryName
}

var store Store

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
	Category        Category
}

// Access via /api/manufacturers/{id}
type Manufacturer struct {
	Make            string `json:"name"`
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

func FetchCarManufacturer(errChan chan<- error, carpointer *Car) {

	if carpointer.DataPerID.ManufactrurerID < 1 {
		errChan <- errors.New("Manufactrurer ID has not been assigned or is invalid. ")
		return
	}

	err := FetchDataFromAPIByRouteAndID(CATEGORIES_ROUTE, carpointer.DataPerID.ManufactrurerID, &carpointer.ManufactDetails)

	if err != nil {
		errChan <- err
		return
	}

	errChan <- nil

}

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

// Initialize the store by fetching all models and categories
func InitStore() error {
	// Fetch all car models
	res, err := http.Get(API_BASE_URL + ALL_MODELS_ROUTE)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("failed to fetch models: status code " + strconv.Itoa(res.StatusCode))
	}

	err = json.NewDecoder(res.Body).Decode(&store.CarModels)
	if err != nil {
		return err
	}

	// Fetch all categories
	store.Categories = make(map[int]string)
	res, err = http.Get(API_BASE_URL + ALL_CATEGORIES_ROUTE)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		var categories []Category
		err = json.NewDecoder(res.Body).Decode(&categories)
		if err == nil {
			for _, cat := range categories {
				store.Categories[cat.ID] = cat.Name
			}
		}
	}

	return nil
}

// Enrich a single car model with manufacturer and category details
func enrich(m CarModel) EnrichedCarModel {
	enriched := EnrichedCarModel{
		CarModel: m,
	}

	// Fetch manufacturer details
	var manufacturer Manufacturer
	if err := FetchDataFromAPIByRouteAndID(MANUFACTURER_ROUTE, m.ManufacturerID, &manufacturer); err == nil {
		enriched.ManufacturerName = manufacturer.Make
		enriched.ManufacturerCountry = manufacturer.CountryOfOrigin
		enriched.FoundingYear = manufacturer.FoundingYear
	}

	// Get category name from store
	enriched.CategoryName = store.Categories[m.CategoryID]

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

// Render a template with the given data
func render(w http.ResponseWriter, templateName string, data any) error {
	tmpl, err := template.ParseFiles(
		"./templates/index.html",
		"./templates/navfooter.html",
		"./templates/"+templateName,
	)
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(w, "index.html", data)
}
