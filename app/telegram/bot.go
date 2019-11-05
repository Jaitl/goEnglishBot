package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/command"
	"log"
)

type Telegram struct {
	userId       int
	bot          *tgbotapi.BotAPI
	updateChanel tgbotapi.UpdatesChannel
}

func New(token string, userId int) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		return nil, err
	}

	return &Telegram{userId: userId, bot: bot, updateChanel: updates}, nil
}

func (t *Telegram) Start(executor *action.Executor) {
	for update := range t.updateChanel {
		log.Printf("[DEBUG] new telegram message: %v", update)
		go t.handleMessage(update, executor)
	}
}

func (t *Telegram) handleMessage(update tgbotapi.Update, executor *action.Executor) {
	cmd, err := command.Parse(update)

	if err != nil {
		msg := fmt.Sprintf("[ERROR] error during parse: %v", err)
		log.Println(msg)
		if update.Message != nil {
			err := t.Send(int(update.Message.Chat.ID), msg)
			if err != nil {
				log.Printf("[ERROR] error during send parse error: %v", err)
			}
		}
		return
	}

	if !t.accessFilter(cmd) {
		err := t.Send(cmd.GetUserId(), "Доступ запрешен")
		if err != nil {
			log.Printf("[ERROR] error check access: %v", err)
		}
		log.Printf("[ERROR] Unknown user: %d", cmd.GetUserId())
		return
	}

	err = executor.Execute(cmd)

	if err != nil {
		msg := fmt.Sprintf("[ERROR] error during execution cmd: %s, err: %v", cmd.GetType(), err)
		log.Println(msg)
		err := t.Send(cmd.GetUserId(), msg)
		if err != nil {
			log.Printf("[ERROR] error during send execution error: %v", err)
		}
		return
	}
}

func (t *Telegram) accessFilter(cmd command.Command) bool {
	switch cmd.(type) {
	case *command.MeCommand:
		return true
	default:
		if t.userId != cmd.GetUserId() {
			return false
		}

		return true
	}
}

func (t *Telegram) SendWithKeyboard(chatId int, message string, keyboard map[ButtonValue]ButtonName) error {
	msg := tgbotapi.NewMessage(int64(chatId), message)

	keys := CreateKeyboard(keyboard)
	msg.ReplyMarkup = keys

	_, err := t.bot.Send(msg)

	return err
}

func (t *Telegram) Send(chatId int, message string) error {
	msg := tgbotapi.NewMessage(int64(chatId), message)

	_, err := t.bot.Send(msg)

	return err
}

func (t *Telegram) SendMarkdown(chatId int, message string) error {
	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           int64(chatId),
			ReplyToMessageID: 0,
		},
		Text:                  message,
		ParseMode:             tgbotapi.ModeMarkdown,
		DisableWebPagePreview: false,
	}

	_, err := t.bot.Send(msg)

	return err
}

func (t *Telegram) SendAudio(chatId int, pathToVoice string) (string, error) {
	msg := tgbotapi.NewAudioUpload(int64(chatId), pathToVoice)

	message, err := t.bot.Send(msg)

	if err != nil {
		return "", err
	}

	if message.Audio != nil {
		return message.Audio.FileID, nil
	}

	return "", nil
}

func (t *Telegram) SendAudioId(chatId int, audioId string) error {
	msg := tgbotapi.NewAudioShare(int64(chatId), audioId)

	_, err := t.bot.Send(msg)

	return err
}

func (t *Telegram) GetFilePath(fileID string) (string, error) {
	file, err := t.bot.GetFileDirectURL(fileID)

	if err != nil {
		return "", err
	}

	return file, nil
}
