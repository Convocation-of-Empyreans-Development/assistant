package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Convocation-of-Empyreans-Development/MALRO_incursion_bot/esi"
)

// SetAtlantisEntranceLocation sets the current entrance to the Atlantis wormhole.
// This function also calculates the distances to each of the configured "home systems" for later retrieval.
func SetAtlantisEntranceLocation(location string, config *Config) {
	origin := esi.SystemNameToId(config.ESIClient, location)
	config.AtlantisEntrance = location
	distances := make(map[string]int)
	for _, system := range config.HomeSystems {
		destination := esi.SystemNameToId(config.ESIClient, location)
		distances[system] = esi.GetDistanceToSystem(config.ESIClient, origin, destination)
	}
	config.AtlantisDistances = distances
}

// SendAtlantisLocationEmbed produces a message embed containing the Atlantis entrance information and
// sends it in the channel where the command used to request it was issued.
func SendAtlantisLocationEmbed(s *discordgo.Session, m *discordgo.MessageCreate, config Config) {
	embed := &discordgo.MessageEmbed{
		Type:  "",
		Title: "Atlantis entrance",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Entrance system",
				Value:  config.AtlantisEntrance,
				Inline: true,
			},
			{
				Name:   "Distance to home systems",
				Value:  GenerateDistanceString(config.AtlantisDistances),
				Inline: false,
			},
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// GenerateDistanceString creates a human-readable list of the distances to each of the home systems.
func GenerateDistanceString(distances map[string]int) string {
	builder := []string{}
	for system, distance := range distances {
		builder = append(builder, fmt.Sprintf("%s: %d", system, distance))
	}
	return strings.Join(builder, "\n")
}
