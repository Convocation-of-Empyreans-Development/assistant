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
	distances := make(map[string][]int)
	for _, system := range config.HomeSystems {
		destination := esi.SystemNameToId(config.ESIClient, system)
		distances[system] = make([]int, 2)
		distances[system][0] = esi.GetDistanceToSystem(config.ESIClient, origin, destination, esi.Shortest)
		distances[system][1] = esi.GetDistanceToSystem(config.ESIClient, origin, destination, esi.Secure)
	}
	config.AtlantisDistances = distances
}

// SendAtlantisLocationEmbed produces a message embed containing the Atlantis entrance information and
// sends it in the channel where the command used to request it was issued.
func SendAtlantisLocationEmbed(s *discordgo.Session, m *discordgo.MessageCreate, config Config) {
	embed := &discordgo.MessageEmbed{
		Type:        "",
		Title:       fmt.Sprintf("Atlantis entrance: %v", config.AtlantisEntrance),
		Description: GenerateDistanceString(config.AtlantisDistances),
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// GenerateDistanceString creates a human-readable list of the distances to each of the home systems.
func GenerateDistanceString(distances map[string][]int) string {
	builder := []string{}
	for system, distance := range distances {
		builder = append(builder, fmt.Sprintf("%s: short %d, safe %d", system, distance[0], distance[1]))
	}
	return strings.Join(builder, "\n")
}
