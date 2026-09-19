package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vozisov/cargo-bot/order"
)

// handleMessage — точка входа для всех текстовых сообщений
func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	userID := msg.From.ID

	if msg.Text == BtnOrder {
		session := b.sessions.Get(userID)
		b.handleCommand(&tgbotapi.Message{
			From: msg.From,
			Chat: msg.Chat,
			Text: "/order",
		}, session)
		return
	}

	session := b.sessions.Get(userID)

	if msg.IsCommand() {
		b.handleCommand(msg, session)
		return
	}

	b.handleState(msg, session)
}

// handleCommand — обработка команд (/start, /order, /cancel)
func (b *Bot) handleCommand(msg *tgbotapi.Message, session *UserSession) {
	switch msg.Text {
	case "/start", BtnStart:
		b.sessions.Reset(msg.From.ID)
		reply := tgbotapi.NewMessage(msg.Chat.ID,
			"👋 Здравствуйте! Я бот для заказа грузоперевозок.\n\n"+
				"Нажмите «Оформить заказ», чтобы оставить заявку.")
		reply.ReplyMarkup = buildMainKeyboard()
		b.api.Send(reply)

	case "/order", BtnOrder:
		session.State = StateFromCity
		session.Order = order.Order{}
		b.sessions.Set(msg.From.ID, session)
		b.send(msg.Chat.ID, "📍 Откуда везём? Напишите город и адрес:")

	case "/cancel", BtnCancel:
		b.sessions.Reset(msg.From.ID)
		b.send(msg.Chat.ID, "❌ Заявка отменена. Начать заново — /order")

	default:
		b.send(msg.Chat.ID, "Неизвестная команда. Доступно: /order, /cancel")
	}
}

// handleState — обработка шагов диалога
func (b *Bot) handleState(msg *tgbotapi.Message, session *UserSession) {
	switch session.State {
	case StateIdle:
		b.send(msg.Chat.ID, "Чтобы оставить заявку, напишите /order")

	case StateFromCity:
		session.Order.FromCity = msg.Text
		session.State = StateToCity
		b.sessions.Set(msg.From.ID, session)
		b.send(msg.Chat.ID, "🎯 Куда везём? Напишите город и адрес:")

	case StateToCity:
		session.Order.ToCity = msg.Text
		session.State = StateCargo
		b.sessions.Set(msg.From.ID, session)
		b.send(msg.Chat.ID, "📦 Что везём? Опишите груз - переезд, мебель, стиральная машина и пр.")

	case StateCargo:
		session.Order.Cargo = msg.Text
		session.State = StateDateSelect
		b.sessions.Set(msg.From.ID, session)

		text := "📅 Когда нужна перевозка? Выберите из ближайших дней:"
		reply := tgbotapi.NewMessage(msg.Chat.ID, text)
		reply.ReplyMarkup = buildDateKeyboard()
		b.api.Send(reply)

	case StateDateManual:
		date, err := parseDate(msg.Text)
		if err != nil {
			b.send(msg.Chat.ID, "⚠️ Не понял дату. Введите в формате ДД.ММ.ГГГГ, например: 25.09.2026")
			return
		}
		session.Order.Date = date
		session.State = StateContact
		b.sessions.Set(msg.From.ID, session)
		b.send(msg.Chat.ID, fmt.Sprintf("✅ Дата: %s\n\n📞 Оставьте контакт для связи (телефон или @username).", date.Format("02.01.2006")))

	case StateContact:
		if !isValidContact(msg.Text) {
			b.send(msg.Chat.ID,
				"⚠️ Не похоже на контакт.\n\n"+
					"Введите номер телефона (например: +7 999 123-45-67)\n"+
					"или @username (например: @vasya).")
			return
		}
		session.Order.Contact = strings.TrimSpace(msg.Text)
		b.finishOrder(msg, session)
	}
}

// handleCallback — обработка нажатий на inline-кнопки
func (b *Bot) handleCallback(cb *tgbotapi.CallbackQuery) {
	userID := cb.From.ID
	session := b.sessions.Get(userID)

	// Всегда отвечаем на callback, иначе Telegram будет "крутить" кнопку
	callback := tgbotapi.NewCallback(cb.ID, "")
	if _, err := b.api.Request(callback); err != nil {
		log.Printf("ошибка ответа на callback: %v", err)
	}

	// Убираем кнопки, чтобы нельзя было нажать дважды
	edit := tgbotapi.NewEditMessageReplyMarkup(
		cb.Message.Chat.ID,
		cb.Message.MessageID,
		tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}},
	)
	if _, err := b.api.Send(edit); err != nil {
		log.Printf("ошибка удаления кнопок: %v", err)
	}

	// Обрабатываем callback только если ждём выбор даты
	if session.State != StateDateSelect {
		return
	}

	data := cb.Data // "date:2026-09-22" или "date:manual"
	value := strings.TrimPrefix(data, "date:")

	// Пользователь нажал "Другая дата" — переходим к ручному вводу
	if value == "manual" {
		session.State = StateDateManual
		b.sessions.Set(userID, session)
		b.send(cb.Message.Chat.ID, "📅 Введите дату в формате ДД.ММ.ГГГГ, например: 25.09.2026")
		return
	}

	// Пользователь выбрал конкретную дату
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		log.Printf("ошибка парсинга даты из callback: %v", err)
		b.send(cb.Message.Chat.ID, "⚠️ Что-то пошло не так с датой. Попробуйте ещё раз.")
		return
	}

	session.Order.Date = date
	session.State = StateContact
	b.sessions.Set(userID, session)

	b.send(cb.Message.Chat.ID, fmt.Sprintf("✅ Дата: %s\n\n📞 Оставьте контакт для связи (телефон или @username).", date.Format("02.01.2006")))
}

// finishOrder — финализация заказа и отправка менеджеру
func (b *Bot) finishOrder(msg *tgbotapi.Message, session *UserSession) {
	// TODO: сохранить в БД (когда добавим)
	session.Order.CreatedAt = time.Now()

	// Отправляем менеджеру
	managerMsg := tgbotapi.NewMessage(b.config.ManagerChatID, session.Order.String())
	managerMsg.ParseMode = "Markdown"

	if _, err := b.api.Send(managerMsg); err != nil {
		log.Printf("Ошибка отправки менеджеру: %v", err)
		b.send(msg.Chat.ID, "⚠️ Не удалось отправить заявку. Попробуйте позже.")
		return
	}

	b.sessions.Reset(msg.From.ID)
	b.send(msg.Chat.ID, "✅ Заявка отправлена! Менеджер свяжется с вами.")
}

// parseDate — парсит дату из нескольких форматов
func parseDate(input string) (time.Time, error) {
	formats := []string{"02.01.2006", "02.01.06", "2.1.2006", "2006-01-02"}
	for _, f := range formats {
		if t, err := time.Parse(f, input); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("неверный формат даты")
}

// send — вспомогательный метод для отправки простого сообщения
func (b *Bot) send(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Ошибка отправки: %v", err)
	}
}
