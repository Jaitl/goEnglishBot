package telegram

import "github.com/go-telegram-bot-api/telegram-bot-api"

type ButtonName string
type ButtonValue string

const (
	RowSize int = 3
)

func CreateKeyboard(keys map[ButtonValue]ButtonName) tgbotapi.InlineKeyboardMarkup {
	var keyboard [][]tgbotapi.InlineKeyboardButton
	row := make([]tgbotapi.InlineKeyboardButton, 0, RowSize)

	curRowCount := 0
	for btValue, btName := range keys {
		curRowCount += 1
		btn := tgbotapi.NewInlineKeyboardButtonData(string(btName), string(btValue))
		row = append(row, btn)
		if curRowCount >= RowSize {
			curRowCount = 0
			keyboard = append(keyboard, row)
			row = make([]tgbotapi.InlineKeyboardButton, 0, RowSize)
		}
	}

	if curRowCount > 0 {
		keyboard = append(keyboard, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}
