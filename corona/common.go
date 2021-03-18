package corona

import (
	"errors"
	"fmt"
	"net/url"
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

// ParseScope query into two dates, or an error.
func ParseScope(qs *url.URL) (upper, lower time.Time, err error) {
	scope := qs.Query().Get("scope")

	// If the query was not given, return some sensible default
	if scope == "" {
		upper = time.Unix(0, 0)
		lower = time.Now()
		return
	}

	parts := strings.Split(scope, "-")

	// Check if all the parts are present
	if len(parts) != 6 {
		err = errors.New("Incorrect date format in scope query")
		return
	}

	// Parse into timepoints
	upper, err = time.Parse(time.RFC3339, fmt.Sprintf("%s-%s-%sT00:00:00Z", parts[0], parts[1], parts[2]))
	lower, err = time.Parse(time.RFC3339, fmt.Sprintf("%s-%s-%sT00:00:00Z", parts[3], parts[4], parts[5]))

	return
}
