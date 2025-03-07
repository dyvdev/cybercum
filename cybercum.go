package cybercum

import (
	"github.com/dyvdev/cybercum/tgbot"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func RunBot() {
	log.Println("starting...")
	done := make(chan bool)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM)
	bot := tgbot.NewBot()
	if bot == nil {
		log.Println("failed to start bot...")
		return
	}
	bot.Update(done)
	bot.Dumper(done)
	<-sigc
	done <- true
}
