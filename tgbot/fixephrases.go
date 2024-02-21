package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
)

func (bot *Bot) SendFixedPhrase(message *tgbotapi.Message) {
	chat := bot.Chats[message.Chat.ID]
	if len(chat.FixedPhrases) != 0 {
		threadId := 0
		if message.Chat.IsForum && message.MessageThreadID != 0 {
			threadId = message.MessageThreadID
		}
		bot.SendMessage(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           message.Chat.ID,
				MessageThreadID:  threadId,
				ReplyToMessageID: 0,
			},
			Text:                  chat.FixedPhrases[rand.Intn(len(chat.FixedPhrases)-1)],
			DisableWebPagePreview: false,
		})
	}
}
func (bot *Bot) AddFixedPhrase(chat *Chat, str string) int {
	if str != "" {
		chat.FixedPhrases = append(chat.FixedPhrases, str)
		bot.SaveDump()
		return len(chat.FixedPhrases) - 1
	}
	return -1
}

func (bot *Bot) RemoveFixedPhrase(chat *Chat, id int) int {
	if id > -1 && id < len(chat.FixedPhrases) {
		chat.FixedPhrases = append(chat.FixedPhrases[:id], chat.FixedPhrases[id+1:]...)
		bot.SaveDump()
		return len(chat.FixedPhrases)
	}
	return -1
}
