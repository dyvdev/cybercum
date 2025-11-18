package tgbot

import (
	"log"
	"math/rand"
	"strconv"
	"strings"

	tgbotapi "github.com/dyvdev/telegram-bot-api"
	"modernc.org/mathutil"
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
		roll - написать число от 0 до N
	*/
	commandAddCum        = "add_cum"
	commandEnableSemen   = "enable_semen"
	commandEnablePhrases = "enable_phrases"
	commandFixed         = "phrase"
	commandFixedRemove   = "phrase_remove"
	commandRatio         = "ratio"
	commandLength        = "length"
	commandRoll          = "roll"
)

func (bot *Bot) Commands(update tgbotapi.Update) {
	chat := bot.Chats[update.FromChat().ID]
	switch update.Message.Command() {
	case commandRoll:
		n, err := strconv.Atoi(strings.TrimSpace(update.Message.CommandArguments()))
		if err == nil && n > 0 {
			bot.Reply(strconv.Itoa(rand.Intn(n)), update.Message)
		} else {
			bot.Reply(strconv.Itoa(rand.Intn(100000000)), update.Message)
		}
	}
	if bot.IsCum(update.Message.Chat.ID, update.Message.From.ID) {
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
		case commandFixed:
		case commandFixedRemove:
		default:
			//bot.Reply("не понял" + update.Message.Command(), update.Message)
			log.Println(update.FromChat().UserName + " " + update.Message.From.UserName + " " + update.Message.Command())
		}
	}
}

func (bot *Bot) IsCum(chatId int64, userId int64) bool {
	mmb, err := bot.BotApi.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID:             chatId,
			SuperGroupUsername: "",
			UserID:             userId},
	})
	if err == nil {
		chat := bot.Chats[chatId]
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
