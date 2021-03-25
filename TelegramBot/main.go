package main

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v5"
	"net/http"
	"time"
)

const (
	BotToken   = "1615859986:AAE_4fUhJebvm2HA0w-KaeNcUe1JgsSdsQo"
	WebhookURL = "https://466c1020453e.ngrok.io"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	updates := bot.ListenForWebhook("/")

	go http.ListenAndServe(":8088", nil)
	fmt.Println("Start listen :8088")

	for update := range updates {
		if update.Message.Text == "hello" {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "hi"))
		} else if update.Message.Text == "time" {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, time.Now().String()))
		} else {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда"))
		}
	}
}
