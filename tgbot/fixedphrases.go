package tgbot

import (
	"encoding/json"
	tgbotapi "github.com/dyvdev/telegram-bot-api"
	"log"
	"math/rand"
	"os"
	"strings"
)

func (bot *Bot) SendFixedPhrase(message *tgbotapi.Message, searchPhrase string) bool {
	chat := bot.Chats[message.Chat.ID]
	txt := AnswerWithFixedPhrase(chat.Filename, searchPhrase)
	threadId := 0
	if txt == "" {
		return false
	}
	if message.Chat.IsForum && message.MessageThreadID != 0 {
		threadId = message.MessageThreadID
	}
	if strings.Contains(txt, "sticker:") {
		txt = strings.Replace(txt, "sticker:", "", 1)
		bot.SendMessage(tgbotapi.NewSticker(message.Chat.ID, tgbotapi.FileID(txt)))
	} else {
		bot.SendMessage(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           message.Chat.ID,
				MessageThreadID:  threadId,
				ReplyToMessageID: 0,
			},
			Text:                  txt,
			DisableWebPagePreview: false,
		})
	}
	return true
}

type Phrase struct {
	Chance  int
	Phrases []string
}

func LoadPhrases(filename string) (phrases map[string][]*Phrase) {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	err = json.Unmarshal(content, &phrases)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return
}

func AnswerWithFixedPhrase(filename string, text string) string {
	phrases := LoadPhrases(filename)
	if text == "" {
		return GetWeightedAnswer(phrases[""])
	}
	for key, _ := range phrases {
		if key != "" {
			keys := strings.Split(key, "|")
			for _, k := range keys {
				if strings.Contains(text, k) {
					return GetWeightedAnswer(phrases[k])
				}
			}
		}
	}
	return ""
}

func GetWeightedAnswer(phrases []*Phrase) (str string) {
	rnd := rand.Intn(100)
	for _, pp := range phrases {
		if rnd < pp.Chance {
			n := rand.Intn(len(pp.Phrases))
			return pp.Phrases[n]
		}
		rnd = rnd - pp.Chance
	}
	return ""
}
