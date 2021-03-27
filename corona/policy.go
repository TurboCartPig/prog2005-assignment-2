package corona

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type PolicyResponse struct {
	Country    string  `json:"country"`
	Scope      string  `json:"scope"`
	Stringency float64 `json:"stringency"`
	Trend      float64 `json:"trend"`
}

// covidTrackerAPIStringencyData contains stringency data for a single country at a single date.
// According to the docs Stringency takes the value of StringencyActual if the latter is available,
// so that should satisfy spec requirements.
type covidTrackerAPIStringencyData struct {
	// Date in yyyy-mm-dd format.
	Date string `json:"date_value"`
	// CountryCode in alpha3Code format.
	CountryCode string `json:"country_code"`
	// Stringency value for the given country at the given date.
	Stringency float64 `json:"stringency"`
}

// covidTrackerApiResponse is the response from http requests to the CovidTrackerApi.
type covidTrackerAPIResponse struct {
	Data covidTrackerAPIStringencyData `json:"stringencyData"`
}

// getStringency for a given country's alpha3 code at a given date.
func getStringency(code, date string) (response covidTrackerAPIResponse, err *ServerError) {
	res, geterr := http.Get(CovidTrackerAPIRootPath + "/stringency/actions/" + code + "/" + date)
	if geterr != nil {
		err = &ServerError{"Failed to get cases for country", res.StatusCode}
		return
	}

	// Set default value if not available
	response.Data.Stringency = -1

	decerr := json.NewDecoder(res.Body).Decode(&response)
	if decerr != nil {
		err = &ServerError{"Failed to decode response from remote", http.StatusInternalServerError}
		return
	}
	res.Body.Close()

	return
}

func PolicyHandler(rw http.ResponseWriter, r *http.Request) {
	var response PolicyResponse
	var scoped bool

	country := chi.URLParam(r, "country")
	upper, lower, err := ParseScope(r.URL)
	if err != nil {
		log.Printf("Invalid request received: %s", err.Error())
		http.Error(rw, "Bad request: check the scope query.", http.StatusBadRequest)
		return
	}
	if upper == nil {
		scoped = false
	} else {
		scoped = true
	}

	alpha3code, serverErr := GetCountryCode(country)
	if serverErr != nil {
		http.Error(rw, serverErr.Error(), serverErr.StatusCode)
		return
	}

	// If scope query was passed, fetch data for all the dates in range
	if scoped {
		// Fetch stringency info for the two dates
		upperRes, serverErr := getStringency(alpha3code, TimeAsString(*upper))
		if serverErr != nil {
			http.Error(rw, serverErr.Error(), serverErr.StatusCode)
			return
		}
		lowerRes, serverErr := getStringency(alpha3code, TimeAsString(*lower))
		if serverErr != nil {
			http.Error(rw, serverErr.Error(), serverErr.StatusCode)
			return
		}

		// Fill out response data
		response.Country = country
		response.Scope = TimeAsString(*upper) + "-" + TimeAsString(*lower)

		// Take current stringency from latest of the two
		response.Stringency = lowerRes.Data.Stringency

		// Trend is difference in stringency from first date to last data in scope
		response.Trend = lowerRes.Data.Stringency - upperRes.Data.Stringency
	} else { // Fetch data for latest available date
		res, serverErr := getStringency(alpha3code, TimeAsString(time.Now().AddDate(0, 0, -2)))
		if serverErr != nil {
			http.Error(rw, serverErr.Error(), serverErr.StatusCode)
			return
		}

		// Fill out response data
		response.Country = country
		response.Scope = "total"
		response.Stringency = res.Data.Stringency
		response.Trend = 0
	}

	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		log.Printf("Something went wrong: %s", err.Error())
		http.Error(rw, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
