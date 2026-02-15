package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type pingInput struct {
	Messsage string `query:"message" example:"ping" description:"Request message" required:"false"`
}
type pingOutput struct {
	Body struct {
		Message string `json:"message" example:"pong" description:"Response message"`
	}
}

func (h *Handler) Ping(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "museum-ping",
			Method:      http.MethodGet,
			Path:        "/ping",
			Summary:     "Ping",
			Description: "A simple endpoint to check if the server is running.",
			Tags: []string{
				"Test",
			},
		},
		func(ctx context.Context, req *pingInput) (res *pingOutput, err error) {
			h.log.Debug("ping request received", slog.String("message", req.Messsage))
			res = &pingOutput{}
			res.Body.Message = "pong"
			return res, nil
		},
	)
}
