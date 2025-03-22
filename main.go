package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Конфигурация приложения
type Config struct {
	Token    string
	AudioURL string
}

// Обработчик запросов с зависимостями
type Handler struct {
	cfg *Config
}

// Загрузка конфигурации из переменных окружения
func loadConfig() (*Config, error) {
	token := os.Getenv("API_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("API_TOKEN is required")
	}

	audioURL := os.Getenv("AUDIO_URL")
	if audioURL == "" {
		return nil, fmt.Errorf("AUDIO_URL is required")
	}

	return &Config{
		Token:    token,
		AudioURL: audioURL,
	}, nil
}

// Фильтрация аудиофайлов по запросу
func filterAudioFiles(files []string, query string) []string {
	filtered := make([]string, 0)

	for _, file := range files {
		// Убираем расширение .ogg и заменяем подчеркивания на пробелы
		title := strings.TrimSuffix(file, ".ogg")
		title = strings.ReplaceAll(title, "_", " ")

		// Ищем совпадение (без учета регистра)
		if strings.Contains(strings.ToLower(title), strings.ToLower(query)) {
			filtered = append(filtered, file)
		}
	}

	return filtered
}

// Генерация аудио-объектов на основе списка файлов
func generateAudioResults(files []string) []models.InlineQueryResult {
	results := make([]models.InlineQueryResult, 0, len(files))

	cfg, _ := loadConfig()

	for i, file := range files {
		// Убираем расширение .ogg для формирования Title
		title := strings.TrimSuffix(file, ".ogg")
		// Заменяем подчеркивания на пробелы и делаем первую букву заглавной
		title = strings.ReplaceAll(title, "_", " ")
		title = strings.Title(title)

		results = append(results, &models.InlineQueryResultVoice{
			ID:       fmt.Sprintf("%d", i+1), // ID начинается с 1
			VoiceURL: fmt.Sprintf("%s/%s", cfg.AudioURL, file),
			Title:    title,
		})
	}

	return results
}

// Загрузка списка файлов из текстового файла
func loadAudioFiles(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	files := strings.Split(string(data), "\n")
	// Убираем пустые строки
	cleanFiles := make([]string, 0, len(files))
	for _, file := range files {
		file = strings.TrimSpace(file)
		if file != "" {
			cleanFiles = append(cleanFiles, file)
		}
	}

	return cleanFiles, nil
}

// Middleware для логгирования входящих запросов
func loggingMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.InlineQuery != nil {
			slog.Info("Incoming inline query",
				"user_id", update.InlineQuery.From.ID,
				"query", update.InlineQuery.Query,
			)
		}
		next(ctx, b, update)
	}
}

func main() {
	// Настройка логгера
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// Загрузка конфигурации
	cfg, err := loadConfig()
	if err != nil {
		slog.Error("Configuration error", "error", err)
		os.Exit(1)
	}

	// Создание обработчика
	handler := &Handler{cfg: cfg}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Настройки бота
	opts := []bot.Option{
		bot.WithDefaultHandler(handler.Handle),
		bot.WithMiddlewares(loggingMiddleware),
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

// Обработчик входящих сообщений
func (h *Handler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.InlineQuery == nil {
		return
	}

	offset, _ := strconv.Atoi(update.InlineQuery.Offset)

	// Формирование ответа
	audioFiles, err := loadAudioFiles("./audiofiles.txt")
	if err != nil {
		slog.Error("Failed to load audio files", "error", err)
		os.Exit(1)
	}

	query := update.InlineQuery.Query

	filteredFiles := filterAudioFiles(audioFiles, query)

	allResults := generateAudioResults(filteredFiles)

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
