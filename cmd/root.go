package cmd

import (
	"log/slog"
	"os"

	"github.com/WhiCu/school-museum/cmd/server"
	"github.com/WhiCu/school-museum/internal/config"
	"github.com/WhiCu/school-museum/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	cfg *config.Config
	log *slog.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "school-museum",
	Short: "",
	Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		ft, err := cmd.Flags().GetString("filetype")
		if err != nil {
			return err
		}
		filetype, err := config.FileType(ft)
		if err != nil {
			return err
		}
		cfg, err = config.Load[config.Config](filetype)
		if err != nil {
			return err
		}

		log = logger.GetLogger(&cfg.Logger)
		log.Info("logger created",
			slog.String("level", cfg.Logger.Level),
			slog.String("path", cfg.Logger.Path),
			slog.Int("size", cfg.Logger.Size),
			slog.Bool("compress", cfg.Logger.Compress),
		)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return server.Run(cfg, log)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.Flags().StringP("filetype", "t", "yaml", "")
}
