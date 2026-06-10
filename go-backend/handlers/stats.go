package handlers

import (
	"net/http"
)

type StatsData struct {
	Page                              string
	TotalModels, TotalMfgs, TotalCats int
	MaxHp, MinHp                      int
	MaxCar, MinCar                    string
	MostCommonCategory                string
	MaxCategoryCount                  int
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	data := buildStatsData()
	render(w, "stats.html", data)
}

func buildStatsData() StatsData {
	mu.RLock()
	defer mu.RUnlock()

	totalModels := len(store.CarModels)
	totalMfgs := len(store.Manufacturers)
	totalCats := len(store.Categories)

	var maxHp, minHp int
	var maxCar, minCar string
	categoryCount := make(map[int]int)

	//find max/min hp/car and most common categoryID
	for _, m := range store.CarModels {
		hp := m.Specifications.Horsepower
		if hp > maxHp {
			maxHp = hp
			maxCar = m.Name
		}
		if hp < minHp || minHp == 0 {
			minHp = hp
			minCar = m.Name
		}
		categoryCount[m.CategoryID]++
	}

	maxCount := 0
	//match categoryID to category name
	var topCatID int
	for categoryID, count := range categoryCount {
		if count > maxCount {
			maxCount, topCatID = count, categoryID
		}
	}
	mostCommonCategory := ""
	for _, cat := range store.Categories {
		if cat.ID == topCatID {
			mostCommonCategory = cat.Name
			break
		}
	}
	return StatsData{
		Page:               "stats",
		TotalModels:        totalModels,
		TotalMfgs:          totalMfgs,
		TotalCats:          totalCats,
		MaxHp:              maxHp,
		MaxCar:             maxCar,
		MinHp:              minHp,
		MinCar:             minCar,
		MostCommonCategory: mostCommonCategory,
		MaxCategoryCount:   maxCount,
	}
}
