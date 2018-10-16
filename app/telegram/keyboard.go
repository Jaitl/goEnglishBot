package telegram

import "github.com/go-telegram-bot-api/telegram-bot-api"

type ButtonName string
type ButtonValue string

func CreateKeyboard(keys map[ButtonValue]ButtonName) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.InlineKeyboardMarkup{}

	var row []tgbotapi.InlineKeyboardButton

	for btValue, btName := range keys {
		btn := tgbotapi.NewInlineKeyboardButtonData(string(btName), string(btValue))
		row = append(row, btn)
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	return keyboard
}
