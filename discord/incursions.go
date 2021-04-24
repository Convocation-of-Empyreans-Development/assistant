package discord

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antihax/goesi"
	"github.com/bwmarrin/discordgo"

	"github.com/Convocation-of-Empyreans-Development/MALRO_incursion_bot/esi"
)

// SendSelectedIncursionDataEmbed searches for an incursion in the selected constellation.
// If an incursion is present and active in the selected constellation, the bot will send an embed
// containing the relevant information. Otherwise, the bot will output an error message.
func SendSelectedIncursionDataEmbed(s *discordgo.Session, m *discordgo.MessageCreate, client *goesi.APIClient) {
	// Split the command from the first and only argument (i.e. the constellation)
	command := strings.SplitN(m.Content, " ", 2)
	if len(command) != 2 {
		return
	}
	incursions := esi.GetIncursions(client)
	found := false
	for _, incursion := range incursions {
		// Perform a case-insensitive equality check for the selected constellation
		if strings.EqualFold(incursion.Constellation, command[1]) {
			found = true
			embed := CreateIncursionEmbed(incursion)
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			break
		}
	}
	if !found {
		s.ChannelMessageSend(m.ChannelID, "No incursion found in selected location.")
	}
}

// MessageInApprovedChannels checks whether a received message came from one of the specified channels.
// We use a naive linear search, O(n), since we know the list of approved channels will be very small.
// If there are no approved channels in the list, we assume that the command can be used everywhere,
// and thus return true.
func MessageInApprovedChannels(channels []string, id string) bool {
	if len(channels) == 0 {
		return true
	}
	for _, channel := range channels {
		if channel == id {
			return true
		}
	}
	return false
}

// PickColorBySecurityStatus chooses a colour for the Discord message embed based on the
// incursion's system security status
func PickColorBySecurityStatus(securitystatus float32) int {
	var color string
	if securitystatus > 0.45 {
		color = "04ff00" // high-security: green
	} else if securitystatus < 0.45 && securitystatus > 0 {
		color = "ff8400" // low-security: orange
	} else {
		color = "ff0000" // null-security: red
	}

	// convert hex #RRGGBB to int (required by discordgo)
	color64, _ := strconv.ParseInt(color, 16, 64)
	return int(color64)
}

// SendIncursionDataEmbed fetches the latest Incursion data from the ESI API,
// and converts it into some easy-to-read embedded messages sent as a reply
// in the requested channel.
func SendIncursionDataEmbed(s *discordgo.Session, m *discordgo.MessageCreate, client *goesi.APIClient) {
	incursions := esi.GetIncursions(client)
	for _, incursion := range incursions {
		embed := CreateIncursionEmbed(incursion)
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
		if err != nil {
			panic(err)
		}
	}
}

// CreateIncursionEmbed takes processed incursion data from the ESI API and creates
// a Discord embed with the relevant information.
func CreateIncursionEmbed(incursion esi.IncursionData) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: PickColorBySecurityStatus(incursion.SecurityStatus),
		Title: fmt.Sprintf("Incursion in %v", incursion.Constellation),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Staging system",
				Value:  incursion.StagingSolarSystem,
				Inline: true,
			},
			{
				Name:   "Influence",
				Value:  fmt.Sprintf("%.1f%%", incursion.Influence*100),
				Inline: true,
			},
			{
				Name:  "Infested systems",
				Value: strings.Join(incursion.InfestedSolarSystems, ", "),
			},
		},
	}
}
