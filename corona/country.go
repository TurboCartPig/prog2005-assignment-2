package corona

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// countryResponse is the response object from the country endpoint.
type countryResponse struct {
	Country              string  `json:"country"`
	Continent            string  `json:"continent"`
	Scope                string  `json:"scope"`
	Confirmed            float64 `json:"confirmed"`
	Recovered            float64 `json:"recovered"`
	PopulationPercentage float64 `json:"population_percentage"`
}

// caseHistory for one country as reported by mmediagroup.
// Some fields are omitted because we don't need them.
type caseHistory struct {
	Country    string             `json:"country"`
	Continent  string             `json:"continent"`
	Population float64            `json:"population"`
	Dates      map[string]float64 `json:"dates"`
}

// Count all the cases within a scope in time.
func (cases *caseHistory) countInScope(upper, lower time.Time) float64 {
	start := TimeAsString(upper)
	end := TimeAsString(lower)

	return cases.Dates[end] - cases.Dates[start]
}

// latestCount gets the latest count of cases.
func (cases *caseHistory) latestCount() float64 {
	key := LatestDateInDateFloatMap(&cases.Dates)
	return cases.Dates[key]
}

func getCases(country string) (confirmed, recovered caseHistory, err *ServerError) {
	root := MMediaGroupAPIRootPath + "/history?country=" + country

	cases := make(map[string]caseHistory)

	// Get confirmed cases
	res, geterr := http.Get(root + "&status=Confirmed")
	if geterr != nil {
		err = &ServerError{"Failed to get cases for country", res.StatusCode}
	}

	decerr := json.NewDecoder(res.Body).Decode(&cases)
	if decerr != nil {
		err = &ServerError{"Failed to decode response from remote", http.StatusInternalServerError}
	}
	res.Body.Close()

	confirmed = cases["All"]

	// Get recovered cases
	res, geterr = http.Get(root + "&status=Recovered")
	if geterr != nil {
		err = &ServerError{"Failed to get cases for country", res.StatusCode}
	}

	decerr = json.NewDecoder(res.Body).Decode(&cases)
	if decerr != nil {
		err = &ServerError{"Failed to decode response from remote", http.StatusInternalServerError}
	}
	res.Body.Close()

	recovered = cases["All"]

	return
}

// CountryHandler is the handler for the country endpoint.
// TODO: Handle no scope query
func CountryHandler(rw http.ResponseWriter, r *http.Request) {
	var response countryResponse
	var scoped bool

	country := chi.URLParam(r, "country")
	// Parse the scope query into two dates
	upper, lower, err := ParseScope(r.URL)
	if err != nil {
		log.Printf("Invalid request received: %s", err.Error())
		http.Error(rw, "Bad request: check the scope query.", http.StatusBadRequest)
	}
	if upper == nil {
		scoped = false
	} else {
		scoped = true
	}

	confirmed, recovered, serverErr := getCases(country)
	if serverErr != nil {
		http.Error(rw, serverErr.Error(), serverErr.StatusCode)
	}

	response.Country = confirmed.Country
	response.Continent = confirmed.Continent

	if scoped {
		response.Scope = TimeAsString(*upper) + "-" + TimeAsString(*lower)
		response.Confirmed = confirmed.countInScope(*upper, *lower)
		response.Recovered = recovered.countInScope(*upper, *lower)
	} else {
		response.Scope = "total"
		response.Confirmed = confirmed.latestCount()
		response.Recovered = recovered.latestCount()
	}

	response.PopulationPercentage = math.Round(response.Confirmed/confirmed.Population*100) / 100

	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		log.Printf("Something went wrong: %s", err.Error())
		http.Error(rw, "Something went wrong", http.StatusInternalServerError)
	}
}
