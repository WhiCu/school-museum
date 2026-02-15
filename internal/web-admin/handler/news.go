package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/WhiCu/school-museum/db/storage"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// CreateNews - создание новой новости.
type createNewsInput struct {
	Body struct {
		Title   string `json:"title" minLength:"1" doc:"Заголовок новости"`
		Content string `json:"content" doc:"Содержание новости"`
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
			n := h.service.CreateNews(req.Body.Title, req.Body.Content)
			return &createNewsOutput{Body: n}, nil
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
			if err := h.service.DeleteNews(req.ID); err != nil {
				if errors.Is(err, storage.ErrNotFound) {
					return nil, huma.Error404NotFound("новость не найдена")
				}
				return nil, huma.Error500InternalServerError("внутренняя ошибка")
			}
			return nil, nil
		},
	)
}
