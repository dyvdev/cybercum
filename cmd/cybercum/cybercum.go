package main

import (
	"flag"
	cum "github.com/dyvdev/cybercum/internal"
	"github.com/dyvdev/cybercum/internal/config"
	"github.com/dyvdev/cybercum/internal/swatter"
	"github.com/dyvdev/cybercum/internal/utils"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	cfgPath := flag.String("c", "", "path to config file")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	//cum.ReadBot()
	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		log.Fatalf("cant load config: %s", err.Error())
	}

	cum.RunBot(cfg)
	//cum.CleanBot()
	//test()
	//testChatHistoryGen()
}

func testChatHistoryGen() {
	sw := &swatter.DataStorage{}
	data := utils.GetTgData("tghistory.json")
	file, err := os.Create("tghistory.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	for _, str := range data {
		sw.ParseText(str)
		file.WriteString(str)
	}
	log.Print(sw.GenerateText("", 15))
}
func test() {
	sw := &swatter.DataStorage{}
	sw.ReadFile("mh.txt")

	log.Print(sw.GenerateText("кум", 5))
	log.Print(sw.GenerateText("рома", 10))
	log.Print(sw.GenerateText("да", 15))
	log.Print(sw.GenerateText("нет", 25))
}
