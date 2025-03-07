package main

import (
	cum "github.com/dyvdev/cybercum"
	"github.com/dyvdev/cybercum/swatter"
	"github.com/dyvdev/cybercum/utils"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	cum.RunBot()
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
