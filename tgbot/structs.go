package tgbot

import (
	"time"

	"github.com/dyvdev/cybercum/swatter"
	tgbotapi "github.com/dyvdev/telegram-bot-api"
)

const (
	maxLength = 100
)

type Config struct {
	BotId   string // айди бота от ОтцаБОтов
	MainCum string // ник владельца

	EnablePhrases          bool   // включить фиксированные фразы
	DefaultPhrasesFilename string // список фраз

	EnableSemen     bool // включить генерацию фраз
	EnableReactions bool // включить реакции
	EnableNeuro     bool // включить нейронку
	EnableGame      bool // включить игру
	Ratio           int  // количество сообщений между ответами бота
	Length          int  // длина сообщений генератоа цепей
}
type Chat struct {
	ChatName         string
	CanTalkSemen     bool
	CanTalkNeuro     bool
	CanTalkPhrases   bool
	CanSendReactions bool
	Ratio            int //количество сообщений между ответами бота
	Counter          int //счетчик сообщений в чате
	SemenLength      int
	Filename         string
	Cums             []string
	Context          []string
	NeuroPrompt      string

	Gamers      map[int64]*Gamer
	CanPlayGame bool
}

type Bot struct {
	BotApi *tgbotapi.BotAPI
	Timer  time.Time
	Pause  time.Duration
	Cfg    Config

	Chats map[int64]*Chat

	Swatter *swatter.DataStorage

	GamingChan   chan tgbotapi.Update
	CurrentGames map[int64]*Game
}
