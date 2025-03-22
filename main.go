package main

import (
	"context"
	"kingpin_bot/bot_models"
	"kingpin_bot/handlers"
	"kingpin_bot/utils"
	"log/slog"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
)

func main() {
	// Настройка логгера
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// Загрузка конфигурации
	cfg, err := bot_models.LoadConfig()
	if err != nil {
		slog.Error("Configuration error", "error", err)
		os.Exit(1)
	}

	// Создание обработчика
	handler := handlers.NewHandler(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Настройки бота
	opts := []bot.Option{
		bot.WithDefaultHandler(handler.Handle),
		bot.WithMiddlewares(utils.LoggingMiddleware),
	}

	// Инициализация бота
	b, err := bot.New(cfg.Token, opts...)
	if err != nil {
		slog.Error("Bot initialization failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting bot")
	b.Start(ctx)
}
