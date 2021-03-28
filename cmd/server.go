package main

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"assignment-2/corona"
	"assignment-2/notifications"
	mymw "assignment-2/middleware"
	fs "assignment-2/firestore"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Globals
// -------------------------------------------------------------------------------------------

// StartTime is instant the server was started
var StartTime time.Time = time.Now()

// DefaultPort is the default port number if no other port number is specified via the $PORT environment variable
const DefaultPort int = 3000

// Functions
// -------------------------------------------------------------------------------------------

// Get the port from environment variable $PORT, or use default if the variable is not set
func port() int {
	if port := os.Getenv("PORT"); port != "" {
		port, _ := strconv.Atoi(port)
		return port
	}
	return DefaultPort
}

// Serve the resources as defined by routes in `r`
func serve(r *chi.Mux, wg *sync.WaitGroup) {
	port := port()
	addr := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalf("Error while serving http: %s", err.Error())
	}

	wg.Done()
}

// Setup all the top level routes the server serves on
func setupRoutes(fs *firestore.Client, registerChan chan<- string) *chi.Mux {
	r := chi.NewRouter()

	// Use middleware
	r.Use(middleware.Logger)
	r.Use(middleware.RedirectSlashes)
	r.Use(mymw.ReturnJSON)

	// Define endpoints
	r.Get(corona.DiagRootPath, corona.NewDiagHandler(fs, StartTime))
	r.Get(corona.CountryRootPath+"/{country:[a-zA-Z]+}", corona.CountryHandler)
	r.Get(corona.PolicyRootPath+"/{country:[a-zA-Z]+}", corona.PolicyHandler)

	// Define webhook endpoints in a subroute
	r.Route(notifications.RootPath, func(r chi.Router) {
		r.Post("/", notifications.NewCreateHandler(fs, registerChan))
		r.Get("/", notifications.NewReadAllHandler(fs))
		r.Delete(notifications.IDPattern, notifications.NewDeleteHandler(fs))
		r.Get(notifications.IDPattern, notifications.NewReadHandler(fs))
	})

	return r
}

func main() {
	// Initialize a firestore client
	fs := fs.NewFirestoreClient()
	defer fs.Close()

	registerChan := make(chan string)

	wg := &sync.WaitGroup{}
	wg.Add(2) //nolint:gomnd // How many goroutines we are about to launch

	r := setupRoutes(fs, registerChan)
	go serve(r, wg)
	go notifications.InvokeLoop(fs, registerChan, wg)

	wg.Wait()
}
