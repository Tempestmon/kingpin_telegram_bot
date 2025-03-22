package handlers

import (
	"context"
	"kingpin_bot/bot_models"
	"kingpin_bot/utils"
	"log/slog"
	"os"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Обработчик запросов с зависимостями
type Handler struct {
	cfg *bot_models.Config
}

func NewHandler(cfg *bot_models.Config) *Handler {
	return &Handler{cfg: cfg}
}

// Обработчик входящих сообщений
func (h *Handler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.InlineQuery == nil {
		return
	}

	offset, _ := strconv.Atoi(update.InlineQuery.Offset)

	// Формирование ответа
	audioFiles, err := utils.LoadAudioFiles(h.cfg.AudiofilesPath)
	if err != nil {
		slog.Error("Failed to load audio files", "error", err)
		os.Exit(1)
	}

	query := update.InlineQuery.Query

	filteredFiles := utils.FilterAudioFiles(audioFiles, query)

	allResults := utils.GenerateAudioResults(filteredFiles)

	limit := 50
	start := offset
	end := min(start+limit, len(allResults))

	results := allResults[start:end]

	nextOffset := ""
	if end < len(allResults) {
		nextOffset = strconv.Itoa(end)
	}
	// Отправка ответа
	_, err = b.AnswerInlineQuery(ctx, &bot.AnswerInlineQueryParams{
		InlineQueryID: update.InlineQuery.ID,
		Results:       results,
		NextOffset:    nextOffset,
	})
	if err != nil {
		slog.Error("Failed to answer query",
			"query_id", update.InlineQuery.ID,
			"error", err,
		)
		return
	}

	slog.Info("Answered query",
		"query_id", update.InlineQuery.ID,
		"user_id", update.InlineQuery.From.ID,
		"offset", update.InlineQuery.From.ID,
		"results_count", len(results),
		"next_offset", nextOffset,
	)
}
