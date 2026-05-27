package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

const (
	API_BASE_URL       string = "localhost:3000"
	MODELS_ROUTE       string = "/models/"
	MANUFACTURER_ROUTE string = "/manufacturers/"
)

// DATA STRUCTS
type Car struct {
	DataPerID       CarSpecs
	ManufactDetails Manufacturer
}

// Access via /api/manufacturers/{id}
type Manufacturer struct {
	Make            string `json:"name"`
	CountryOfOrigin string `json:"country"`
	FoundingYear    string `json:"foundingYear"`
}

// Access via /api/models/{id}
type CarSpecs struct {
	CarID           string `json:"id"`             // The cars own individual unique id
	ManufactrurerID string `json:"manufacturerId"` // Holds data to the car manufacturer
	CategoryID      string `json:"categoryId"`

	MakeModel    string `json:"name"`
	Year         string `json:"year"`
	Engine       string `json:"specifications.engine"`
	EngineHP     string `json:"specifications.horsepower"`
	Transmission string `json:"specifications.transmission"`
	Drivetrain   string `json:"specifications.drivetrain"`
	ImgSrc       string `json:"image"`
}

func FetchDataFromAPIByRouteAndID(route, id string, DataModel any) error {

	endpoint := API_BASE_URL + route + id

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

// Use to fetch all related information to Car page. Second request needs manufacturer id so it cannot be concurrent
func FetchCar(car_id string, errChan chan<- error, carpointer *Car) {

	err := FetchDataFromAPIByRouteAndID(MODELS_ROUTE, car_id, &carpointer.DataPerID)

	if err != nil {
		errChan <- err
		return
	}

	FetchDataFromAPIByRouteAndID(MANUFACTURER_ROUTE, carpointer.DataPerID.ManufactrurerID, &carpointer.ManufactDetails)

	if err != nil {
		errChan <- err
		return
	}
}
