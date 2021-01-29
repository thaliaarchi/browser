package firefox

import "github.com/andrewarchi/browser/timefmt"

// Times contains installation times in times.json.
type Times struct {
	Created  timefmt.UnixMilli `json:"created"`
	FirstUse timefmt.UnixMilli `json:"firstUse"`
}

// ParseTimes parses the times.json file in a Firefox profile.
func ParseTimes(filename string) (*Times, error) {
	var times Times
	if err := parseJSON(filename, &times); err != nil {
		return nil, err
	}
	return &times, nil
}
