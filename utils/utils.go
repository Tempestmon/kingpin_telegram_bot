package utils

import (
	"context"
	"fmt"
	"kingpin_bot/bot_models"
	"log/slog"
	"os"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Фильтрация аудиофайлов по запросу
func FilterAudioFiles(files []string, query string) []string {
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
func GenerateAudioResults(files []string) []models.InlineQueryResult {
	results := make([]models.InlineQueryResult, 0, len(files))

	cfg, _ := bot_models.LoadConfig()

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
func LoadAudioFiles(filename string) ([]string, error) {
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
func LoggingMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
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
