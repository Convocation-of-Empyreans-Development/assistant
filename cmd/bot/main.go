package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/antihax/goesi"
	"github.com/bwmarrin/discordgo"

	bot "github.com/Convocation-of-Empyreans-Development/MALRO_incursion_bot/discord"
	"github.com/Convocation-of-Empyreans-Development/MALRO_incursion_bot/esi"
)

var configFilename = flag.String("config", "config.json", "path to the bot configuration file")
var memcachedAddress = flag.String("memcached-address", "",
	"address (host:port) for memcached instance used by ESI API client")
var userAgent = flag.String("user-agent", "MALRO Incursions Monitor",
	"User agent identifying the ESI API client to CCP")

//creates a websocket to connect to the bot
func main() {
	var client *goesi.APIClient
	flag.Parse()
	if memcachedAddress != nil {
		// Create ESI client with caching
		client = esi.CreateCachingESIClient(*userAgent, *memcachedAddress)
	} else {
		client = esi.CreateESIClient(*userAgent)
	}

	config := bot.ReadConfig(*configFilename)
	config.ESIClient = client

	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	fmt.Println("Connection successful")

	// Add a handler for received messages.
	dg.AddHandler(bot.HandleMessageCreate(config))

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
