package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ericogr/go-counter-online/core"
	"github.com/ericogr/go-counter-online/counter"
	"github.com/ericogr/go-counter-online/storage"
	"github.com/gorilla/mux"
)

// GetCountRoute is the route that handles the count requests using uuid v5 pattern (ex AAAAAAAA-AAAA-5AAA-AAAA-AAAAAAAAAAAA)
// ex: curl 'http://localhost:8080/count/AAAAAAAA-AAAA-5AAA-AAAA-AAAAAAAAAAAA'
func GetCountRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	counter := counter.Counter{UUID: vars["uuid"]}
	counterData, err := storage.GetCounterInstance(core.Options.Datastore, core.Options.ExtraParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %s", err)
		return
	}

	counter, err = counterData.Increment(counter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %s", err)
		return
	}
	json.NewEncoder(w).Encode(counter)
}

// CreateCountRoute is the route that handles the create requests using uuid v5 pattern
// ex: curl -X POST 'http://localhost:8080/count/AAAAAAAA-AAAA-5AAA-AAAA-AAAAAAAAAAAA/foobar'
func CreateCountRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var counter counter.Counter
	counter.UUID = vars["uuid"]
	counter.Name = vars["name"]
	counter.Date = time.Now()

	counterData, err := storage.GetCounterInstance(core.Options.Datastore, core.Options.ExtraParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %s", err)
		return
	}
	counter, err = counterData.Create(counter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %s", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
