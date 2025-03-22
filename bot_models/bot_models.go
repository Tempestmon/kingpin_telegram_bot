package bot_models

import (
	"fmt"
	"os"
)

// Конфигурация приложения
type Config struct {
	Token          string
	AudioURL       string
	AudiofilesPath string
}

// Загрузка конфигурации из переменных окружения
func LoadConfig() (*Config, error) {
	token := os.Getenv("API_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("API_TOKEN is required")
	}

	audioURL := os.Getenv("AUDIO_URL")
	if audioURL == "" {
		return nil, fmt.Errorf("AUDIO_URL is required")
	}

	audiofiles := os.Getenv("AUDIOFILES_PATH")
	if audiofiles == "" {
		return nil, fmt.Errorf("AUDIOFILES_PATH is required")
	}

	return &Config{
		Token:          token,
		AudioURL:       audioURL,
		AudiofilesPath: audiofiles,
	}, nil
}
