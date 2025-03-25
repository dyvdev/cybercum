package tgbot

import (
	"encoding/json"
	tgbotapi "github.com/dyvdev/telegram-bot-api"
	"log"
	"math/rand"
	"os"
	"strings"
)

func (bot *Bot) SendFixedPhrase(message *tgbotapi.Message, text string) bool {
	chat := bot.Chats[message.Chat.ID]
	txt := AnswerWithFixedPhrase(chat.Filename, text)
	if txt == "" {
		return false
	}
	if strings.Contains(txt, "sticker:") {
		txt = strings.Replace(txt, "sticker:", "", 1)
		bot.SendMessage(tgbotapi.NewSticker(message.Chat.ID, tgbotapi.FileID(txt)))
	} else {
		if text == "" {
			bot.SendText(txt, message)
		} else {
			bot.Reply(txt, message)
		}
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
		log.Fatal("Error when opening phrases file: ", filename, err)
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
