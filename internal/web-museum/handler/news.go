package handler

import (
	"context"
	"net/http"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// GetAllNews - получение списка всех новостей.
type getAllNewsOutput struct {
	Body []model.News `json:"news"`
}

func (h *Handler) GetAllNews(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "get-all-news",
			Method:      http.MethodGet,
			Path:        "/news",
			Summary:     "Получить все новости",
			Description: "Возвращает список всех новостей музея.",
			Tags:        []string{"News"},
		},
		func(ctx context.Context, req *struct{}) (*getAllNewsOutput, error) {
			news := h.service.GetAllNews()
			return &getAllNewsOutput{Body: news}, nil
		},
	)
}

// GetNewsByID - получение конкретной новости по ID.
type getNewsByIDInput struct {
	ID uuid.UUID `path:"id" format:"uuid" doc:"ID новости"`
}

type getNewsByIDOutput struct {
	Body model.News `json:"news"`
}

func (h *Handler) GetNewsByID(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "get-news-by-id",
			Method:      http.MethodGet,
			Path:        "/news/{id}",
			Summary:     "Получить новость по ID",
			Description: "Возвращает конкретную новость по её идентификатору.",
			Tags:        []string{"News"},
		},
		func(ctx context.Context, req *getNewsByIDInput) (*getNewsByIDOutput, error) {
			n, ok := h.service.GetNewsByID(req.ID)
			if !ok {
				return nil, huma.Error404NotFound("новость не найдена")
			}
			return &getNewsByIDOutput{Body: n}, nil
		},
	)
}
