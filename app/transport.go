package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
)

type countCreateRequest struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Context context.Context
}

type countIncrementRequest struct {
	UUID    string `json:"uuid"`
	Context context.Context
}

type createResponse struct {
	Counter Counter `json:"counter"`
	Error   string  `json:"err,omitempty"`
}

func makeCountCreateEndpoint(svc CounterService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(countCreateRequest)
		counter, err := svc.Create(req.Context, req.UUID, req.Name)
		if err != nil {
			return createResponse{
				Counter: counter,
				Error:   err.Error(),
			}, nil
		}
		return createResponse{
			Counter: counter,
			Error:   "",
		}, nil
	}
}

func makeCountIncrementEndpoint(svc CounterService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(countIncrementRequest)
		counter, err := svc.Increment(req.Context, req.UUID)
		if err != nil {
			return createResponse{
				Counter: counter,
				Error:   err.Error(),
			}, nil
		}
		return createResponse{
			Counter: counter,
			Error:   "",
		}, nil
	}
}

func decodeCountCreateRequest(c context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	return countCreateRequest{
		UUID:    vars["uuid"],
		Name:    vars["name"],
		Context: c,
	}, nil
}

func decodeCountIncrementRequest(c context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	return countIncrementRequest{
		UUID:    vars["uuid"],
		Context: c,
	}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
