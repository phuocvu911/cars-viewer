package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// those tructs mirror the JSONs, can be moved to ./model
type Manufacturer struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	FoundingYear int    `json:"foundingYear"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Specifications struct {
	Engine       string `json:"engine"`
	Horsepower   int    `json:"horsepower"`
	Transmission string `json:"transmission"`
	Drivetrain   string `json:"drivetrain"`
}

type CarModel struct {
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	ManufacturerID int            `json:"manufacturerId"`
	CategoryID     int            `json:"categoryId"`
	Year           int            `json:"year"`
	Specifications Specifications `json:"specifications"`
	Image          string         `json:"image"`
}

// unifies add tables
type EnrichedCarModel struct {
	CarModel
	ManufacturerName    string
	ManufacturerCountry string
	FoundingYear        int
	CategoryName        string
}

type DataStore struct {
	Manufacturers []Manufacturer `json:"manufacturers"`
	Categories    []Category     `json:"categories"`
	CarModels     []CarModel     `json:"carModels"`
}

var store DataStore

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	// Fetch data from Cars API at localhost:3000
	apiURL := "http://localhost:3000"

	fetchAPI := func(path string, v any) error {
		resp, err := http.Get(apiURL + path)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}
		return json.NewDecoder(resp.Body).Decode(v)
	}

	errM := fetchAPI("/api/models", &store.CarModels)
	errMf := fetchAPI("/api/manufacturers", &store.Manufacturers)
	errC := fetchAPI("/api/categories", &store.Categories)

	if errM != nil || errMf != nil || errC != nil {
		log.Printf("Failed to pull data from %s. Make sure the Cars API is running! Models:%v Mfgs:%v Cats:%v", apiURL, errM, errMf, errC)
	} else {
		log.Printf("Pulled %d models, %d manufacturers, %d categories from Cars API",
			len(store.CarModels), len(store.Manufacturers), len(store.Categories))
	}

	// Initialize and cache templates
	initTemplates()

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Proxy image requests to localhost:3000
	remoteURL, _ := url.Parse(apiURL)
	proxy := httputil.NewSingleHostReverseProxy(remoteURL)
	mux.HandleFunc("/api/images/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	// Page routes
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/car/", handleCar)
	mux.HandleFunc("/compare", handleCompare)
	mux.HandleFunc("/recommend", handleRecommend)

	log.Printf("AutoVault running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
