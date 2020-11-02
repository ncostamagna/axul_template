package contacts

import (
	"context"
	"fmt"
	"time"

	"github.com/ncostamagna/rerrors"

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
		req := request.(ContactRequest)
		birthday, err := time.Parse(layoutISO, fmt.Sprintf("%s 17:00:00", req.Birthday))

		if err != nil {
			rerr := rerrors.NewBadRequestError(err)
			resp := response.NewResponse(rerr.Message(), rerr.Status(), "", nil)
			return resp, nil
		}

		c := Contact{
			Firstname: req.Firstname,
			Lastname:  req.Lastname,
			Nickname:  req.Nickname,
			Gender:    req.Gender,
			Phone:     req.Phone,
			Birthday:  birthday,
		}

		fmt.Println(birthday)

		if rerr := s.Create(ctx, &c); rerr != nil {
			resp := response.NewResponse(rerr.Message(), rerr.Status(), "", nil)
			return resp, nil
		}

		resp := response.NewResponse("Success", 200, "", c)
		return resp, nil

	}
}

func makeGetAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(getRequest)
		var contacts []Contact
		fmt.Println(req)

		if rerr := s.GetAll(ctx, &contacts, req.birthday); rerr != nil {
			resp := response.NewResponse(rerr.Message(), rerr.Status(), "", nil)
			return resp, nil
		}

		return response.NewResponse("Success", 200, "", contacts), nil
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
