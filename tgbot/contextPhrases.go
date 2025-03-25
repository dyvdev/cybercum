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
	return true
}
func (bot *Bot) SendAnswer(update tgbotapi.Update) bool {
	if strings.Index(update.Message.Text[1:], "?")+1 != 0 {
		str := "да"
		if rand.Intn(10) == 1 {
			str = "ебёт тебя?"
		} else if rand.Intn(5) == 1 {
			str = "да ты заебал уже всех со своими вопросами"
		} else if rand.Intn(2) == 0 {
			str = "нет"
		}
		bot.Reply(str, update.Message)
		return true
	}
	return false
}
