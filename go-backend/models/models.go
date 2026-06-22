package models

type EnrichedCarModel struct {
	CarModel
	ManufacturerName    string
	ManufacturerCountry string
	FoundingYear        int
	CategoryName        string
}

// Global store for all models and categories
type DataStore struct {
	Manufacturers []Manufacturer `json:"manufacturers"`
	Categories    []Category     `json:"categories"`
	CarModels     []CarModel     `json:"carModels"`
}

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
