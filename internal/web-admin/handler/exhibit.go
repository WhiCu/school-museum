package handler

import (
	"context"
	"net/http"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// --- Exhibits ---

// CreateExhibit - создание нового экспоната.
type createExhibitInput struct {
	Body struct {
		ExhibitionID uuid.UUID `json:"exhibition_id" format:"uuid" doc:"ID экспозиции"`
		Title        string    `json:"title" minLength:"1" doc:"Название экспоната"`
		Description  string    `json:"description" doc:"Описание экспоната"`
		ImageURLs    []string  `json:"image_urls" doc:"URLs изображений экспоната"`
	}
}

type createExhibitOutput struct {
	Body model.Exhibit
}

func (h *Handler) CreateExhibit(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "create-exhibit",
			Method:      http.MethodPost,
			Path:        "/exhibits",
			Summary:     "Создать экспонат",
			Description: "Создаёт новый экспонат в указанной экспозиции.",
			Tags: []string{
				"Admin", "Exhibits",
			},
		},
		func(ctx context.Context, req *createExhibitInput) (*createExhibitOutput, error) {
			if req.Body.ExhibitionID == uuid.Nil {
				return nil, huma.Error400BadRequest("exhibition_id обязателен")
			}
			ex, err := h.service.CreateExhibit(ctx, model.Exhibit{
				ExhibitionID: req.Body.ExhibitionID,
				Title:        req.Body.Title,
				Description:  req.Body.Description,
				ImageURLs:    req.Body.ImageURLs,
			})
			if err != nil {
				return nil, huma.Error500InternalServerError("не удалось создать экспонат")
			}
			return &createExhibitOutput{Body: ex}, nil
		},
	)
}

// UpdateExhibit - обновление экспоната.
type updateExhibitInput struct {
	ID   uuid.UUID `path:"id" format:"uuid" doc:"ID экспоната"`
	Body struct {
		Title       string   `json:"title" doc:"Название экспоната"`
		Description string   `json:"description" doc:"Описание экспоната"`
		ImageURLs   []string `json:"image_urls" doc:"URLs изображений экспоната"`
	}
}

type updateExhibitOutput struct {
	Body model.Exhibit
}

func (h *Handler) UpdateExhibit(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "update-exhibit",
			Method:      http.MethodPut,
			Path:        "/exhibits/{id}",
			Summary:     "Обновить экспонат",
			Description: "Обновляет данные существующего экспоната.",
			Tags: []string{
				"Admin",
				"Exhibits",
			},
		},
		func(ctx context.Context, req *updateExhibitInput) (*updateExhibitOutput, error) {
			ex, err := h.service.UpdateExhibit(ctx, model.Exhibit{
				ID:          req.ID,
				Title:       req.Body.Title,
				Description: req.Body.Description,
				ImageURLs:   req.Body.ImageURLs,
			})
			if err != nil {
				return nil, huma.Error500InternalServerError("не удалось обновить экспонат")
			}
			return &updateExhibitOutput{Body: ex}, nil
		},
	)
}

// DeleteExhibit - удаление экспоната.
type deleteExhibitInput struct {
	ID uuid.UUID `path:"id" format:"uuid" doc:"ID экспоната"`
}

func (h *Handler) DeleteExhibit(api huma.API) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "delete-exhibit",
			Method:      http.MethodDelete,
			Path:        "/exhibits/{id}",
			Summary:     "Удалить экспонат",
			Description: "Удаляет экспонат по его идентификатору.",
			Tags:        []string{"Admin", "Exhibits"},
		},
		func(ctx context.Context, req *deleteExhibitInput) (*struct{}, error) {
			if err := h.service.DeleteExhibit(ctx, req.ID); err != nil {
				return nil, huma.Error500InternalServerError("не удалось удалить экспонат")
			}
			return nil, nil
		},
	)
}
