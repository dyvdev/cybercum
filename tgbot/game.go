package tgbot

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	tgbotapi "github.com/dyvdev/telegram-bot-api"
)

type Gamer struct {
	Streak int
	Wins   int
	Loses  int
	Score  int
}
type CurrentGamer struct {
	Id    int64
	Stake int
}
type Game struct {
	MessageId     int
	Started       bool
	FirstPlayerId int64
	Stake         int
}

func (bot *Bot) StartGaming() {
	bot.CurrentGames = map[int64]*Game{}
	bot.GamingChan = make(chan tgbotapi.Update)
	go func() {
		for {
			select {
			case u := <-bot.GamingChan:
				game, ok := bot.CurrentGames[u.FromChat().ID]
				if ok {
					if u.CallbackQuery != nil {
						log.Println("msg id:", bot.CurrentGames[u.FromChat().ID].MessageId)
						if bot.CurrentGames[u.FromChat().ID].FirstPlayerId == u.CallbackQuery.From.ID {
							break
						}
						if bot.gameAccept(u, game) {
							// если игра закончилась, почистим
							delete(bot.CurrentGames, u.FromChat().ID)
						} else {
							go func() {
								time.Sleep(30 * time.Minute)
								//time.Sleep(5 * time.Second)
								game, ok := bot.CurrentGames[u.FromChat().ID]
								if ok {
									if bot.gameAccept(tgbotapi.Update{
										CallbackQuery: &tgbotapi.CallbackQuery{
											ID:   "1",
											From: &bot.BotApi.Self,
											Data: strconv.Itoa(rand.Intn(4) + 1),
											Message: &tgbotapi.Message{
												MessageID: game.MessageId,
												Chat:      u.FromChat(),
											},
										},
									}, game) {
										delete(bot.CurrentGames, u.FromChat().ID)
									}
								}
							}()
						}
					}
				} else {
					if u.Message != nil {
						bot.CurrentGames[u.FromChat().ID] = &Game{
							MessageId:     u.Message.MessageID,
							Started:       false,
							FirstPlayerId: 0,
						}
						bot.newGameInvite(u, bot.CurrentGames[u.FromChat().ID])
					}
				}
			}
		}
	}()
}
func (bot *Bot) GameUpdate(update tgbotapi.Update) {
	bot.GamingChan <- update
}

func (bot *Bot) gameAccept(update tgbotapi.Update, currentGame *Game) bool {
	chat := bot.Chats[update.FromChat().ID]
	stake, _ := strconv.Atoi(update.CallbackQuery.Data)
	gamerId := update.CallbackQuery.From.ID
	if chat.Gamers == nil {
		chat.Gamers = map[int64]*Gamer{}
	}
	_, ok := chat.Gamers[gamerId]
	if !ok {
		chat.Gamers[gamerId] = &Gamer{
			Wins:  0,
			Loses: 0,
		}
	}
	if !currentGame.Started {
		currentGame.FirstPlayerId = gamerId
		currentGame.Stake = stake
		currentGame.Started = true
		keyboard := tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{{
				tgbotapi.NewInlineKeyboardButtonData("✂️", "1"),
				tgbotapi.NewInlineKeyboardButtonData("🪨", "2"),
				tgbotapi.NewInlineKeyboardButtonData("🧻", "3")}},
		}
		for i := range keyboard.InlineKeyboard[0] {
			j := rand.Intn(i + 1)
			keyboard.InlineKeyboard[0][i], keyboard.InlineKeyboard[0][j] = keyboard.InlineKeyboard[0][j], keyboard.InlineKeyboard[0][i]
		}
		msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, fmt.Sprintf("%s\nбросает вызов чату, выбери своё оружие и сразись!",
			GetPlayerString(update.CallbackQuery.From, chat.Gamers[gamerId])),
			keyboard)
		c, err := bot.BotApi.Send(msg)
		if err != nil {
			log.Println("Error sending message: ", err)
		}
		currentGame.MessageId = c.MessageID
		return false
	}

	winnerId := currentGame.FirstPlayerId
	loserId := gamerId
	winner, err := bot.BotApi.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID:             update.FromChat().ID,
			SuperGroupUsername: "",
			UserID:             winnerId},
	})
	if err != nil {
		return false
	}
	loser, err := bot.BotApi.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID:             update.FromChat().ID,
			SuperGroupUsername: "",
			UserID:             loserId},
	})
	if err != nil {
		return false
	}
	log.Println("1 ", GetName(winner.User), chat.Gamers[winnerId], currentGame.Stake)
	log.Println("2 ", GetName(loser.User), chat.Gamers[loserId], stake)
	winSmile := GetStake(currentGame.Stake)
	loseSmile := GetStake(stake)
	beatenStreak := 0
	if stake == currentGame.Stake {
		keyboard := tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{{
				tgbotapi.NewInlineKeyboardButtonData("✂️", "1"),
				tgbotapi.NewInlineKeyboardButtonData("🪨", "2"),
				tgbotapi.NewInlineKeyboardButtonData("🧻", "3")}},
		}
		for i := range keyboard.InlineKeyboard[0] {
			j := rand.Intn(i + 1)
			keyboard.InlineKeyboard[0][i], keyboard.InlineKeyboard[0][j] = keyboard.InlineKeyboard[0][j], keyboard.InlineKeyboard[0][i]
		}
		msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, fmt.Sprintf("%s\nразошлись миром 🐷%s🐔\n%s\nЖмакай кнопки, чтобы сыграть заново",
			GetPlayerString(winner.User, chat.Gamers[winnerId]),
			GetStake(stake),
			GetPlayerString(loser.User, chat.Gamers[loserId])),
			keyboard)
		c, err := bot.BotApi.Send(msg)
		if err != nil {
			log.Println("Error sending message: ", err)
		}
		currentGame.MessageId = c.MessageID
		currentGame.Stake = 0
		currentGame.Started = false
		currentGame.FirstPlayerId = 0
		return false
	} else if (stake == 3 && currentGame.Stake == 2) || (stake == 2 && currentGame.Stake == 1) || (stake == 1 && currentGame.Stake == 3) {
		winnerId, loserId = gamerId, currentGame.FirstPlayerId
		winner, loser = loser, winner
		winSmile, loseSmile = loseSmile, winSmile
		chat.Gamers[gamerId].Wins++
		chat.Gamers[gamerId].Streak++
		if chat.Gamers[gamerId].Score < chat.Gamers[gamerId].Streak {
			chat.Gamers[gamerId].Score = chat.Gamers[gamerId].Streak
		}
		chat.Gamers[currentGame.FirstPlayerId].Loses++
		beatenStreak = chat.Gamers[currentGame.FirstPlayerId].Streak
		chat.Gamers[currentGame.FirstPlayerId].Streak = 0
	} else {
		chat.Gamers[currentGame.FirstPlayerId].Wins++
		chat.Gamers[currentGame.FirstPlayerId].Streak++
		if chat.Gamers[currentGame.FirstPlayerId].Score < chat.Gamers[currentGame.FirstPlayerId].Streak {
			chat.Gamers[currentGame.FirstPlayerId].Score = chat.Gamers[currentGame.FirstPlayerId].Streak
		}
		chat.Gamers[gamerId].Loses++
		beatenStreak = chat.Gamers[gamerId].Streak
		chat.Gamers[gamerId].Streak = 0
	}
	log.Println("winner ", GetName(winner.User), chat.Gamers[winnerId])
	log.Println("loser ", GetName(loser.User), chat.Gamers[loserId])

	actions := []string{
		"мягко ляпает лапкой",
		"угробил ладошкой",
		"мутузит писюном",
		"ставит подножку",
		"дает отеческого леща",
		"вероломно нападает с тыла",
		"наносит удар в псину",
		"даёт щелбан",
		"пробивает лося",
		"тыкает пальчиком в пупок",
		"шлёпает по попе",
		"заводит за щеку",
		"яростно квокает на",
		"грозит пальчиком",
		"рякает на",
		"делает крапивку",
		"вдувает по самые помидоры",
		"проводит славянский зажим яйцами",
	}
	msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, fmt.Sprintf("🐔%s%s\n%s\n🐷%s%s",
		GetPlayerString(winner.User, chat.Gamers[winnerId]),
		winSmile,
		actions[rand.Intn(len(actions)-1)],
		GetPlayerString(loser.User, chat.Gamers[loserId]),
		loseSmile))
	if beatenStreak > 2 {
		msg.Text += "\nи заканчивает его серию из " + strconv.Itoa(beatenStreak) + " побед"
	}
	switch chat.Gamers[winnerId].Streak {
	case 1:
	case 2:
	case 3:
		msg.Text += "\n😱Уже третья победа подряд!😱"
	case 4:
		msg.Text += "\n😱Ладно, шутки шутками, но 4 победы подряд?!😱"
	case 5:
		msg.Text += "\n😱ИИИиии пятый фраг подряд зарабатывает в свою копилку молодой игрок на 🐔!😱"
	case 6:
		msg.Text += "\n😱ШЕСТЬ ПОБЕД ПОДРЯД! 🐔 НЕ ОСТАНОВИТЬ!😱"
	case 7:
		msg.Text += "\n😱7(семь) ПОБЕД.. чтоооОО?!😱"
	case 8:
		msg.Text += "\n😱Восьмая победа подряд! (хотите так же? ссылка на донат в описании)😱"
	case 9:
		msg.Text += "\n😱Никто в это не верил и вот - ДЕВЯТАЯ ПОБЕДА ПОДРЯД!😱"
	case 10:
		msg.Text += "\n10 wins in a row, please contact cums with a bug report"
	case 11:
		msg.Text += "\n11 побед? или сколько? я со счету сбился.."
	default:
		msg.Text += "\n" + strconv.Itoa(chat.Gamers[winnerId].Streak) + " победа.. кому не пофиг?"
	}
	_, err = bot.BotApi.Send(msg)
	if err != nil {
		log.Println("Error sending message: ", err)
		log.Println("msg: ", msg.MessageID)
	}
	return true
}

func (bot *Bot) newGameInvite(update tgbotapi.Update, currentGame *Game) {
	msg := tgbotapi.NewMessage(update.FromChat().ID, "Жмякай кнопку, чтобы принять участие в Игре!")
	keyboard := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{{
			tgbotapi.NewInlineKeyboardButtonData("✂️", "1"),
			tgbotapi.NewInlineKeyboardButtonData("🪨", "2"),
			tgbotapi.NewInlineKeyboardButtonData("🧻", "3")}},
	}
	for i := range keyboard.InlineKeyboard[0] {
		j := rand.Intn(i + 1)
		keyboard.InlineKeyboard[0][i], keyboard.InlineKeyboard[0][j] = keyboard.InlineKeyboard[0][j], keyboard.InlineKeyboard[0][i]
	}
	msg.ReplyMarkup = keyboard
	c, err := bot.BotApi.Send(msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
	currentGame.MessageId = c.MessageID
}

func GetName(user *tgbotapi.User) string {
	name := user.FirstName + " "
	if user.LastName != "" {
		name += user.LastName + " "
	}
	if user.UserName != "" {
		name += "@" + user.UserName + " "
	}
	return name
}
func GetPlayerString(user *tgbotapi.User, g *Gamer) string {
	return GetName(user) + "(" + strconv.Itoa(g.Wins) + "💰" + strconv.Itoa(g.Loses) + "⚰️)"
}
func GetStake(i int) string {
	if i == 1 {
		return "✂️"
	}
	if i == 2 {
		return "🪨"
	}
	if i == 3 {
		return "🧻"
	}
	return ""
}
