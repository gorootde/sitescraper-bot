package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func NewBot(token string) *Bot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Printf("Authorized on account %s", bot.Self.UserName)
	bot.Debug = true
	return &Bot{bot}
}

func (bot *Bot) Send(text string) {
	msg := tgbotapi.NewMessage(-301873798, text)
	msg.ParseMode = "HTML"
	bot.api.Send(msg)
}

//u := tgbotapi.NewUpdate(0)
//u.Timeout = 60

// updates, err := bot.GetUpdatesChan(u)

// for update := range updates {
// 	if update.Message == nil { // ignore any non-Message Updates
// 		continue
// 	}

// 	logrus.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
// 	msg.ReplyToMessageID = update.Message.MessageID

// 	bot.Send(msg)
// }
