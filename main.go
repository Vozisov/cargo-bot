package main

import (
	"log"

	"github.com/vozisov/cargo-bot/bot"
	"github.com/vozisov/cargo-bot/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка конфига: %v", err)
	}

	b, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("Ошибка запуска бота: %v", err)
	}

	log.Println("Бот запущен. Нажмите Ctrl+C для остановки.")
	b.Run()
}
