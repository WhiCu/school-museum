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

// --- Exhibitions ---

// CreateExhibition - создание новой экспозиции.
type createExhibitionInput struct {
	Body struct {
		Title       string `json:"title" minLength:"1" doc:"Название экспозиции"`
		Description string `json:"description" doc:"Описание экспозиции"`
	}
}

type createExhibitionOutput struct {
	Body model.Exhibition
}

func (h *Handler) CreateExhibition(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "create-exhibition",
			Method:      http.MethodPost,
			Path:        "/exhibitions",
			Summary:     "Создать экспозицию",
			Description: "Создаёт новую экспозицию в музее.",
			Tags:        []string{"Admin", "Exhibitions"},
		},
		func(ctx context.Context, req *createExhibitionInput) (*createExhibitionOutput, error) {
			ex := h.service.CreateExhibition(req.Body.Title, req.Body.Description)
			return &createExhibitionOutput{Body: ex}, nil
		},
	)
}

// UpdateExhibition - обновление экспозиции.
type updateExhibitionInput struct {
	ID   uuid.UUID `path:"id" format:"uuid" doc:"ID экспозиции"`
	Body struct {
		Title       string `json:"title" doc:"Название экспозиции"`
		Description string `json:"description" doc:"Описание экспозиции"`
	}
}

type updateExhibitionOutput struct {
	Body model.Exhibition
}

func (h *Handler) UpdateExhibition(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "update-exhibition",
			Method:      http.MethodPut,
			Path:        "/exhibitions/{id}",
			Summary:     "Обновить экспозицию",
			Description: "Обновляет данные существующей экспозиции.",
			Tags: []string{
				"Admin",
				"Exhibitions",
			},
		},
		func(ctx context.Context, req *updateExhibitionInput) (*updateExhibitionOutput, error) {
			ex, err := h.service.UpdateExhibition(req.ID, req.Body.Title, req.Body.Description)
			if err != nil {
				if errors.Is(err, storage.ErrNotFound) {
					return nil, huma.Error404NotFound("экспозиция не найдена")
				}
				return nil, huma.Error500InternalServerError("внутренняя ошибка")
			}
			return &updateExhibitionOutput{Body: ex}, nil
		},
	)
}

// DeleteExhibition - удаление экспозиции (и всех её экспонатов).
type deleteExhibitionInput struct {
	ID uuid.UUID `path:"id" format:"uuid" doc:"ID экспозиции"`
}

func (h *Handler) DeleteExhibition(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "delete-exhibition",
			Method:      http.MethodDelete,
			Path:        "/exhibitions/{id}",
			Summary:     "Удалить экспозицию",
			Description: "Удаляет экспозицию и все связанные с ней экспонаты.",
			Tags:        []string{"Admin", "Exhibitions"},
		},
		func(ctx context.Context, req *deleteExhibitionInput) (*struct{}, error) {
			if err := h.service.DeleteExhibition(req.ID); err != nil {
				if errors.Is(err, storage.ErrNotFound) {
					return nil, huma.Error404NotFound("экспозиция не найдена")
				}
				return nil, huma.Error500InternalServerError("внутренняя ошибка")
			}
			return nil, nil
		},
	)
}
