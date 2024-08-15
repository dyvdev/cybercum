package tgbot

import (
	"encoding/json"
	"github.com/dyvdev/cybercum/swatter"
	"log"
	"os"
	"time"
)

const (
	dumpTick = time.Hour
	saveDump = "blobs/save.blob"
)

func (bot *Bot) Dumper(done <-chan bool) {
	ticker := time.NewTicker(dumpTick)
	go func() {
		for {
			select {
			case <-done:
				bot.SaveDump()
				bot.BotApi.StopReceivingUpdates()
				return
			case <-ticker.C:
				bot.SaveDump()
			}
		}
	}()
}

func (bot *Bot) SaveDump() {
	cfgJson, _ := json.Marshal(bot.Chats)
	err := os.WriteFile("chats.json", cfgJson, 0644)
	if err != nil {
		log.Fatal("Error during saving chats: ", err)
	}
	bot.Swatter.SaveDump(saveDump)
}

func (bot *Bot) LoadDump() {
	log.Println("reading chats...")
	content, err := os.ReadFile("chats.json")
	if err != nil {
		log.Println("Error when opening chats.json: ", err)
		cfgJson, _ := json.Marshal(bot.Chats)
		err := os.WriteFile("chats.json", cfgJson, 0644)
		if err != nil {
			log.Fatal("Error during saving chats: ", err)
		}
		log.Println("...created new")
		return
	}
	err = json.Unmarshal(content, &bot.Chats)
	if err != nil {
		log.Println("Empty chats or error during Unmarshal(): ", err)
		return
	}

	log.Println("reading dump...")
	var needToSave bool
	bot.Swatter, err = swatter.NewFromDump(saveDump)
	if err != nil {
		bot.Swatter, err = swatter.NewFromTextFile(bot.Cfg.DefaultDataFileName)
		if err != nil {
			log.Fatal("Error creating new swatter ", err)
		}
		needToSave = true
	}
	if needToSave {
		log.Println("saving dump...")
		bot.SaveDump()
	}
	//bot.FixChats()
	log.Println("reading chats...done")
}
