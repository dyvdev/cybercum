package tgbot

import (
	"encoding/json"
	"github.com/dyvdev/cybercum/swatter"
	tgbotapi "github.com/dyvdev/telegram-bot-api"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"sync/atomic"
	"time"
)

const (
	maxLength = 100
	nefren    = "CAACAgIAAx0CTK3KYQACAQNjDKmYViPp5K-PWxuUKUDpwg0vQQAC9hEAAqx6iEqOhkQYAe2vbSkE"
)

type Config struct {
	BotId   string // айди бота от ОтцаБОтов
	MainCum string // ник владельца

	EnablePhrases  bool     // включить фиксированные фразы
	DefaultPhrases []string // список фраз

	EnableSemen         bool   // включить генерацию фраз
	Ratio               int    // количество сообщений между ответами бота
	Length              int    // длина сообщений генератоа цепей
	DefaultDataFileName string // текстовый файл из которого берутся базовые данные
}
type Chat struct {
	ChatName       string
	CanTalkSemen   bool
	CanTalkPhrases bool
	Ratio          int //количество сообщений между ответами бота
	Counter        int //счетчик сообщений в чате
	SemenLength    int
	FixedPhrases   []string
	Cums           []string
	lastMessageId  atomic.Uint64
}

type Bot struct {
	BotApi *tgbotapi.BotAPI
	Timer  time.Time
	Pause  time.Duration
	Cfg    Config

	Chats map[int64]*Chat

	Swatter *swatter.DataStorage
}

func NewBot() *Bot {
	bot := Bot{}
	bot.LoadConfig()
	log.Println(bot.Cfg)
	if bot.Cfg.BotId == "" {
		panic("error creating new bot")
	}
	bapi, err := tgbotapi.NewBotAPI(bot.Cfg.BotId)
	if err != nil {
		log.Println("id: ", bot.Cfg.BotId)
		log.Fatal("starting tg bot error: ", err)
		return nil
	}
	bot.BotApi = bapi
	bot.Swatter = &swatter.DataStorage{}
	bot.Chats = map[int64]*Chat{}
	bot.LoadDump()
	bot.Pause = 15 * time.Second
	bot.Timer = time.Now().UTC().Add(bot.Pause)
	return &bot
}

func (bot *Bot) Update(done <-chan bool) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{"message", "message_reaction", "message_reaction_count"}
	updates := bot.BotApi.GetUpdatesChan(u)
	ticker := time.NewTicker(25 * time.Millisecond)
	go func() {
		for update := range updates {
			select {
			case <-done:
				bot.SaveDump()
				return
			case <-ticker.C:
				if update.FromChat().Type == "channel" {
					log.Println("channel update!")
					log.Println(update.FromChat())
					bot.BotApi.Leave(update.FromChat())
					continue
				}
				//continue
				if update.MessageReaction != nil {
					bot.ProcessReaction(update)
				}
				if update.Message != nil {
					if update.Message.Photo != nil && rand.Intn(15) == 1 {
						bot.SendPhotoReaction(update)
					} else if update.Message.Text != "" {
						bot.CheckChatSettings(update)
						if update.Message.IsCommand() {
							bot.Commands(update)
						} else {
							bot.ProcessMessage(update)
						}
					}
				}
			}
		}
	}()
}

func (bot *Bot) ProcessMessage(update tgbotapi.Update) {
	chat := bot.Chats[update.FromChat().ID]
	chat.Counter++
	isTimeToTalk := chat.Ratio == 0 || (chat.Counter > chat.Ratio && bot.Tick()) //|| bot.IsCum(update.Message)
	if update.FromChat().IsPrivate() {
		msg := bot.GenerateMessage(update.Message)
		if msg == nil {
			return
		}
		bot.SendMessage(msg)
		return
	}

	if isTimeToTalk && chat.CanTalkPhrases {
		bot.SendFixedPhrase(update.Message)
		chat.Counter = 0
	} else if chat.CanTalkSemen {
		isReply := update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.UserName == bot.BotApi.Self.UserName
		isMessageToMe := bot.BotApi.IsMessageToMe(*update.Message)
		if isTimeToTalk || isReply || isMessageToMe {
			// всегда отвечаем на вопрос к нам
			if (isMessageToMe || isReply) && bot.SendAnswer(update) {
				return
			}
			chat.Counter = 0
			if rand.Intn(10) == 1 {
				// если есть вопрос, ответим
				if bot.SendAnswer(update) {
					return
				}
				txt := strings.ToLower(regexp.MustCompile(`\.|,|;|!|\?`).ReplaceAllString(update.Message.Text, ""))
				// попытаемся скаламбурить
				txt = shakeSpear(txt)
				if txt == "" {
					// если не вышло, просто генерим фразу как обычно
					bot.SemenMessageSend(update, isReply)
					return
				}
				msg := tgbotapi.NewMessage(update.FromChat().ID, txt+"😁")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.BotApi.Send(msg)
			} else {
				bot.SemenMessageSend(update, isReply)
			}
			bot.Learning(update.Message)
		} else {
			bot.SendRandomReaction(update)
		}
	}
}

func (bot *Bot) ProcessReaction(update tgbotapi.Update) {
	if update.MessageReaction.NewReaction != nil && update.MessageReaction.NewReaction[0].Emoji == "❤" {
		//
	}
}
func (bot *Bot) SemenMessageSend(update tgbotapi.Update, isReply bool) {
	msg := bot.GenerateMessage(update.Message)
	if msg == nil {
		return
	}
	if isReply { // тут конверт доделать
		switch concrete := msg.(type) {
		case tgbotapi.MessageConfig:
			concrete.ReplyToMessageID = update.Message.MessageID
			bot.SendMessage(concrete)
		case tgbotapi.StickerConfig:
			concrete.ReplyToMessageID = update.Message.MessageID
			bot.SendMessage(concrete)
		default:
			log.Println("ошибка")
		}
	} else {
		bot.SendMessage(msg)
	}
}

func (bot *Bot) GenerateMessage(message *tgbotapi.Message) tgbotapi.Chattable {
	msg := bot.Swatter.GenerateText(message.Text, bot.Chats[message.Chat.ID].SemenLength)
	threadId := 0
	if message.Chat.IsForum && message.MessageThreadID != 0 {
		threadId = message.MessageThreadID
	}
	return tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           message.Chat.ID,
			MessageThreadID:  threadId,
			ReplyToMessageID: 0,
		},
		Text:                  msg,
		DisableWebPagePreview: false,
	}
	//else {
	//    return tgbotapi.NewSticker(message.Chat.ID, tgbotapi.FileID(nefren))
	//}
}

func (bot *Bot) Learning(message *tgbotapi.Message) {
	bot.Swatter.ParseText(message.Text)
}

func (bot *Bot) SendMessage(message tgbotapi.Chattable) {
	switch concrete := message.(type) {
	case tgbotapi.MessageConfig:
		if concrete.Text == "" {
			return
		}
	}
	_, err := bot.BotApi.Send(message)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}

func (bot *Bot) Reply(text string, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID
	bot.SendMessage(msg)
}

func (bot *Bot) ReplyNefren(message *tgbotapi.Message) {
	msg := tgbotapi.NewSticker(message.Chat.ID, tgbotapi.FileID(nefren))
	msg.ReplyToMessageID = message.MessageID
	bot.SendMessage(msg)
}

func (bot *Bot) Tick() bool {
	isReady := time.Now().UTC().After(bot.Timer)
	if isReady {
		bot.Timer = time.Now().UTC().Add(bot.Pause)
	}
	return isReady
}

func (bot *Bot) LoadConfig() {
	log.Println("reading config...")
	content, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	err = json.Unmarshal(content, &bot.Cfg)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	log.Println("reading config...done")
}

func (bot *Bot) SaveConfig() {
	log.Println("saving config...")
	cfgJson, _ := json.Marshal(bot.Cfg)
	err := os.WriteFile("config.json", cfgJson, 0644)
	if err != nil {
		log.Fatal("Error during saving config: ", err)
	}
	log.Println("saving config...done")
}

func (bot *Bot) FixChats() {
	for id, c := range bot.Chats {
		chat, err := bot.BotApi.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: id, SuperGroupUsername: ""}})
		if err == nil {
			//log.Println(chat)
			//if chat.IsPrivate() {
			//	log.Println("deleting " + c.ChatName)
			//	delete(bot.Chats, id)
			//}
			c.ChatName = chat.Title
			if c.ChatName == "" {
				c.ChatName = chat.UserName
			}
			if c.ChatName == "" {
				c.ChatName = chat.FirstName + " " + chat.LastName
			}
			log.Println("title " + c.ChatName)
		} else {
			log.Println("deleting err ")
			delete(bot.Chats, id)
		}
	}
	bot.SaveDump()
}

func (bot *Bot) Clean() {
	bot.Swatter.Clean()
}

func (bot *Bot) CheckChatSettings(update tgbotapi.Update) {
	_, ok := bot.Chats[update.FromChat().ID]
	// если впервые в чате, зададим дефолтный конфиг
	if !ok {
		log.Println("new chat: ", update.FromChat().ID)
		chatName := update.FromChat().Title
		if chatName == "" {
			chatName = update.FromChat().UserName
		}
		bot.Chats[update.FromChat().ID] = &Chat{
			ChatName:       chatName,
			CanTalkSemen:   bot.Cfg.EnableSemen,
			CanTalkPhrases: bot.Cfg.EnablePhrases,
			Ratio:          bot.Cfg.Ratio,
			Counter:        0,
			SemenLength:    bot.Cfg.Length,
			FixedPhrases:   bot.Cfg.DefaultPhrases,
			Cums:           []string{bot.Cfg.MainCum},
			lastMessageId:  atomic.Uint64{},
		}
		var err error
		bot.Swatter, err = swatter.NewFromTextFile(bot.Cfg.DefaultDataFileName)
		if err != nil {
			log.Fatal("Error creating new swatter ", err)
		}
		bot.SaveDump()
	}
	bot.Chats[update.FromChat().ID].ChatName = update.FromChat().Title
}
