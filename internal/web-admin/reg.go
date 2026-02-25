package webadmin

import (
	"log/slog"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/WhiCu/school-museum/db/storage"
	"github.com/WhiCu/school-museum/internal/web-admin/client"
	"github.com/WhiCu/school-museum/internal/web-admin/handler"
	"github.com/WhiCu/school-museum/internal/web-admin/service"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterHandlers(
	api huma.API,
	news storage.Storage[model.News],
	exhibitions *storage.ExhibitionStorage,
	exhibits storage.Storage[model.Exhibit],
	visits *storage.VisitStorage,
	log *slog.Logger) {
	stg := client.NewStorage(news, exhibitions, exhibits, visits, log.WithGroup("storage"))
	srv := service.NewService(stg, log.WithGroup("service"))
	h := handler.NewHandler(srv, log.WithGroup("handler"))

	h.Ping(api)

	// News
	h.CreateNews(api)
	h.UpdateNews(api)
	h.DeleteNews(api)

	// Exhibitions
	h.CreateExhibition(api)
	h.UpdateExhibition(api)
	h.DeleteExhibition(api)
	h.SetExhibitionPreview(api)

	// Exhibits
	h.CreateExhibit(api)
	h.UpdateExhibit(api)
	h.DeleteExhibit(api)

	// Stats
	h.GetStats(api)
}
