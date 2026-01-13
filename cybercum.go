package cybercum

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dyvdev/cybercum/tgbot"
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
	wg := &sync.WaitGroup{}
	wg.Add(2)
	bot.Update(wg, done)
	bot.Dumper(wg, done)
	select {
	case <-sigc:
		close(done)
	}
	wg.Wait()
}
