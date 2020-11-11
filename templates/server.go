package templates

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/ncostamagna/response"
)

type (
	entityRequest struct {
		ID       uint   `json:"id"`
		Type     string `json:"type"`
		Template string `json:"template"`
	}

	getRequest struct {
		id uint
	}
)

//NewHTTPServer is a server handler
func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Handle("/templates/", httptransport.NewServer(
		endpoints.Create,
		decodeCreate,
		encodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/templates/", httptransport.NewServer(
		endpoints.GetAll,
		decodeGet,
		encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/templates/{id}", httptransport.NewServer(
		endpoints.Get,
		decodeGet,
		encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/templates/{id}", httptransport.NewServer(
		endpoints.Update,
		nil,
		encodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/templates/{id}", httptransport.NewServer(
		nil,
		decodeCreate,
		encodeResponse,
		opts...,
	)).Methods("DELETE")

	return r

}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.WriteHeader(r.GetStatusCode())
	return json.NewEncoder(w).Encode(r)
}

func decodeCreate(ctx context.Context, r *http.Request) (interface{}, error) {

	var req entityRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGet(ctx context.Context, r *http.Request) (interface{}, error) {

	vars := mux.Vars(r)
	idVar, ok := vars["id"]
	req := getRequest{}

	if !ok {
		return req, nil
	}

	id, err := strconv.Atoi(idVar)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, errors.New("Invalid id")
	}

	req.id = uint(id)
	return req, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := response.NewResponse(err.Error(), 500, "", nil)
	w.WriteHeader(resp.GetStatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
