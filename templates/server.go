package templates

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	gt "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	pb "github.com/ncostamagna/axul_template/templatespb"
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

type gRPCServer struct {
	getTemplate gt.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC StatsServiceServer.
func NewGRPCServer(ctx context.Context, endpoints Endpoints) pb.TemplatesServiceServer {
	return &gRPCServer{
		getTemplate: gt.NewServer(
			endpoints.Get,
			decodeGetRequest,
			encodeGetResponse,
		),
	}
}

func (s *gRPCServer) GetTemplate(ctx context.Context, req *pb.TemplateRequest) (*pb.Template, error) {
	_, resp, err := s.getTemplate.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.Template), nil
}

func decodeGetRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.TemplateRequest)
	return getRequest{id: uint(req.Id)}, nil
}

func encodeGetResponse(_ context.Context, resp interface{}) (interface{}, error) {
	r := resp.(response.Response)
	d := r.GetData()

	if d == nil {
		return nil, errors.New("Entity doesn't exists")
	}

	entity := d.(Template)
	template := &pb.Template{
		ID:       uint32(entity.ID),
		Type:     entity.Type,
		Template: entity.Template,
	}

	return template, nil
}
