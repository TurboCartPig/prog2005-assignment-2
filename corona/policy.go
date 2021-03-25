package corona

import "net/http"

type policyResponse struct {
	Country    string  `json:"country"`
	Scope      string  `json:"scope"`
	Stringency float64 `json:"stringency"`
	Trend      float64 `json:"trend"`
}

// covidTrackerAPICountryData contains stringency data for a single country at a single date.
// According to the docs Stringency takes the value of StringencyActual if the latter is available,
// so that should satisfy spec requirements.
type covidTrackerAPICountryData struct {
	// Date in yyyy-mm-dd format.
	Date string `json:"date_value"`
	// CountryCode in alpha3Code format.
	CountryCode string `json:"country_code"`
	// Stringency value for the given country at the given date.
	Stringency float64 `json:"stringency"`
}

type covidTrackerApiResponse struct {
	// All the countries with data for the given scope, in alpha3Code format.
	Countries []string `json:"countries"`
	// Date -> Country -> Country Data
	Data map[string]map[string]covidTrackerAPICountryData `json:"data"`
}

func PolicyHandler(w http.ResponseWriter, r *http.Request) {

}
