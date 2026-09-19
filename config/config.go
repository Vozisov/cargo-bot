package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	BotToken      string
	ManagerChatID int64
}

func Load() (*Config, error) {
	if err := loadEnv(".env"); err != nil {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("BOT_TOKEN не задан")
	}

	chatIDStr := os.Getenv("MANAGER_CHAT_ID")
	if chatIDStr == "" {
		return nil, fmt.Errorf("MANAGER_CHAT_ID не задан")
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("MANAGER_CHAT_ID не число: %w", err)
	}

	return &Config{
		BotToken:      token,
		ManagerChatID: chatID,
	}, nil
}

// loadEnv — простой парсер .env без внешних библиотек
func loadEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}
	return scanner.Err()
}
