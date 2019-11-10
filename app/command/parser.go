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
			incNumber, err := strconv.ParseInt(update.Message.Text, 10, 32)
			if err == nil {
				return parseNumberCommand(update.Message.From.ID, int(incNumber))
			}

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
		case "/me":
			return &MeCommand{userId}, nil
		case "/add":
			text := strings.Join(parts[1:], " ")
			return &AddCommand{userId, text}, nil
		case "/list", "/l":
			return &ListCommand{userId}, nil
		case "/remove", "/r":
			incNumber, err := parseIncNumber(parts)
			if err != nil {
				return nil, err
			}
			return &RemoveCommand{userId, int(*incNumber)}, nil
		case "/voice", "/v":
			incNumber, err := parseIncNumber(parts)
			if err != nil {
				return nil, err
			}
			return &VoiceCommand{userId, int(*incNumber)}, nil
		case "/puzzleAudio", "/pa":
			incNumber, err := parseIncNumber(parts)
			if err != nil {
				return nil, err
			}
			return &PuzzleAudioCommand{userId, int(*incNumber)}, nil
		case "/puzzleTrans", "/pt":
			incNumber, err := parseIncNumber(parts)
			if err != nil {
				return nil, err
			}
			return &PuzzleTransCommand{userId, int(*incNumber)}, nil
		case "/writeAudio", "/wa":
			incNumber, err := parseIncNumber(parts)
			if err != nil {
				return nil, err
			}
			return &WriteAudioCommand{userId, int(*incNumber)}, nil
		case "/writeTrans", "/wt":
			incNumber, err := parseIncNumber(parts)
			if err != nil {
				return nil, err
			}
			return &WriteTransCommand{userId, int(*incNumber)}, nil
		default:
			return nil, fmt.Errorf("unknown command: %+v", cmd)
		}
	}

	return &TextCommand{userId, cmd}, nil
}

func parseNumberCommand(userId, incNumber int) (Command, error) {
	return &NumberCommand{userId, incNumber}, nil
}

func parseVoiceCommand(userId int, voice *tgbotapi.Voice) (Command, error) {
	return &ReceivedVoiceCommand{UserId: userId, FileID: voice.FileID}, nil
}

func parseIncNumber(parts []string) (*int64, error) {
	if len(parts) != 2 {
		return nil, errors.New("not enough arguments for the command")
	}
	incNumber, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, err
	}

	return &incNumber, nil
}
