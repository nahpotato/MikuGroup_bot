package main

import (
	"log"
	"net/http"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type bot struct {
	token   string
	botAPI  *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
}

func botNew(token string) *bot {
	bot := new(bot)

	bot.token = token

	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.botAPI = botAPI

	return bot
}

func (bot *bot) run(port int) {
	_, err := bot.botAPI.SetWebhook(tgbotapi.NewWebhook("https://ponyrevolution-bot.herokuapp.com/" + bot.token))
	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.botAPI.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	} else if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	bot.updates = bot.botAPI.ListenForWebhook("/" + bot.token)
	go http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), nil)

	bot.listenToCommands()
}

func (bot *bot) listenToCommands() {
	for update := range bot.updates {
		if update.Message != nil {
			continue
		}

		if update.Message.Chat.Type == "supergroup" {
			if !update.Message.IsCommand() {
				continue
			}

			switch update.Message.Command() {
			case "ban":
				ban(bot.botAPI, update)
				break
			}
		}
	}
}
