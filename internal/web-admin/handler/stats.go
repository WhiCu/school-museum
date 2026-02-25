package handler

import (
	"context"
	"net/http"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/danielgtaylor/huma/v2"
)

// GetStats — returns aggregated visit statistics for the admin panel.
type getStatsOutput struct {
	Body model.VisitStats `json:"stats"`
}

func (h *Handler) GetStats(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "admin-get-stats",
			Method:      http.MethodGet,
			Path:        "/stats",
			Summary:     "Получить статистику",
			Description: "Возвращает статистику посещений и количество объектов.",
			Tags:        []string{"Admin", "Stats"},
		},
		func(ctx context.Context, req *struct{}) (*getStatsOutput, error) {
			stats, err := h.service.GetStats(ctx)
			if err != nil {
				return nil, huma.Error500InternalServerError("не удалось получить статистику")
			}
			return &getStatsOutput{Body: stats}, nil
		},
	)
}
