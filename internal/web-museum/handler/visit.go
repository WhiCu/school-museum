package handler

import (
	"context"
	"net/http"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/danielgtaylor/huma/v2"
)

// RecordVisit — tracks a page view from the public site.
// IP and User-Agent are injected into the context by the visitTrackingMiddleware.
// Client-side data (screen, language, referrer) comes from the request body.
type recordVisitInput struct {
	Body struct {
		Page         string `json:"page" example:"/" doc:"Page path"`
		Referrer     string `json:"referrer,omitempty" doc:"HTTP referrer"`
		ScreenWidth  int    `json:"screen_width,omitempty" doc:"Screen width in px"`
		ScreenHeight int    `json:"screen_height,omitempty" doc:"Screen height in px"`
		Language     string `json:"language,omitempty" doc:"Browser language"`
	}
}

type recordVisitOutput struct {
	Body struct {
		OK bool `json:"ok"`
	}
}

func (h *Handler) RecordVisit(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "record-visit",
			Method:      http.MethodPost,
			Path:        "/visit",
			Summary:     "Record a page visit",
			Description: "Tracks a page view. IP and User-Agent captured by middleware; client sends page, referrer, screen size, language.",
			Tags:        []string{"Visit"},
		},
		func(ctx context.Context, req *recordVisitInput) (*recordVisitOutput, error) {
			ip, _ := ctx.Value(model.CtxKeyVisitorIP).(string)
			ua, _ := ctx.Value(model.CtxKeyVisitorUA).(string)

			v := model.Visitor{
				IP:           ip,
				UserAgent:    ua,
				Page:         req.Body.Page,
				Referrer:     req.Body.Referrer,
				ScreenWidth:  req.Body.ScreenWidth,
				ScreenHeight: req.Body.ScreenHeight,
				Language:     req.Body.Language,
			}

			if err := h.service.RecordVisit(ctx, v); err != nil {
				h.log.Error("failed to record visit")
				return nil, huma.Error500InternalServerError("failed to record visit")
			}
			out := &recordVisitOutput{}
			out.Body.OK = true
			return out, nil
		},
	)
}
