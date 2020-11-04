package templates

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/ncostamagna/response"
)

type (
	TemplateRequest struct {
		ID       uint   `json:"id"`
		Type     string `json:"type"`
		Template string `json:"template"`
	}

	getRequest struct {
		id       string
		birthday string
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
		decodeCreateContact,
		encodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/templates/", httptransport.NewServer(
		endpoints.GetAll,
		decodeGetContact,
		encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/templates/{id}", httptransport.NewServer(
		endpoints.Get,
		nil,
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
		decodeCreateContact,
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

func decodeCreateContact(ctx context.Context, r *http.Request) (interface{}, error) {

	var req TemplateRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetContact(ctx context.Context, r *http.Request) (interface{}, error) {

	v := r.URL.Query()

	req := getRequest{
		birthday: v.Get("birthday"),
	}
	return req, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := response.NewResponse(err.Error(), 500, "", nil)
	w.WriteHeader(resp.GetStatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
