package discord

import (
	"github.com/antihax/goesi"

	"github.com/Convocation-of-Empyreans-Development/MALRO_incursion_bot/esi"
)

// SetAtlantisEntranceLocation sets the current entrance to the Atlantis wormhole.
// This function also calculates the distances to each of the configured "home systems" for later retrieval.
func SetAtlantisEntranceLocation(client *goesi.APIClient, location string, config Config) {
	origin := esi.SystemNameToId(client, location)
	config.AtlantisEntrance = location
	distances := make(map[string]int)
	for _, system := range config.HomeSystems {
		destination := esi.SystemNameToId(client, location)
		distances[system] = esi.GetDistanceToSystem(client, origin, destination)
	}
	config.AtlantisDistances = distances
}
