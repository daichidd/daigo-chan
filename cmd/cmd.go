package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/daichidd/daigo-chan/handler"
)

func Run(token string) error {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}
	defer dg.Close()

	dg.AddHandler(handler.MessageCreate)

	if err := dg.Open(); err != nil {
		return err
	}
	log.Println("daigo-chan is running")

	// press ctrl-c to exit
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	return nil
}
