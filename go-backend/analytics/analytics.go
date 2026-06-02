package analytics

import (
	"sync"
)

const (
	ANALYTICS_MAX_ROWS  int    = 100_000 // do not start counting large datasets
	RUN_ANALYTICS       bool   = true    // Enable or disable analytics
	ANALYTICS_FILE_PATH string = "./suggestions-data.jsonl"
)

type Entry struct {
	Brand   *string `json:"brand"`
	Chassis *string `json:"chassis"`
	Id      *string `json:"id,omitempty"`
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

// Adds entry to the JSON and to inmemory struct
func (self *CookieData) AddEntry(id, brand, chassis string) error {
	self.mu.Lock()
	defer self.mu.Unlock()

	new_data := Entry{Brand: &brand, Chassis: &chassis}
	self.Preferences = append(self.Preferences, new_data)

	new_data.Id = &id

	if err := AppendJSONL(ANALYTICS_FILE_PATH, new_data); err != nil {
		return err
	}

	self.unsafeUpdateCommonMetrics()
	return nil
}

func (self *CookieData) unsafeUpdateCommonMetrics() {

	if len(self.Preferences) == 0 || len(self.Preferences) > ANALYTICS_MAX_ROWS {
		self.UsualBrand = "Undetermined"
		self.UsualChassis = "Undetermined"
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

	maxBrand, maxBrandCount := "Undetermined", 0
	for key, value := range brandCounter {
		if value > maxBrandCount {
			maxBrand, maxBrandCount = key, value
		}
	}

	maxChassis, maxChassisCount := "Undetermined", 0
	for key, value := range chassisCounter {
		if value > maxChassisCount {
			maxChassis, maxChassisCount = key, value
		}
	}

	self.UsualBrand = maxBrand
	self.UsualChassis = maxChassis
}
