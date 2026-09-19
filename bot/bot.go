package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/vozisov/cargo-bot/config"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	config   *config.Config
	sessions *SessionStore
}

func New(cfg *config.Config) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	api.Debug = false
	log.Printf("Бот авторизован: @%s", api.Self.UserName)

	return &Bot{
		api:      api,
		config:   cfg,
		sessions: NewSessionStore(),
	}, nil
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			go b.handleMessage(update.Message)
		}
		if update.CallbackQuery != nil {
			go b.handleCallback(update.CallbackQuery)
		}
	}
}
