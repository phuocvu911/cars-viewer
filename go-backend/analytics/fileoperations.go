package analytics

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
)

func AppendJSONL(filePath string, record Entry) error {

	if record.ShortID == nil || record.LongID == nil {
		return errors.New("Entry does not include ID.")
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(record)

	if err != nil {
		return err
	}

	return nil
}

func LoadAndAggregate(filePath string) (*UserPreferences, error) {

	file, err := os.Open(filePath)

	if err != nil {

		if errors.Is(err, fs.ErrPermission) {
			return nil, errors.New("No permissions for filepath. ")
		}

		if errors.Is(err, fs.ErrNotExist) {
			return &UserPreferences{Data: make(map[string]*CookieData)}, nil
		}
		return nil, err
	}
	defer file.Close()

	userMap := make(map[string]*CookieData)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()

		if len(line) == 0 {
			continue
		}

		var entry Entry
		err := json.Unmarshal(line, &entry)

		if err != nil {
			return nil, err
		}

		if entry.ShortID == nil {
			continue
		}

		userId := *entry.ShortID

		_, found := userMap[userId]

		if !found {
			userMap[userId] = &CookieData{
				Preferences: []Entry{},
			}
		}

		userMap[userId].Preferences = append(userMap[userId].Preferences, entry)
	}

	err = scanner.Err()

	if err != nil {
		return nil, err
	}

	for _, cookieData := range userMap {
		cookieData.unsafeUpdateCommonMetrics()
	}

	return &UserPreferences{Data: userMap}, nil
}
