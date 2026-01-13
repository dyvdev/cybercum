package tgbot

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/dyvdev/cybercum/swatter"
)

const (
	dumpTick        = time.Hour
	saveDump        = "blobs/save.blob"
	saveDumpDefault = "blobs/default.blob"
)

func (bot *Bot) Dumper(wg *sync.WaitGroup, done <-chan bool) {
	ticker := time.NewTicker(dumpTick)
	go func() {
		defer wg.Done()
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
	cfgJson, err := json.Marshal(bot.Chats)
	if err != nil {
		log.Println("error saving dump: ", err)
		return
	}
	err = os.WriteFile("chats.json", cfgJson, 0644)
	if err != nil {
		log.Fatal("Error during saving chats: ", err)
	}
	bot.Swatter.SaveDump(saveDump)
}

func (bot *Bot) LoadDump() error {
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
		return err
	}
	err = json.Unmarshal(content, &bot.Chats)
	if err != nil {
		log.Println("Empty chats or error during Unmarshal(): ", err)
		return err
	}

	log.Println("reading dump...")
	var needToSave bool
	bot.Swatter, err = swatter.NewFromDump(saveDump)
	if err != nil {
		bot.Swatter, err = swatter.NewFromDump(saveDumpDefault)
		if err != nil {
			log.Println("Error creating new swatter ", err)
		} else {
			needToSave = true
		}
	}
	if needToSave {
		log.Println("saving dump...")
		bot.SaveDump()
	}
	//bot.FixChats()
	log.Println("reading chats...done")
	return nil
}
