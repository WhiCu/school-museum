package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/WhiCu/school-museum/internal/config"
	"github.com/WhiCu/school-museum/internal/server"
)

func Run(cfg *config.Config, log *slog.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server := server.NewApp(ctx, cfg, log)

	if err := server.Run(ctx); err != nil {
		log.Error("Server failed to start", slog.String("ERR", err.Error()))
		fmt.Printf(`
		==================================
		=                                =
		=    Server failed to start     =
		=                                =
		==================================
	%v
			`, err)
		return err
	}
	log.Info("server successfully stopped")
	fmt.Println(`
	==================================
	=                                =
	=    Server successfully stop    =
	=                                =
	==================================
	`)
	return nil
}
