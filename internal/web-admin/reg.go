package webadmin

import (
	"log/slog"

	"github.com/WhiCu/school-museum/internal/store"
	"github.com/WhiCu/school-museum/internal/web-admin/client"
	"github.com/WhiCu/school-museum/internal/web-admin/handler"
	"github.com/WhiCu/school-museum/internal/web-admin/service"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterHandlers(api huma.API, s *store.Store, log *slog.Logger) {
	stg := client.NewStorage(s, log.WithGroup("storage"))
	srv := service.NewService(stg, log.WithGroup("service"))
	h := handler.NewHandler(srv, log.WithGroup("handler"))

	h.Ping(api)

	// News
	h.CreateNews(api)
	h.DeleteNews(api)

	// Exhibitions
	h.CreateExhibition(api)
	h.UpdateExhibition(api)
	h.DeleteExhibition(api)

	// Exhibits
	h.CreateExhibit(api)
	h.UpdateExhibit(api)
	h.DeleteExhibit(api)
}
