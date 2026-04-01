package handler

import (
	"context"
	"net/http"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// CreateNews - создание новой новости.
type createNewsInput struct {
	Body struct {
		Title     string   `json:"title" minLength:"1" doc:"Заголовок новости"`
		Content   string   `json:"content" doc:"Содержание новости"`
		ImageURLs []string `json:"image_urls" doc:"URLs изображений новости"`
	}
}

type createNewsOutput struct {
	Body model.News
}

func (h *Handler) CreateNews(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "create-news",
			Method:      http.MethodPost,
			Path:        "/news",
			Summary:     "Создать новость",
			Description: "Создаёт новую новость в музее.",
			Tags:        []string{"Admin", "News"},
		},
		func(ctx context.Context, req *createNewsInput) (*createNewsOutput, error) {
			n, err := h.service.CreateNews(ctx, model.News{
				Title:     req.Body.Title,
				Content:   req.Body.Content,
				ImageURLs: req.Body.ImageURLs,
			})
			if err != nil {
				return nil, huma.Error500InternalServerError("не удалось создать новость")
			}
			return &createNewsOutput{Body: n}, nil
		},
	)
}

// UpdateNews - обновление новости.
type updateNewsInput struct {
	ID   uuid.UUID `path:"id" format:"uuid" doc:"ID новости"`
	Body struct {
		Title     string   `json:"title" doc:"Заголовок новости"`
		Content   string   `json:"content" doc:"Содержание новости"`
		ImageURLs []string `json:"image_urls" doc:"URLs изображений новости"`
	}
}

type updateNewsOutput struct {
	Body model.News
}

func (h *Handler) UpdateNews(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "update-news",
			Method:      http.MethodPut,
			Path:        "/news/{id}",
			Summary:     "Обновить новость",
			Description: "Обновляет данные существующей новости.",
			Tags:        []string{"Admin", "News"},
		},
		func(ctx context.Context, req *updateNewsInput) (*updateNewsOutput, error) {
			n, err := h.service.UpdateNews(ctx, model.News{
				ID:        req.ID,
				Title:     req.Body.Title,
				Content:   req.Body.Content,
				ImageURLs: req.Body.ImageURLs,
			})
			if err != nil {
				return nil, huma.Error500InternalServerError("не удалось обновить новость")
			}
			return &updateNewsOutput{Body: n}, nil
		},
	)
}

// DeleteNews - удаление новости по ID.
type deleteNewsInput struct {
	ID uuid.UUID `path:"id" format:"uuid" doc:"ID новости"`
}

func (h *Handler) DeleteNews(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "delete-news",
			Method:      http.MethodDelete,
			Path:        "/news/{id}",
			Summary:     "Удалить новость",
			Description: "Удаляет новость по её идентификатору.",
			Tags: []string{
				"Admin", "News",
			},
		},
		func(ctx context.Context, req *deleteNewsInput) (*struct{}, error) {
			if err := h.service.DeleteNews(ctx, req.ID); err != nil {
				return nil, huma.Error500InternalServerError("не удалось удалить новость")
			}
			return nil, nil
		},
	)
}
