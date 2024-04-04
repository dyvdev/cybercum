package tgbot

import (
	tgbotapi "github.com/dyvdev/telegram-bot-api"
	"math/rand"
	"regexp"
	"strings"
)

func (bot *Bot) SendContextPhrase(update tgbotapi.Update) bool {
	if update.Message == nil {
		return false
	}
	msg := strings.ToLower(regexp.MustCompile(`[.,;!?]`).ReplaceAllString(update.Message.Text, ""))
	words := strings.Split(msg, " ")
	if len(words) == 0 || len(words) > 25 {
		return false
	}
	//chat := bot.Chats[update.FromChat().ID]
	return true
}
func (bot *Bot) SendAnswer(update tgbotapi.Update) bool {
	if strings.Index(update.Message.Text[1:], "?")+1 != 0 {
		str := "да"
		if rand.Intn(10) == 1 {
			str = "ебёт тебя?"
		} else if rand.Intn(2) == 0 {
			str = "нет"
		}
		msg := tgbotapi.NewMessage(update.FromChat().ID, str)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.BotApi.Send(msg)
		return true
	}
	return false
}
