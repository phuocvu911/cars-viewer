package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const (
	API_BASE_URL       string = "http://localhost:3000"
	MODELS_ROUTE       string = "/api/models/"
	MANUFACTURER_ROUTE string = "/api/manufacturers/"
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
	FoundingYear    int    `json:"foundingYear"`
}

// Access via /api/models/{id}
type CarSpecs struct {
	CarID           int `json:"id"`             // The cars own individual unique id
	ManufactrurerID int `json:"manufacturerId"` // Holds data to the car manufacturer
	CategoryID      int `json:"categoryId"`

	MakeModel    string `json:"name"` // e.g. "Audi Q5"
	Year         int    `json:"year"`
	Engine       string `json:"specifications.engine"`
	EngineHP     string `json:"specifications.horsepower"`
	Transmission string `json:"specifications.transmission"`
	Drivetrain   string `json:"specifications.drivetrain"`
	ImgSrc       string `json:"image"`
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

// Use to fetch all related information to Car page. Second request needs manufacturer id so it cannot be concurrent
func FetchCar(car_id string, errChan chan<- error, carpointer *Car) {

	int_id, err := strconv.Atoi(car_id)

	if err != nil {
		errChan <- err
		return
	}

	err = FetchDataFromAPIByRouteAndID(MODELS_ROUTE, int_id, &carpointer.DataPerID)
	fmt.Println(carpointer)
	if err != nil {
		errChan <- err
		return
	}

	err = FetchDataFromAPIByRouteAndID(MANUFACTURER_ROUTE, carpointer.DataPerID.ManufactrurerID, &carpointer.ManufactDetails)
	fmt.Println(carpointer)
	if err != nil {
		errChan <- err
		return
	}
}
