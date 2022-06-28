package tgbot

import (
    "log"
    "main/semen"
    "math/rand"
    "regexp"
    "encoding/json"
    "strings"
    "time"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
    BotApi *tgbotapi.BotAPI
    BotUsername string
    Semen semen.Semen
}

func NewBot(botId string, botName string, saveFile string) *Bot {
    if botId == "" || botName == "" || saveFile == "" {
        panic ("error creating new bot")
    }
    bot := Bot{}
    bapi, err := tgbotapi.NewBotAPI(botId)
    if err != nil {
        log.Fatal(err)
        return nil
    }
    bot.BotApi = bapi
    bot.Semen = *semen.NewFromDump(saveFile)
    bot.BotUsername = botName
    return &bot
}

func (bot Bot) Update() {
    ratio:= 5
    pauseSeconds:=15 * time.Second
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60
    updates := bot.BotApi.GetUpdatesChan(u)
    next:= time.Now().UTC()
    for update := range updates {
        if update.Message != nil && update.Message.Text != "" {

            if false {
                log.Printf("[%d] [%s] %s", update.UpdateID, update.Message.From.UserName, update.Message.Text)
                out, err := json.Marshal(update.Message)
                if err != nil {
                    panic (err)
                }
                log.Println(string(out))
                //log.Printf("ready %t mention %t reply %t u: %s", isReady, isMention, isReply, bot.BotUsername)
            }
            isReply := update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.UserName == bot.BotUsername
            isMention := strings.Contains(update.Message.Text, "@" + bot.BotUsername)
            isReady := time.Now().UTC().After(next)
            isTimeToTalk := rand.Intn(100) < ratio && isReady
            if  isReady {
                next = time.Now().UTC().Add(pauseSeconds)
            }

            if  isMention || isReply || isTimeToTalk {
                reg, err := regexp.Compile(`\p{Cyrillic}+`)
                if err != nil {
                    log.Fatal(err)
                }
                words := reg.FindAllString(update.Message.Text, -1)
                if len(words) != 0 {
                    bot.Semen.Learning(words)
                    msg := tgbotapi.NewMessage(
                        update.Message.Chat.ID,
                        //CHAT_ID,
                         bot.Semen.Talk(words[rand.Intn(len(words))], 5))
                    if !isTimeToTalk {
                        msg.ReplyToMessageID = update.Message.MessageID
                    }
                    bot.BotApi.Send(msg)
                }
            }
        }
    }
}
