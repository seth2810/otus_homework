package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/seth2810/otus_homework/hw12_13_14_15_calendar/internal/calendar"
	"github.com/seth2810/otus_homework/hw12_13_14_15_calendar/internal/commands"
	"github.com/seth2810/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/seth2810/otus_homework/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/seth2810/otus_homework/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/seth2810/otus_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/cobra"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

var configPath string

var rootCmd = &cobra.Command{
	Use: "calendar",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &calendar.Config{}

		if err := config.ReadConfig(cfg, configPath); err != nil {
			return fmt.Errorf("failed to read config: %w", err)
		}

		return startApp(cmd.Context(), cfg)
	},
}

func init() {
	rootCmd.Flags().StringVar(&configPath, "config", "/etc/calendar/config.yaml", "Path to configuration file")
	rootCmd.AddCommand(commands.NewVersionCmd(release, buildDate, gitHash))
}

func main() {
	ctx, cancelFn := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP,
	)

	defer cancelFn()

	cobra.CheckErr(rootCmd.ExecuteContext(ctx))
}

func initStorage(ctx context.Context, cfg calendar.StorageConfig) (calendar.Storage, error) {
	switch cfg.Type {
	case "memory":
		return memorystorage.New(), nil
	case "sql":
		return sqlstorage.Init(ctx, cfg.Database)
	default:
		return nil, fmt.Errorf("unrecognized type: %q", cfg.Type)
	}
}

func startApp(ctx context.Context, cfg *calendar.Config) error {
	logger, err := logger.New(cfg.Logger.Level, cfg.Logger.File)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	storage, err := initStorage(ctx, cfg.Storage)
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	app := calendar.New(logger, storage)

	return app.Serve(ctx, cfg.Server)
}
