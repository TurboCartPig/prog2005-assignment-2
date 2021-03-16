package corona

const (
	MMediaGroupAPIRootPath  string = "https://covid-api.mmediagroup.fr/v1"
	CovidTrackerAPIRootPath string = "https://covidtrackerapi.bsg.ox.ac.uk/api/v2"

	// Version of the api
	Version string = "v1"

	// RootPath common for all the endpoints
	RootPath string = "/corona/" + Version

	// CountryRootPath for the country endpoint
	CountryRootPath string = RootPath + "/country"

	// PolicyRootPath for the policy endpoint
	PolicyRootPath string = RootPath + "/policy"

	// DiagRootPath for the diag endpoint
	DiagRootPath string = RootPath + "/diag"

	// NotificationsRootPath for the notifications endpoint
	NotificationsRootPath string = RootPath + "/notifications"
)
