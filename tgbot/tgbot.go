package tgbot

import (
	"encoding/json"
	"github.com/dyvdev/cybercum/swatter"
	"github.com/dyvdev/cybercum/utils"
	tgbotapi "github.com/dyvdev/telegram-bot-api"
	"log"
	"math/rand"
	"os"
	"sync/atomic"
	"time"
)

// NewBot конструктор нового бота, загрузка данных чатов
// NewBot конструктор нового бота, загрузка данных чатов
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
	if err = bot.LoadDump(); err != nil {
		return nil
	}
	bot.DropCounter()
	bot.Pause = 15 * time.Second
	bot.Timer = time.Now().UTC().Add(bot.Pause)
	return &bot
}

// Update обработка событий
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
				// админ
				bot.NewUserMessageRemover(update)

				//========= раскоментировать для дебага:
				//bot.CheckChatSettings(update)
				//continue
				//=========

				bot.CheckChatSettings(update)
				if update.MessageReaction != nil {
					bot.ProcessReaction(update)
				}
				if update.Message != nil {
					if update.FromChat().IsPrivate() {
						logMessage(update.Message)
						msg := bot.GenerateMessage(update.Message)
						if msg != nil {
							bot.SendMessage(msg)
						}
					} else if update.Message.Text != "" {
						if update.Message.IsCommand() {
							bot.Commands(update)
						} else {
							bot.ProcessMessage(update)
						}
					}
					if update.Message.Photo != nil {
						bot.SendPhotoReaction(update)
					}
				}
			}
		}
	}()
}

// ProcessMessage обработка сообщений
func (bot *Bot) ProcessMessage(update tgbotapi.Update) {
	chat := bot.Chats[update.FromChat().ID]
	chat.Counter++
	isTimeToTalk := chat.Ratio == 0 || (chat.Counter > chat.Ratio && bot.Tick())
	if chat.CanSendReactions && bot.SendRandomReaction(update) {
		return
	}
	if utils.CheckForUrls(update.Message) {
		return
	}
	if isTimeToTalk && chat.CanTalkPhrases && bot.SendFixedPhrase(update.Message, "") {
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
			if bot.SendAnswer(update) {
				return
			} else if bot.Shakspearing(update) {
				return
			} else {
				bot.SemenMessageSend(update, isReply)
			}
		}
		bot.Learning(update.Message)
	}
}

// Shakspearing попытка скаламбурить, true если получилось
func (bot *Bot) Shakspearing(update tgbotapi.Update) bool {
	if rand.Intn(10) == 1 {
		txt := utils.CleanText(update.Message.Text)
		// попытаемся скаламбурить
		txt = shakeSpear(txt)
		if txt == "" {
			return false
		}
		msg := tgbotapi.NewMessage(update.FromChat().ID, txt+"😁")
		msg.ReplyToMessageID = update.Message.MessageID
		bot.BotApi.Send(msg)
		return true
	}
	return false
}

// ProcessReaction обработка реакций на сообщения бота
func (bot *Bot) ProcessReaction(update tgbotapi.Update) {
	if update.MessageReaction.NewReaction != nil && update.MessageReaction.NewReaction[0].Emoji == "❤" {
		log.Println("reaction message: ", update.Message)
	}
}

// SemenMessageSend отправка генерируемых сообщений
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

// GenerateMessage генерация сообщения
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
}

// Learning обучение
func (bot *Bot) Learning(message *tgbotapi.Message) {
	bot.Swatter.ParseText(message.Text)
}

// SendMessage отправка сообщения
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

// Reply отправка ответа на сообщение
func (bot *Bot) Reply(text string, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID
	bot.SendMessage(msg)
}

// Tick таймер
func (bot *Bot) Tick() bool {
	isReady := time.Now().UTC().After(bot.Timer)
	if isReady {
		bot.Timer = time.Now().UTC().Add(bot.Pause)
	}
	return isReady
}

// LoadConfig загрузка конфига
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

// SaveConfig сохранение конфига
func (bot *Bot) SaveConfig() {
	log.Println("saving config...")
	cfgJson, _ := json.Marshal(bot.Cfg)
	err := os.WriteFile("config.json", cfgJson, 0644)
	if err != nil {
		log.Fatal("Error during saving config: ", err)
	}
	log.Println("saving config...done")
}

// DropCounter сброс счетчика для всех чатов
func (bot *Bot) DropCounter() {
	for _, c := range bot.Chats {
		c.Counter = 0
	}
}

// CheckChatSettings проверка чата на предмет кривого названия или установка настроек по умолчанию для нового чата
func (bot *Bot) CheckChatSettings(update tgbotapi.Update) {
	_, ok := bot.Chats[update.FromChat().ID]
	chatName := update.FromChat().Title
	if chatName == "" {
		chatName = update.FromChat().UserName
	}
	// если впервые в чате, зададим дефолтный конфиг
	if !ok {
		log.Println("new chat: ", update.FromChat().ID, chatName)
		bot.Chats[update.FromChat().ID] = &Chat{
			ChatName:         chatName,
			CanTalkSemen:     bot.Cfg.EnableSemen,
			CanTalkPhrases:   bot.Cfg.EnablePhrases,
			CanSendReactions: bot.Cfg.EnableReactions,
			Ratio:            bot.Cfg.Ratio,
			Counter:          0,
			SemenLength:      bot.Cfg.Length,
			Filename:         bot.Cfg.DefaultPhrasesFilename,
			Cums:             []string{bot.Cfg.MainCum},
			lastMessageId:    atomic.Uint64{},
		}
		bot.SaveDump()
	}
	bot.Chats[update.FromChat().ID].ChatName = chatName
}

func logMessage(message *tgbotapi.Message) {
	chatName := message.Chat.Title
	if chatName == "" {
		chatName = message.Chat.UserName
	}
	log.Println("message log start")
	log.Println("chatname: ", chatName)
	if message.Text != "" {
		log.Println("text: ", message.Text)
	}
	if message.Sticker != nil {
		log.Println("sticker id: ", message.Sticker.FileID)
	}
	if message.Audio != nil {
		log.Println("sticker id: ", message.Audio.FileID)
	}
	log.Println("message log end")
}
