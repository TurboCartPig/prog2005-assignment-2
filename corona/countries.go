package corona

import (
	"encoding/json"
	"net/http"
)

// country represents a country as given by `restcountries.eu`.
type country struct {
	Name       string `json:"name"`
	Alpha3Code string `json:"alpha3Code"`
}

// GetCountryCode gets the alpha 3 code of a given country.
func GetCountryCode(name string) (string, *ServerError) {
	var country country

	res, err := http.Get(RestCountriesRootPath + "/name/" + name + "?fullText=true&fields=name,alpha3Code")
	if err != nil {
		return "", &ServerError{"Get country failed with: " + err.Error(), res.StatusCode}
	}

	err = json.NewDecoder(res.Body).Decode(&country)
	if err != nil {
		return "", &ServerError{"Failed to decode json response from restcountries.eu", http.StatusInternalServerError}
	}
	res.Body.Close()

	return country.Alpha3Code, nil
}
