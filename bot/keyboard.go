package bot

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	BtnOrder  = "Оформить заказ"
	BtnStart  = "Старт"
	BtnCancel = "Отменить"
)

// buildDateKeyboard — создаёт клавиатуру с датами на неделю вперёд
func buildDateKeyboard() tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	today := time.Now()
	var currentRow []tgbotapi.InlineKeyboardButton

	for i := 0; i < 7; i++ {
		date := today.AddDate(0, 0, i)
		label := formatDateLabel(date, i)

		btn := tgbotapi.NewInlineKeyboardButtonData(
			label,
			"date:"+date.Format("2006-01-02"), // callback data
		)
		currentRow = append(currentRow, btn)

		// По 3 кнопки в ряд
		if len(currentRow) == 3 || i == 6 {
			rows = append(rows, currentRow)
			currentRow = nil
		}
	}

	// Кнопка "Другая дата"
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📅 Другая дата", "date:manual"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// formatDateLabel — красивая подпись кнопки
func formatDateLabel(date time.Time, offset int) string {
	switch offset {
	case 0:
		return "Сегодня"
	case 1:
		return "Завтра"
	default:
		// "Пн 22.09"
		weekdays := map[time.Weekday]string{
			time.Monday:    "Пн",
			time.Tuesday:   "Вт",
			time.Wednesday: "Ср",
			time.Thursday:  "Чт",
			time.Friday:    "Пт",
			time.Saturday:  "Сб",
			time.Sunday:    "Вс",
		}
		return fmt.Sprintf("%s %s", weekdays[date.Weekday()], date.Format("02.01"))
	}
}

// buildMainKeyboard — постоянная клавиатура с кнопкой заказа
func buildMainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnOrder),
		),
	)
	keyboard.ResizeKeyboard = true
	return keyboard
}
