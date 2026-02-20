package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/viktorselivanov/homework/hw12_13_14_15_calendar/internal/app"
	"github.com/viktorselivanov/homework/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/viktorselivanov/homework/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/viktorselivanov/homework/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/viktorselivanov/homework/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return nil
	}

	cfg, err := NewConfigFromFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logg := logger.New(cfg.Logger.Level)

	var storage app.Storage
	switch cfg.Storage.Type {
	case "sql":
		sql := sqlstorage.New(cfg.DB.DSN)
		if err := sql.Connect(context.Background()); err != nil {
			return fmt.Errorf("failed to connect to db: %w", err)
		}
		storage = sql
	default:
		storage = memorystorage.New()
	}

	calendar := app.New(logg, storage)
	server := internalhttp.NewServer(logg, calendar, cfg.Server.Host, cfg.Server.Port)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer shutdownCancel()

		if err := server.Stop(shutdownCtx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	return nil
}
