package handler

import (
	"context"
	"net/http"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// GetAllExhibitions - получение списка всех экспозиций.
type getAllExhibitionsOutput struct {
	Body []model.Exhibition `json:"exhibitions"`
}

func (h *Handler) GetAllExhibitions(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "get-all-exhibitions",
			Method:      http.MethodGet,
			Path:        "/exhibitions",
			Summary:     "Получить все экспозиции",
			Description: "Возвращает список всех экспозиций музея.",
			Tags:        []string{"Exhibitions"},
		},
		func(ctx context.Context, req *struct{}) (*getAllExhibitionsOutput, error) {
			exhibitions := h.service.GetAllExhibitions()
			return &getAllExhibitionsOutput{Body: exhibitions}, nil
		},
	)
}

// GetExhibitionByID - получение экспозиции с её экспонатами по ID.
type getExhibitionByIDInput struct {
	ID uuid.UUID `path:"id" format:"uuid" doc:"ID экспозиции"`
}

type getExhibitionByIDOutput struct {
	Body model.Exhibition `json:"exhibition"`
}

func (h *Handler) GetExhibitionByID(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "get-exhibition-by-id",
			Method:      http.MethodGet,
			Path:        "/exhibitions/{id}",
			Summary:     "Получить экспозицию по ID",
			Description: "Возвращает конкретную экспозицию со всеми её экспонатами.",
			Tags:        []string{"Exhibitions"},
		},
		func(ctx context.Context, req *getExhibitionByIDInput) (*getExhibitionByIDOutput, error) {
			detail, ok := h.service.GetExhibitionByID(req.ID)
			if !ok {
				return nil, huma.Error404NotFound("экспозиция не найдена")
			}
			return &getExhibitionByIDOutput{Body: detail}, nil
		},
	)
}
