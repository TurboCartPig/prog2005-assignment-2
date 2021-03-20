package corona

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

// ServerError describes an internal server error and what http status code it should return.
type ServerError struct {
	error string
	// StatusCode is the http status code that should be returned by the server when handling this error.
	StatusCode int
}

func (e *ServerError) Error() string {
	return e.error
}

// TimeAsString converts a timepoint to a string width format yyyy-mm-dd.
func TimeAsString(t time.Time) string {
	return fmt.Sprintf("%.4d-%.2d-%.2d", t.Year(), t.Month(), t.Day())
}

// LatestDateInDateFloatMap returns the latest date in a map where key = date (as strings with format "yyyy-mm-dd").
// The naming reflects the stupidity of go's type system not being able to express this function generically.
func LatestDateInDateFloatMap(m *map[string]float64) string {
	// Get the keys in the map
	keys := make([]string, 0, len(*m))
	for k := range *m {
		keys = append(keys, k)
	}
	// Sort them alphabetically
	sort.Strings(keys)
	// Pick the last one
	latest := keys[len(keys)-1]
	return latest
}

// LatestDateInDateMapStringDataMap returns the latest date in a map where key = date (as strings with format "yyyy-mm-dd").
// The naming reflects the stupidity of go's type system not being able to express this function generically.
// YES THIS IS THE SAME FUNCTION TWICE. FUCK GO THAT'S WHY.
func LatestDateInDateMapStringDataMap(m *map[string]map[string]covidTrackerAPICountryData) string {
	// Get the keys in the map
	keys := make([]string, 0, len(*m))
	for k := range *m {
		keys = append(keys, k)
	}
	// Sort them alphabetically
	sort.Strings(keys)
	// Pick the last one
	latest := keys[len(keys)-1]
	return latest
}

// ParseScope query into two dates, or an error.
func ParseScope(qs *url.URL) (*time.Time, *time.Time, error) {
	scope := qs.Query().Get("scope")

	// If the query was not given, return nil, but no error
	if scope == "" {
		return nil, nil, nil
	}

	parts := strings.Split(scope, "-")

	// Check if all the parts are present
	if len(parts) != 6 {
		err := errors.New("incorrect date format in scope query")
		return nil, nil, err
	}

	// Parse into timepoints
	upper, err := time.Parse(time.RFC3339, fmt.Sprintf("%s-%s-%sT00:00:00Z", parts[0], parts[1], parts[2]))
	if err != nil {
		return nil, nil, err
	}

	lower, err := time.Parse(time.RFC3339, fmt.Sprintf("%s-%s-%sT00:00:00Z", parts[3], parts[4], parts[5]))

	return &upper, &lower, err
}
