package esi

import (
	"context"
	"net/http"

	"github.com/antihax/goesi"
	"github.com/antihax/goesi/esi"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gregjones/httpcache"
	httpmemcache "github.com/gregjones/httpcache/memcache"
)

// IncursionData holds processed data on the state of a currently active incursion.
type IncursionData struct {
	// Constellation represents the constellation under incursion.
	Constellation string
	// Faction represents the faction responsible for the incursion.
	Faction string
	// HasBoss represents whether the incursion boss is present in the HQ system.
	HasBoss bool
	// InfestedSolarSystems represents the list of systems currently under incursion.
	InfestedSolarSystems []string
	// Influence represents the current level of influence on a scale from 0 to 1.
	Influence float32
	// StagingSolarSystem represents the uninvaded system acting as the staging point for the incursion.
	StagingSolarSystem string
	// State represents the current state of the incursion.
	State string
	// Type represents the type of incursion created by the server.
	Type string
	// SecurityStatus represents the security status of the staging system.
	SecurityStatus float32
}

// CreateCachingESIClient creates a new ESI client, backed by a connection to a memcached server.
// This should reduce the number of requests made to the ESI API.
func CreateCachingESIClient(userAgent string, memcachedAddress string) *goesi.APIClient {
	// Connect to the memcached server
	cache := memcache.New(memcachedAddress)

	// Create a memcached http client for the ESI APIs.
	transport := httpcache.NewTransport(httpmemcache.NewWithClient(cache))
	transport.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
	client := &http.Client{Transport: transport}

	// Get our API Client.
	return goesi.NewAPIClient(client, userAgent)
}

// GetIncursionData gets raw incursion data from the ESI API.
func GetIncursionData(client *goesi.APIClient) []esi.GetIncursions200Ok {
	incursions, response, err := client.ESI.IncursionsApi.GetIncursions(context.TODO(), nil)
	CheckESIResponse(err, response)
	return incursions
}

// ProcessIncursionData processes the raw incursion data returned by the API into a human-readable form.
// IDs for systems, constellations and factions are converted into their equivalent names.
func ProcessIncursionData(client *goesi.APIClient, data []esi.GetIncursions200Ok) (incursions []IncursionData) {
	for _, incursion := range data {
		processedIncursion := IncursionData{
			Constellation:        IdToName(client, incursion.ConstellationId),
			Faction:              IdToName(client, incursion.FactionId),
			HasBoss:              incursion.HasBoss,
			InfestedSolarSystems: IdsToNames(client, incursion.InfestedSolarSystems),
			Influence:            incursion.Influence,
			StagingSolarSystem:   IdToName(client, incursion.StagingSolarSystemId),
			State:                incursion.State,
			Type:                 incursion.Type_,
			SecurityStatus:       GetSecurityStatus(client, incursion.StagingSolarSystemId),
		}
		incursions = append(incursions, processedIncursion)
	}
	return incursions
}

// GetIncursions fetches the latest incursion data from ESI, processes it and returns it to the caller.
func GetIncursions(client *goesi.APIClient) []IncursionData {
	rawData := GetIncursionData(client)
	return ProcessIncursionData(client, rawData)
}
