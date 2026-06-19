package handlers

import (
	"cars-viewer/analytics"
	"cars-viewer/cookies"
	"errors"
	"math/rand/v2"
	"net/http"
	"sync"
)

const (
	RECOMMENDATIONS_MAX_COUNT int = 3
)

func MergeAndReturnUnique(list_1, list_2 []CarSpecs) []CarSpecs {
	idsInList1 := make(map[int]bool)
	idsInList2 := make(map[int]bool)

	for _, car := range list_1 {
		idsInList1[car.CarID] = true
	}
	for _, car := range list_2 {
		idsInList2[car.CarID] = true
	}

	var result []CarSpecs

	for _, car := range list_1 {
		result = append(result, car)
	}

	for _, car := range list_2 {
		if !idsInList1[car.CarID] {
			result = append(result, car)
		}
	}

	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	return result
}

func FetchRecommendations(r *http.Request) ([]CarSpecs, error) {

	cookieCtx, no_problem := r.Context().Value(cookies.CookieCtxKey{}).(cookies.CookieCtx)

	if !no_problem {
		return nil, errors.New("Failed to read cookies. ")
	}
	var user_preferred_brand, user_preferred_chassis string

	// Check if the cookie exists
	if analytics.LiveCookieData.Data[cookieCtx.ShortCookie.Value] != nil {
		if analytics.LiveCookieData.Data[cookieCtx.ShortCookie.Value].UsualBrand != "" {
			user_preferred_brand = analytics.LiveCookieData.Data[cookieCtx.ShortCookie.Value].UsualBrand
		}

		if analytics.LiveCookieData.Data[cookieCtx.ShortCookie.Value].UsualChassis != "" {
			user_preferred_chassis = analytics.LiveCookieData.Data[cookieCtx.ShortCookie.Value].UsualChassis
		}
	} else {
		return nil, http.ErrNoCookie
	}

	var wg sync.WaitGroup
	var list1, list2 []CarSpecs
	var get_brands, get_chassis bool
	var threads int

	if user_preferred_brand != "" && user_preferred_brand != analytics.UNDETERMINED_VALUE_NAME {
		get_brands = true
		threads++
	}
	if user_preferred_chassis != "" && user_preferred_chassis != analytics.UNDETERMINED_VALUE_NAME {
		get_chassis = true
		threads++
	}

	if threads == 0 {
		// No data
		return nil, errors.New("Broken data for the user.")
	}

	var errChan chan (error) = make(chan (error), threads)

	if get_chassis {

		wg.Go(
			func() {
				var tmp []CarSpecs
				err := fetchDataFromAPI(MODELS__BY_CHASSIS_ROUTE+user_preferred_chassis, &tmp)
				if err == nil {
					list1 = tmp
				}
				errChan <- err
			})
	}

	if get_brands {

		wg.Go(
			func() {
				var tmp []CarSpecs
				err := fetchDataFromAPI(MODELS__BY_BRAND_ROUTE+user_preferred_brand, &tmp)
				if err == nil {
					list2 = tmp
				}
				errChan <- err
			})
	}
	wg.Wait()

	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	// Return always list1

	list1 = MergeAndReturnUnique(list2, list1)

	if len(list1) > RECOMMENDATIONS_MAX_COUNT {
		list1 = list1[:RECOMMENDATIONS_MAX_COUNT]
	}

	return list1, nil
}
