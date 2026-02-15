package webmuseum

import (
	"log/slog"

	"github.com/WhiCu/school-museum/internal/store"
	"github.com/WhiCu/school-museum/internal/web-museum/client"
	"github.com/WhiCu/school-museum/internal/web-museum/handler"
	"github.com/WhiCu/school-museum/internal/web-museum/service"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterHandlers(api huma.API, s *store.Store, log *slog.Logger) {
	stg := client.NewStorage(s, log.WithGroup("storage"))
	srv := service.NewService(stg, log.WithGroup("service"))
	h := handler.NewHandler(srv, log.WithGroup("handler"))

	h.Ping(api)
	h.GetAllNews(api)
	h.GetNewsByID(api)
	h.GetAllExhibitions(api)
	h.GetExhibitionByID(api)
}
