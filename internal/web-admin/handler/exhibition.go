package handler

import (
	"context"
	"net/http"

	"github.com/WhiCu/school-museum/db/model"
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
			ex, err := h.service.CreateExhibition(ctx, model.Exhibition{
				Title:       req.Body.Title,
				Description: req.Body.Description,
			})
			if err != nil {
				return nil, huma.Error500InternalServerError("не удалось создать экспозицию")
			}
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
			ex, err := h.service.UpdateExhibition(ctx, model.Exhibition{
				ID:          req.ID,
				Title:       req.Body.Title,
				Description: req.Body.Description,
			})
			if err != nil {
				return nil, huma.Error500InternalServerError("не удалось обновить экспозицию")
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
			if err := h.service.DeleteExhibition(ctx, req.ID); err != nil {
				return nil, huma.Error500InternalServerError("не удалось удалить экспозицию")
			}
			return nil, nil
		},
	)
}

// SetExhibitionPreview - установка превью-экспоната для экспозиции.
type setExhibitionPreviewInput struct {
	ID   uuid.UUID `path:"id" format:"uuid" doc:"ID экспозиции"`
	Body struct {
		ExhibitID string `json:"exhibit_id" doc:"ID экспоната-превью (пустая строка для сброса)"`
	}
}

type setExhibitionPreviewOutput struct {
	Body model.Exhibition
}

func (h *Handler) SetExhibitionPreview(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "set-exhibition-preview",
			Method:      http.MethodPut,
			Path:        "/exhibitions/{id}/preview",
			Summary:     "Установить превью экспозиции",
			Description: "Устанавливает экспонат, изображение которого будет использоваться как превью экспозиции.",
			Tags:        []string{"Admin", "Exhibitions"},
		},
		func(ctx context.Context, req *setExhibitionPreviewInput) (*setExhibitionPreviewOutput, error) {
			var exhibitID *uuid.UUID
			if req.Body.ExhibitID != "" {
				parsed, err := uuid.Parse(req.Body.ExhibitID)
				if err != nil {
					return nil, huma.Error422UnprocessableEntity("неверный формат ID экспоната")
				}
				exhibitID = &parsed
			}
			ex, err := h.service.SetExhibitionPreview(ctx, req.ID, exhibitID)
			if err != nil {
				return nil, huma.Error500InternalServerError("не удалось установить превью")
			}
			return &setExhibitionPreviewOutput{Body: ex}, nil
		},
	)
}
