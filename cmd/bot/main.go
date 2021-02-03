package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	bot "github.com/lichgrave/MALRO_incursion_bot/discord"
)

var configFilename = flag.String("config", "config.json", "path to the bot configuration file")

//creates a websocket to connect to the bot
func main() {
	flag.Parse()
	config := bot.ReadConfig(*configFilename)
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	fmt.Println("Connection successful")

	// Add a handler for received messages.
	dg.AddHandler(bot.HandleMessageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other process termination signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
