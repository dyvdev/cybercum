package tgbot

import (
	tgbotapi "github.com/dyvdev/telegram-bot-api"
	"log"
	"math/rand"
)

func (bot *Bot) SendPhotoReaction(update tgbotapi.Update) {
	chat := bot.Chats[update.FromChat().ID]
	if chat.CanTalkSemen {
		emojis := []tgbotapi.Emoji{"💩", "❤️", "🔥", "🥰", "😁", "🤔", "🤯", "😱", "🥱"}
		emoji := emojis[rand.Intn(len(emojis))]
		react := tgbotapi.SetMessageReaction(update.FromChat().ID, update.Message.MessageID, emoji)
		_, err := bot.BotApi.Send(react)
		if err != nil {
			log.Println(err)
		}
	}
}
func (bot *Bot) SendRandomReaction(update tgbotapi.Update) {
	emojis := []tgbotapi.Emoji{"🤡", "🤔", "😁"}
	emoji := emojis[rand.Intn(len(emojis))]
	if rand.Intn(99) == 1 {
		react := tgbotapi.SetMessageReaction(update.FromChat().ID, update.Message.MessageID, emoji)
		_, err := bot.BotApi.Send(react)
		if err != nil {
			log.Println(err)
		}
	}
}
