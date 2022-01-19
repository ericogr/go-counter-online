package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ericogr/go-counter-online/core"
	"github.com/ericogr/go-counter-online/routes"
	"github.com/gorilla/mux"
)

var (
	PATH_ROOT       = "/"
	PATH_COUNT_GET  = "/count/{uuid:[0-9A-F]{8}-[0-9A-F]{4}-[5][0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}}"
	PATH_COUNT_POST = "/count/{uuid:[0-9A-F]{8}-[0-9A-F]{4}-[5][0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}}/{name:[a-zA-Z][a-zA-Z0-9]{0,15}}"
)

func init() {
	fmt.Printf("Initializing...\n")

	port := flag.Int("port", 8080, "app port to listen connections")
	datastore := flag.String("datastore", "memory", "datastore name (memory | postgresql)")
	extraParams := flag.String("extra-params", "", "extra parameters")
	hideExtraParams := flag.Bool("hide-extra-params", false, "hide extra parameters")
	flag.Parse()

	core.Options.Port = *port
	core.Options.Datastore = *datastore
	core.Options.ExtraParams = *extraParams
	core.Options.HideExtraParams = *hideExtraParams

	fmt.Printf(core.Options.String(core.Options.HideExtraParams))
}

func main() {
	fmt.Printf("Listening on port %d\n", core.Options.Port)
	doRoute()
}

func doRoute() {
	r := mux.NewRouter()
	r.HandleFunc(PATH_ROOT, routes.DoNothingRoute)
	r.HandleFunc(PATH_COUNT_GET, routes.GetCountRoute).Methods("GET")
	r.HandleFunc(PATH_COUNT_POST, routes.CreateCountRoute).Methods("POST")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", core.Options.Port), r))
}
