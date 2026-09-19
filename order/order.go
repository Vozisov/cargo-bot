package order

import (
	"fmt"
	"time"
)

type Order struct {
	FromCity  string
	ToCity    string
	Cargo     string
	Date      time.Time // ← было string
	Contact   string
	CreatedAt time.Time
}

func (o Order) String() string {
	dateStr := "не указана"
	if !o.Date.IsZero() {
		dateStr = o.Date.Format("02.01.2006")
	}

	return fmt.Sprintf(
		"🚚 *Новая заявка на грузоперевозку*\n\n"+
			"📍 Откуда: %s\n"+
			"🎯 Куда: %s\n"+
			"📦 Груз: %s\n"+
			"📅 Дата: %s\n"+
			"📞 Контакт: %s\n\n"+
			"🕐 Создано: %s",
		o.FromCity,
		o.ToCity,
		o.Cargo,
		dateStr,
		o.Contact,
		o.CreatedAt.Format("02.01.2006 15:04"),
	)
}
