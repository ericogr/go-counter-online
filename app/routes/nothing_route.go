package routes

import (
	"fmt"
	"net/http"
)

func DoNothingRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hi there!\n")
	fmt.Fprintf(w, "- try to call GET /count/{uuidv5} to validate uuid\n")
	fmt.Fprintf(w, "- try to call POST /count/{uuidv5}/{name} to create counter\n")
}
