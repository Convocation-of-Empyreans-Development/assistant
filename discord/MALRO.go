package discord

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/antihax/goesi"
	"github.com/bwmarrin/discordgo"
)

// Config represents the Discord bot configuration.
type Config struct {
	// Token holds the authentication token granted by Discord.
	Token string `json:"token"`
	// ApprovedChannels holds a list of approved channel IDs for which commands can be used.
	ApprovedChannels []string `json:"approved_channels"`
	// HomeSystems holds a list of system IDs pointing to important systems for the alliance.
	HomeSystems []string `json:"home_systems"`
	// AtlantisLocation holds the system ID for the location of the current Atlantis entrance.
	AtlantisEntrance string
	// AtlantisDistances holds the distances to each of the home systems from the Atlantis entrance.
	AtlantisDistances map[string]int
	// ESIClient holds the ESI API client used to make requests.
	ESIClient *goesi.APIClient
}

// ReadConfig reads a JSON file from disk containing the bot configuration
// and attempts to parse it, returning a Config struct if this is the case.
func ReadConfig(filename string) *Config {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var config = &Config{}

	if err = json.Unmarshal(dat, config); err != nil {
		panic(err)
	}
	return config
}

// HandleMessageCreate returns a function used to handle the receipt of new messages and dispatch commands.
// The closure here allows the handler to access the bot's configuration from within its scope.
// We use this method because the bot requires the handler function to have the specified signature.
func HandleMessageCreate(config *Config) func(*discordgo.Session, *discordgo.MessageCreate) {
	// We can access the bot configuration from within this function.
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages sent by the bot
		if m.Author.ID == s.State.User.ID || !MessageInApprovedChannels(config.ApprovedChannels, m.ChannelID) {
			return
		}
		// Message handling: implement commands.
		if m.Content == "!incursions" {
			// !incursions - send embeds containing current incursion data.
			SendIncursionDataEmbed(s, m, config.ESIClient)
		} else if strings.Contains(m.Content, "!info") {
			// !info <constellation> - send embed containing data for incursion in constellation if active.
			SendSelectedIncursionDataEmbed(s, m, config.ESIClient)
		}
	}
}
