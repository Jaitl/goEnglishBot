package command

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

func Parse(update tgbotapi.Update) (Command, error) {
	if update.Message != nil {
		return parseTextCommand(update.Message.From.ID, update.Message.Text), nil
	}

	if update.CallbackQuery != nil {
		return &KeyboardCallbackCommand{update.CallbackQuery.From.ID, update.CallbackQuery.Data}, nil
	}

	return nil, fmt.Errorf("unknown command: %+v", update)
}

func parseTextCommand(userId int, cmd string) Command {
	if strings.HasPrefix(cmd, "/") {
		parts := strings.Split(cmd, " ")
		cmd := parts[0]
		text := strings.Join(parts[1:], " ")

		switch cmd {
		case "/add":
			return &AddCommand{userId, text}
		}
	}

	return &TextCommand{userId, cmd}
}
