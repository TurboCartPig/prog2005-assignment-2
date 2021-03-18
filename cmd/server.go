package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"assignment-2/corona"
	mymw "assignment-2/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Globals
/////////////////////////////////////////////////////////////////////////////////////////////

// The instant the server was started
var StartTime time.Time = time.Now()

// Functions
/////////////////////////////////////////////////////////////////////////////////////////////

// Get the port from environment variable $PORT, or use default if the variable is not set
func port() int {
	if port := os.Getenv("PORT"); port != "" {
		port, _ := strconv.Atoi(port)
		return port
	}
	return 3000
}

// Serve the resources as defined by routes in `r`
func serve(r *chi.Mux) {
	port := port()
	addr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(addr, r)
}

// Setup all the top level routes the server serves on
func setupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Use middleware
	r.Use(middleware.Logger)
	r.Use(middleware.RedirectSlashes)
	r.Use(mymw.ReturnJSON)

	// Define endpoints
	r.Get(corona.DiagRootPath, corona.NewDiagHandler(0, StartTime))
	r.Get(corona.CountryRootPath+"/{country:[a-zA-Z]+}", corona.CountryHandler)

	return r
}

func main() {
	r := setupRoutes()
	serve(r)
}
