package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ericogr/go-counter-online/services"
	"github.com/gorilla/mux"
)

type CounterRest struct {
	CounterService services.CounterService
}

type Counter struct {
	UUID  string    `json:"uuid"`
	Name  string    `json:"name"`
	Count int       `json:"count"`
	Date  time.Time `json:"created_at"`
}

// GetCount is the route that handles the count requests using uuid v5 pattern (ex AAAAAAAA-AAAA-5AAA-AAAA-AAAAAAAAAAAA)
// ex: curl 'http://localhost:8080/count/AAAAAAAA-AAAA-5AAA-AAAA-AAAAAAAAAAAA'
func (cs *CounterRest) GetCount(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	counter, err := cs.CounterService.Increment(services.Counter{
		UUID: vars["uuid"],
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error: %s", err)
		return
	}
	json.NewEncoder(w).Encode(counter)
}

// CreateCount is the route that handles the create requests using uuid v5 pattern
// ex: curl -X POST 'http://localhost:8080/count/AAAAAAAA-AAAA-5AAA-AAAA-AAAAAAAAAAAA/foobar'
func (cs *CounterRest) CreateCount(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	_, err := cs.CounterService.Create(services.Counter{
		UUID: vars["uuid"],
		Name: vars["name"],
		Date: time.Now(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error: %s", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (cs *CounterRest) DoNothing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hi there!\n")
	fmt.Fprintf(w, "- try to call GET /count/{uuidv5} to validate uuid\n")
	fmt.Fprintf(w, "- try to call POST /count/{uuidv5}/{name} to create counter\n")
}
