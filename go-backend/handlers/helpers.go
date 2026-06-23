package handlers

import (
	"cars-viewer/analytics"
	"cars-viewer/cookies"
	"cars-viewer/models"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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
	LOCAL_CARS_ROUTE         string = "/car/" // ROUTES FOR THE :8080/
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
