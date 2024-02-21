package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"modernc.org/mathutil"
	"strconv"
	"strings"
)

const (
	/*
		список команд для ОтцаБотов:
		add_cum - добавить кума
		enable_semen - включить генерацию сообщений(да/нет)
		enable_phrases - включить фиксированные фразы(да/нет)
		phrase - добавить фразу(через пробел после команды, вернёт номер фразы)
		phrase_remove - убрать фразу по номеру
		ratio - частота сообщений(50 значит, что бот будет писать раз в 50 сообщений)
		length - длина сгенерированных сообщений
	*/
	commandAddCum        = "add_cum"
	commandEnableSemen   = "enable_semen"
	commandEnablePhrases = "enable_phrases"
	commandFixed         = "phrase"
	commandFixedRemove   = "phrase_remove"
	commandRatio         = "ratio"
	commandLength        = "length"
)

func (bot *Bot) Commands(update tgbotapi.Update) {
	chat := bot.Chats[update.FromChat().ID]
	if bot.IsCum(update.Message) {
		switch update.Message.Command() {
		case commandAddCum:
			i := bot.AddCum(chat, update.Message.CommandArguments())
			bot.Reply("id:"+strconv.Itoa(i), update.Message)
			bot.SaveDump()
		case commandEnableSemen:
			chat.CanTalkSemen = update.Message.CommandArguments() == "да"
			bot.SaveDump()
		case commandEnablePhrases:
			chat.CanTalkPhrases = update.Message.CommandArguments() == "да"
			bot.SaveDump()
		case commandFixed:
			i := bot.AddFixedPhrase(chat, update.Message.CommandArguments())
			bot.Reply("id:"+strconv.Itoa(i), update.Message)
		case commandFixedRemove:
			id, err := strconv.Atoi(strings.TrimSpace(update.Message.CommandArguments()))
			if err != nil {
				bot.Reply(strconv.Itoa(len(chat.FixedPhrases)), update.Message)
			} else {
				i := bot.RemoveFixedPhrase(chat, id)
				bot.Reply("left:"+strconv.Itoa(i), update.Message)
				bot.SaveDump()
			}
		case commandRatio:
			ratio, err := strconv.Atoi(strings.TrimSpace(update.Message.CommandArguments()))
			if err != nil {
				bot.Reply(strconv.Itoa(chat.Ratio), update.Message)
			} else {
				chat.Ratio = ratio
				bot.SaveDump()
			}
		case commandLength:
			length, err := strconv.Atoi(strings.TrimSpace(update.Message.CommandArguments()))
			if err != nil {
				bot.Reply(strconv.Itoa(chat.SemenLength), update.Message)
			} else {
				chat.SemenLength = mathutil.Clamp(length, 1, maxLength)
				bot.SaveDump()
			}
		default:
			//bot.Reply("не понял" + update.Message.Command(), update.Message)
			log.Println(update.FromChat().UserName + " " + update.Message.From.UserName + " " + update.Message.Command())
		}
	}
}

func (bot *Bot) IsCum(message *tgbotapi.Message) bool {
	mmb, err := bot.BotApi.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID:             message.Chat.ID,
			SuperGroupUsername: "",
			UserID:             message.From.ID},
	})
	if err == nil {
		chat := bot.Chats[message.Chat.ID]
		for _, v := range chat.Cums {
			if v == mmb.User.UserName {
				return true
			}
		}
	}
	return false
}

func (bot *Bot) AddCum(chat *Chat, str string) int {
	if str != "" {
		str = strings.TrimPrefix(str, "@")
		chat.Cums = append(chat.Cums, str)
		bot.SaveDump()
		return len(chat.Cums) - 1
	}
	return -1
}
