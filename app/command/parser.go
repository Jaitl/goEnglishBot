package command

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
		case "/list", "/l":
			incNumber, err := parseOptionalIncNumber(parts)
			if err != nil {
				return nil, err
			}
			return &ListPhrasesCommand{UserId: userId, IncNumber: incNumber}, nil
		case "/remove", "/r":
			incNumber, err := parseIncNumber(parts)
			if err != nil {
				return nil, err
			}
			return &RemovePhraseCommand{userId, *incNumber}, nil
		case "/cat":
			if len(parts) < 2 {
				return nil, errors.New("name is empty")
			}
			name := strings.Join(parts[1:], " ")
			return &AddCategoryCommand{userId, name}, nil
		case "/cats", "/cl":
			return &ListCategoriesCommand{userId}, nil
		case "/set":
			incNumber, err := parseIncNumber(parts)
			if err != nil {
				return nil, err
			}
			return &SetCategoriesCommand{UserId: userId, IncNumber: *incNumber}, nil
		case "/catRm", "/catrm", "/crm":
			incNumber, err := parseIncNumber(parts)
			if err != nil {
				return nil, err
			}
			return &RemoveCategoryCommand{UserId: userId, IncNumber: *incNumber}, nil
		case "/speech", "/sp":
			from, to, err := parseIntRange(parts)
			if err != nil {
				return nil, err
			}
			return &SpeechCommand{UserId: userId, From: from, To: to}, nil
		case "/puzzleAudio", "/pa":
			from, to, err := parseIntRange(parts)
			if err != nil {
				return nil, err
			}
			return &PuzzleAudioCommand{UserId: userId, From: from, To: to}, nil
		case "/puzzleTrans", "/pt":
			from, to, err := parseIntRange(parts)
			if err != nil {
				return nil, err
			}
			return &PuzzleTransCommand{UserId: userId, From: from, To: to}, nil
		case "/writeAudio", "/wa":
			from, to, err := parseIntRange(parts)
			if err != nil {
				return nil, err
			}
			return &WriteAudioCommand{UserId: userId, From: from, To: to}, nil
		case "/writeTrans", "/wt":
			from, to, err := parseIntRange(parts)
			if err != nil {
				return nil, err
			}
			return &WriteTransCommand{UserId: userId, From: from, To: to}, nil
		case "/skip", "/sk":
			return &SkipCommand{UserId: userId}, nil
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

func parseIncNumber(parts []string) (*int, error) {
	if len(parts) != 2 {
		return nil, errors.New("not enough arguments for the command")
	}
	incNumber, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	return &incNumber, nil
}

func parseOptionalIncNumber(parts []string) (*int, error) {
	if len(parts) != 2 {
		return nil, nil
	}

	if len(parts) > 2 {
		return nil, errors.New("too many arguments for the command")
	}

	incNumber, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	return &incNumber, nil
}

func parseIntRange(parts []string) (*int, *int, error) {
	var from, to *int = nil, nil

	if len(parts) > 3 {
		return nil, nil, errors.New("too many arguments for the command")
	}

	if len(parts) >= 2 {
		numb, err := strconv.Atoi(parts[1])

		if err != nil {
			return nil, nil, err
		}

		from = &numb
	}

	if len(parts) == 3 {
		numb, err := strconv.Atoi(parts[2])

		if err != nil {
			return nil, nil, err
		}

		to = &numb

		if *from >= *to {
			return nil, nil, errors.New("'from' cannot be more than 'to'")
		}
	}

	return from, to, nil
}
