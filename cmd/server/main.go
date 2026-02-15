package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/WhiCu/school-museum/internal/config"
	"github.com/WhiCu/school-museum/internal/server"
)

func Run(cfg *config.Config, log *slog.Logger) error {
	server := server.NewApp(cfg, log)

	ctx := context.Background()

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
