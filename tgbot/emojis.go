package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/rand"
)

func (bot *Bot) SendPhotoReaction(update tgbotapi.Update) {
	chat := bot.Chats[update.FromChat().ID]
	if chat.CanTalkSemen {
		emojis := []tgbotapi.Emoji{"🍌", "🌭", "💩", "❤️", "🔥", "🥰", "😁", "🤔", "🤯", "😱", "🥱"}
		emoji := emojis[rand.Intn(len(emojis)-1)]
		react := tgbotapi.SetMessageReaction(update.FromChat().ID, update.Message.MessageID, emoji)
		_, err := bot.BotApi.Send(react)
		if err != nil {
			log.Println(err)
		}
	}
}
func (bot *Bot) SendRandomReaction(update tgbotapi.Update) {
	emojis := []tgbotapi.Emoji{"🤡", "🤔", "😁"}
	emoji := emojis[rand.Intn(len(emojis)-1)]
	if rand.Intn(99) == 1 {
		react := tgbotapi.SetMessageReaction(update.FromChat().ID, update.Message.MessageID, emoji)
		_, err := bot.BotApi.Send(react)
		if err != nil {
			log.Println(err)
		}
	}
}
