package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type resolveMediaInput struct {
	URL string `query:"url" required:"true" doc:"External media page URL"`
}

type resolveMediaOutput struct {
	Body struct {
		URL  string `json:"url"`
		Type string `json:"type"`
	}
}

func (h *Handler) ResolveMedia(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "resolve-media",
			Method:      http.MethodGet,
			Path:        "/media/resolve",
			Summary:     "Resolve external media",
			Description: "Resolves external media page URL to a direct media file URL.",
			Tags:        []string{"Media"},
		},
		func(ctx context.Context, req *resolveMediaInput) (*resolveMediaOutput, error) {
			if req.URL == "" {
				return nil, huma.Error400BadRequest("url is required")
			}

			resolvedURL, mediaType, err := h.service.ResolveExternalMedia(ctx, req.URL)
			if err != nil {
				return nil, huma.Error422UnprocessableEntity(err.Error())
			}

			out := &resolveMediaOutput{}
			out.Body.URL = resolvedURL
			out.Body.Type = mediaType
			return out, nil
		},
	)
}
