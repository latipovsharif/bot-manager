package platforms

import (
	"fmt"
	"strconv"

	"github.com/latipovsharif/bot-manager/base"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Telegram /
type Telegram struct {
	Key     string
	Timeout int
	Debug   bool
}

var numericKeyboard = tg.NewInlineKeyboardMarkup(
	tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButtonData("1", "1"),
	),
)

// Run and listen for messages
func (t *Telegram) Run() error {
	bot, err := tg.NewBotAPI(t.Key)
	if err != nil {
		return errors.Wrap(err, "cannot connect to telegram")
	}

	bot.Debug = t.Debug
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tg.NewUpdate(0)
	u.Timeout = t.Timeout

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		chatID := getChatID(update)

		if update.Message != nil && update.Message.Text == "/start" {
			if _, err := bot.Send(getNextMsg(chatID)); err != nil {
				log.Errorf("cannot send message %v", err)
			}
			continue
		}

		page, err := base.GetCurrentPage(strconv.FormatInt(chatID, 10))
		if err != nil {
			log.Warnf("cannot get current page %d err: %v", chatID, err)
		}

		if page == base.PageOne {
			if messageIsValid(update) {
				if _, err := bot.Send(getNextMsg(chatID)); err != nil {
					log.Errorf("cannot send message %v", err)
				}
			}
		} else if page == base.PageTwo {
			if _, err := bot.Send(getPrevMsg(chatID, update.Message.Text)); err != nil {
				log.Errorf("cannot send message %v", err)
			}
		}

		log.Printf("page is %v, message is valid %v, chat %d", page, messageIsValid(update), chatID)
	}

	return nil
}

func getChatID(u tg.Update) int64 {
	if u.Message != nil {
		return u.Message.Chat.ID
	}

	if u.CallbackQuery != nil {
		return u.CallbackQuery.Message.Chat.ID
	}

	// FIXME what if both nil or Data or Text does not exists?
	return 0
}

func messageIsValid(u tg.Update) bool {
	return u.CallbackQuery != nil && u.CallbackQuery.Data == "1"
}

func getNextMsg(chatID int64) tg.MessageConfig {
	if err := base.SetCurrentPage(strconv.FormatInt(chatID, 10), base.PageTwo); err != nil {
		log.Warnf("cannot set page to second %v", err)
	}
	return tg.NewMessage(chatID, "Введите свое имя")
}

func getPrevMsg(chatID int64, name string) tg.MessageConfig {
	if err := base.SetCurrentPage(strconv.FormatInt(chatID, 10), base.PageOne); err != nil {
		log.Warnf("cannot set page to first %v", err)
	}
	msg := tg.NewMessage(chatID, fmt.Sprintf("Ваше имя %v 1 для возврата", name))
	msg.ReplyMarkup = numericKeyboard
	return msg
}
