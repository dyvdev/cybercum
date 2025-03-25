package tgbot

import (
	tgbotapi "github.com/dyvdev/telegram-bot-api"
	"log"
)

func (bot *Bot) NewUserMessageRemover(update tgbotapi.Update) {
	if update.Message != nil && update.Message.NewChatMembers != nil {
		if !bot.CheckAdminRights(update.FromChat().ID) {
			return
		}
		_, err := bot.BotApi.Send(tgbotapi.NewDeleteMessage(update.FromChat().ID, update.Message.MessageID))
		if err != nil {
			log.Println("error removing message for new user")
		}
	}
}
func (bot *Bot) CheckAdminRights(chatId int64) bool {
	mmb, err := bot.BotApi.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID:             chatId,
			SuperGroupUsername: "",
			UserID:             bot.BotApi.Self.ID},
	})
	if err == nil && (mmb.Status == "administrator" || mmb.Status == "creator") {
		return true
	}
	return false
}
