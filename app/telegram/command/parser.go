package command

import (
	"errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
)

func Parse(update tgbotapi.Update) (Command, error) {
	if update.Message != nil {
		if update.Message.Voice != nil {
			return parseVoiceCommand(update.Message.From.ID, update.Message.Voice)
		} else {
			return parseTextCommand(update.Message.From.ID, update.Message.Text)
		}
	}

	if update.CallbackQuery != nil {
		return &KeyboardCallbackCommand{update.CallbackQuery.From.ID, update.CallbackQuery.Data}, nil
	}

	return nil, fmt.Errorf("unknown command: %+v", update)
}

func parseTextCommand(userId int, cmd string) (Command, error) {
	if strings.HasPrefix(cmd, "/") {
		parts := strings.Split(cmd, " ")
		cmd := parts[0]

		switch cmd {
		case "/add":
			text := strings.Join(parts[1:], " ")
			return &AddCommand{userId, text}, nil
		case "/list":
			return &ListCommand{userId}, nil
		case "/audio":
			if len(parts) != 2 {
				return nil, errors.New("not enough arguments for the audio command")
			}
			incNumber, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return nil, err
			}
			return &AudioCommand{userId, int(incNumber)}, nil
		case "/voice":
			if len(parts) != 2 {
				return nil, errors.New("not enough arguments for the voice command")
			}
			incNumber, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return nil, err
			}
			return &VoiceCommand{userId, int(incNumber)}, nil
		}
	}

	return &TextCommand{userId, cmd}, nil
}

func parseVoiceCommand(userId int, voice *tgbotapi.Voice) (Command, error) {
	return &ReceivedVoiceCommand{userId}, nil
}
