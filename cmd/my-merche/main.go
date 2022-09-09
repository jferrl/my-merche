package main

import (
	"log"
	"os"

	"github.com/yanzay/tbot/v2"
)

var ttoken = os.Getenv("TELEGRAM_TOKEN")

func main() {
	bot := tbot.New(ttoken)
	c := bot.Client()
	bot.HandleMessage("/info", func(m *tbot.Message) {
		c.SendMessage(m.Chat.ID, "hello!")
	})
	err := bot.Start()
	if err != nil {
		log.Fatal(err)
	}
}
