package server

import (
	"context"
	"net/http"

	"github.com/WhiCu/school-museum/internal/telemetry"
	"github.com/danielgtaylor/huma/v2"
)

type analyticsOutput struct {
	Body telemetry.UmamiInfo
}

// analyticsHandler регистрирует эндпоинт, возвращающий конфигурацию Umami
// для подключения клиентского трекинг-скрипта на фронтенде.
func analyticsHandler(api huma.API, info *telemetry.UmamiInfo) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "get-analytics-config",
			Method:      http.MethodGet,
			Path:        "/analytics",
			Summary:     "Получить конфигурацию аналитики",
			Description: "Возвращает URL и WebsiteID для подключения Umami трекинг-скрипта.",
			Tags:        []string{"Analytics"},
		},
		func(ctx context.Context, req *struct{}) (*analyticsOutput, error) {
			return &analyticsOutput{Body: *info}, nil
		},
	)
}
