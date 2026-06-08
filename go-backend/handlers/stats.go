package handlers

import (
	"net/http"
	"strconv"
)

type StatsData struct {
	Page                              string
	TotalModels, TotalMfgs, TotalCats int
	MaxHp, MaxCar                     string
	MinHp, MinCar                     string
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
	for _, m := range allModels {
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

	mostCommonCategory := ""
	maxCount := 0
	//match categoryID to category name
	for categoryID, count := range categoryCount {
		if count > maxCount {
			maxCount = count
			for _, cat := range allCats {
				if cat.ID == categoryID {
					mostCommonCategory = cat.Name
					break
				}
			}
		}
	}
	return StatsData{
		Page:               "stats",
		TotalModels:        totalModels,
		TotalMfgs:          totalMfgs,
		TotalCats:          totalCats,
		MaxHp:              strconv.Itoa(maxHp),
		MaxCar:             maxCar,
		MinHp:              strconv.Itoa(minHp),
		MinCar:             minCar,
		MostCommonCategory: mostCommonCategory,
		MaxCategoryCount:   maxCount,
	}
}
