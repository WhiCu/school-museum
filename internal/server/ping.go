package server

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type pingInput struct {
	Messsage string `query:"message" example:"ping" description:"Request message" required:"true"`
}
type pingOutput struct {
	Body struct {
		Message string `json:"message" example:"pong" description:"Response message"`
	}
}

func pingHandler(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "ping",
			Method:      http.MethodGet,
			Path:        "/ping",
			Summary:     "Ping",
			Description: "A simple endpoint to check if the server is running.",
			Tags: []string{
				"Test",
			},
		},
		func(ctx context.Context, req *pingInput) (res *pingOutput, err error) {
			if req.Messsage != "ping" {
				return nil, huma.Error400BadRequest("Invalid ping message")
			}
			res = &pingOutput{}
			res.Body.Message = "pong"
			return res, nil
		},
	)
}
