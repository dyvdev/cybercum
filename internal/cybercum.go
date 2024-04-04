package internal

import (
	"github.com/dyvdev/cybercum/internal/config"
	"github.com/dyvdev/cybercum/internal/tgbot"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func ReadBot(cfgFile string) {
	log.Println("starting...")
	bot, err := tgbot.NewBot(nil)
	if err != nil {
		log.Fatalf("create bot: %s", err.Error())
	}

	log.Println("reading...")
	//bot.Swatter.ReadFile("mh.txt")
	log.Println("saving...")
	bot.SaveDump()
	log.Println("done...")
}

func CleanBot() {
	bot, err := tgbot.NewBot(nil)
	if err != nil {
		log.Fatalf("create bot: %s", err.Error())
	}

	bot.Clean()
	bot.SaveDump()
}

func RunBot(cfg *config.Config) {
	log.Println("starting...")
	done := make(chan bool)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM)
	bot, err := tgbot.NewBot(cfg)
	if err != nil {
		log.Fatalf("create bot: %s", err.Error())
	}

	bot.Update(done)
	bot.Dumper(done)
	<-sigc
	done <- true
}
