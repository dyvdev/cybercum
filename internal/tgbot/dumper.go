package tgbot

import (
	"encoding/json"
	"github.com/dyvdev/cybercum/internal/swatter"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	dumpTick   = time.Hour
	saveFolder = "blobs/"
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
	for key, _ := range bot.Chats {
		bot.Swatter[key].SaveDump(saveFolder + strconv.Itoa(int(key)) + ".blob")
	}
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
	for key, chat := range bot.Chats {
		log.Println("reading dump... " + saveFolder + strconv.Itoa(int(key)) + ".blob for chat [" + chat.ChatName + "]")

		bot.Swatter[key], err = swatter.NewFromDump(saveFolder + strconv.Itoa(int(key)) + ".blob")
		if err != nil {
			bot.Swatter[key], err = swatter.NewFromTextFile(bot.Cfg.DefaultDataFileName)
			if err != nil {
				log.Fatal("Error creating new swatter ", err)
			}
			needToSave = true
		}
	}
	if needToSave {
		log.Println("saving dump...")
		bot.SaveDump()
	}
	//bot.FixChats()
	log.Println("reading chats...done")
}
