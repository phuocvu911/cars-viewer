package analytics

import (
	"sync"
)

const (
	ANALYTICS_MAX_ROWS  int    = 100_000 // do not start counting large datasets
	ANALYTICS_FILE_PATH string = "./suggestions-data.jsonl"

	// After the user has this amount of clicks to a chassis type or car brand,
	// the program will calculate preferred types.
	RECOMMENDATIONS_THRESHOLD int    = 2
	UNDETERMINED_VALUE_NAME   string = "undetermined"
)

type Entry struct {
	Brand   *string `json:"brand"`
	Chassis *string `json:"chassis"`
	ShortID *string `json:"short_id"`
	LongID  *string `json:"long_id"`
}

type CookieData struct {
	mu           sync.Mutex // disable updating CookieData from 2 threads at the same time
	Preferences  []Entry
	UsualBrand   string
	UsualChassis string
}

type UserPreferences struct {
	Data map[string]*CookieData // Data key is the Cookie given to the user!
}

// Adds entry to the JSONL and to in-memory struct
func (self *CookieData) AddEntry(user_long_id, user_short_id, brand, chassis string) error {

	self.mu.Lock()
	defer self.mu.Unlock()

	new_data := Entry{Brand: &brand, Chassis: &chassis}
	self.Preferences = append(self.Preferences, new_data)

	new_data.ShortID = &user_short_id
	new_data.LongID = &user_long_id

	if err := AppendJSONL(ANALYTICS_FILE_PATH, new_data); err != nil {
		return err
	}

	self.unsafeUpdateCommonMetrics()
	return nil
}

// Only call from a goroutine that has locked the struct!!!
// Calling from multiple goroutines can cause race conditions
// and data issues if the data is not locked.
func (self *CookieData) unsafeUpdateCommonMetrics() {

	if len(self.Preferences) == 0 || len(self.Preferences) > ANALYTICS_MAX_ROWS {
		self.UsualBrand = UNDETERMINED_VALUE_NAME
		self.UsualChassis = UNDETERMINED_VALUE_NAME
		return
	}

	brandCounter := map[string]int{}
	chassisCounter := map[string]int{}

	for _, entry := range self.Preferences {
		if entry.Brand != nil {
			brandCounter[*entry.Brand]++
		}
		if entry.Chassis != nil {
			chassisCounter[*entry.Chassis]++
		}
	}

	maxBrand, maxBrandCount := UNDETERMINED_VALUE_NAME, 0
	for key, value := range brandCounter {
		if value > maxBrandCount {
			maxBrand, maxBrandCount = key, value
		}
	}

	maxChassis, maxChassisCount := UNDETERMINED_VALUE_NAME, 0
	for key, value := range chassisCounter {
		if value > maxChassisCount {
			maxChassis, maxChassisCount = key, value
		}
	}

	if maxBrandCount >= RECOMMENDATIONS_THRESHOLD {
		self.UsualBrand = maxBrand
	} else {
		self.UsualBrand = UNDETERMINED_VALUE_NAME
	}

	if maxChassisCount >= RECOMMENDATIONS_THRESHOLD {
		self.UsualChassis = maxChassis
	} else {
		self.UsualChassis = UNDETERMINED_VALUE_NAME
	}
}

var LiveCookieData, _ = LoadAndAggregate(ANALYTICS_FILE_PATH)
