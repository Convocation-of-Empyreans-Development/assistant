package discord //switch back to discord after testing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/lichgrave/MALRO_incursion_bot/esi"
)

type Config struct {
	Token string `json:"token"`
}

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

func HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// Message handling: implement commands
	if m.Content == "!incursions" {
		SendIncursionDataEmbed(s, m)
	}
}

func PickColorBySecurityStatus(securitystatus float32) int {
	var color string
	if securitystatus > 0.5 {
		color = "04ff00"
	} else if securitystatus < 0.5 && securitystatus > 0 {
		color = "ff8400"
	} else {
		color = "ff0000"
	}

	color64, _ := strconv.ParseInt(color, 16, 64)
	return int(color64)
}

// SendIncursionDataEmbed fetches the latest Incursion data from the ESI API,
// and converts it into some easy-to-read embedded messages sent as a reply
// in the requested channel.
func SendIncursionDataEmbed(s *discordgo.Session, m *discordgo.MessageCreate) {
	incursions := esi.GetIncursions()
	for _, incursion := range incursions {
		embed := &discordgo.MessageEmbed{
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
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
		if err != nil {
			panic(err)
		}
	}
}
