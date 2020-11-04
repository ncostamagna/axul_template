package templates

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/ncostamagna/response"
)

//Endpoints struct
type Endpoints struct {
	Create endpoint.Endpoint
	Update endpoint.Endpoint
	Get    endpoint.Endpoint
	GetAll endpoint.Endpoint
}

//MakeEndpoints handler endpoints
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(TemplateRequest)

		t := Template{
			Type:     req.Type,
			Template: req.Template,
		}

		fmt.Println("Pasa")
		if rerr := s.Create(ctx, &t); rerr != nil {
			resp := response.NewResponse(rerr.Message(), rerr.Status(), "", nil)
			return resp, nil
		}

		resp := response.NewResponse("Success", 200, "", t)
		return resp, nil

	}
}

func makeGetAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(getRequest)
		var templates []Template
		fmt.Println(req)

		if rerr := s.GetAll(ctx, &templates); rerr != nil {
			resp := response.NewResponse(rerr.Message(), rerr.Status(), "", nil)
			return resp, nil
		}

		return response.NewResponse("Success", 200, "", templates), nil
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}
