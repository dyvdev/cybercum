package main

import (
    "main/semen"
    "main/tgbot"
    "time"
    "log"
    "io/ioutil"
    "encoding/json"
    "strconv"
)

type Config struct {
    BotId string
    BotName string
    TxtFile string
    SaveFile string
}
var cfg Config

func main() {
    loadcfg()
    //genSemen()
    runBot()
}

func loadcfg() {
    log.Println("reading config...")
    content, err := ioutil.ReadFile("./config.json")
    if err != nil {
        log.Fatal("Error when opening file: ", err)
    }
    err = json.Unmarshal(content, &cfg)
    if err != nil {
        log.Fatal("Error during Unmarshal(): ", err)
    }
}

func runBot() {
    log.Println("starting...")
    b := tgbot.NewBot(cfg.BotId, cfg.BotName, cfg.SaveFile)
    c := make(chan int)
    go saver(c, b)
    b.Update()
    i := <-c
    log.Println("exit ", i)
}

func genSemen() {
    s := semen.NewFromText(cfg.TxtFile)
    //s.ReadFile("bel.txt")
    for i:=0 ; i< 20; i++ {
        log.Println(s.Talk("", 5))
    }
    s.SaveDump(cfg.SaveFile)
}

func saver(c chan int, bot *tgbot.Bot) {
    for {
        time.Sleep(60 * 60 * time.Second)
        log.Println("saving.. [" + strconv.Itoa(len(bot.Semen)) + "]")
        bot.Semen.SaveDump(cfg.SaveFile)
    }
    c <- 1
}
