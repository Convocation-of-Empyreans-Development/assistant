package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	bot "github.com/lichgrave/MALRO_incursion_bot/discord"
	"os"
	"os/signal"
	"syscall"
)

//creates a websocket to connect to the bot
func main() {
	config := bot.ReadConfig("./config.json")
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	fmt.Println("Connection successful")

	dg.AddHandler(bot.HandleMessageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
