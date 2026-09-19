# Cargo Bot

Telegram-бот для заказа грузоперевозок.

## Что умеет

- 📦 Пошаговый диалог оформления заказа
- 📅 Выбор даты через inline-кнопки (7 дней вперёд) или вручную
- ✅ Валидация контакта (телефон или @username)
- 📨 Отправка заявки менеджеру в Telegram-чат

## Стек

- **Go 1.22+**
- [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) — работа с Telegram Bot API
- **FSM** — управление состояниями диалога

## Структура проекта

cargo-bot/
├── main.go # точка входа
├── config/ # загрузка конфига из .env
├── bot/ # логика бота
│ ├── bot.go # запуск, polling
│ ├── handlers.go # обработка сообщений и callback
│ ├── fsm.go # состояния диалога (FSM)
│ ├── keyboard.go # клавиатуры (reply и inline)
│ └── validators.go # валидация полей
└── order/ # модель заказа
└── order.go


## Запуск

### 1. Клонируй репозиторий

```bash
git clone git@github.com:vozisov/cargo-bot.git
cd cargo-bot

### 2. Создай `.env`

Создай файл `.env` в корне проекта:

```env
BOT_TOKEN=...
MANAGER_CHAT_ID=...

Где взять:

    BOT_TOKEN — у @BotFather → /newbot

    MANAGER_CHAT_ID — ID группы менеджера. Узнать можно через @userinfobot

### 3. Установи зависимости и запусти

go mod download
go run main.go

Команды бота
Команда	Описание
/start	Приветствие, показать главное меню
/order	Оформить заказ
/cancel	Отменить заполнение


Как это работает

    Клиент пишет боту /order

    Бот ведёт диалог по шагам (FSM): откуда → куда → груз → дата → контакт

    Дата выбирается через inline-кнопки (7 дней вперёд) или вручную

    Контакт валидируется (телефон или @username)

    Заявка отправляется в чат менеджера

TODO

    □

    Сохранение заявок в PostgreSQL
    □

    Команды для менеджера (/list, /today)
    □

    Статусы заявок (новая, в работе, выполнена)
    □

    Docker + деплой на VPS
    □

    Кнопка «Отмена» на каждом шаге диалога

Лицензия

MIT