package esi

import (
	"context"
	"github.com/antihax/goesi"
	"github.com/antihax/goesi/esi"
	"net/http"

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

// CheckESIResponse checks whether an error was returned by the ESI API, and panics if this is the case.
func CheckESIResponse(err error, response *http.Response) {
	if err != nil || response.StatusCode != http.StatusOK {
		panic(err)
	}
}

// CreateESIClient creates a new ESI client.
func CreateESIClient(userAgent string) *goesi.APIClient {
	client := goesi.NewAPIClient(&http.Client{}, userAgent)
	return client
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

// IdToName converts an int32 ID to the corresponding name using the ESI API.
func IdToName(client *goesi.APIClient, id int32) string {
	names, response, err := client.ESI.UniverseApi.PostUniverseNames(context.TODO(), []int32{id}, nil)
	CheckESIResponse(err, response)
	return names[0].Name
}

// IdsToNames converts a list of int32 IDs to their corresponding names using the ESI API.
func IdsToNames(client *goesi.APIClient, ids []int32) (names []string) {
	apiNames, response, err := client.ESI.UniverseApi.PostUniverseNames(context.TODO(), ids, nil)
	CheckESIResponse(err, response)
	for _, item := range apiNames {
		names = append(names, item.Name)
	}
	return names
}

// GetSecurityStatus gets the security status of a given system ID from the ESI API.
func GetSecurityStatus(client *goesi.APIClient, systemID int32) float32 {
	system, response, err := client.ESI.UniverseApi.GetUniverseSystemsSystemId(context.TODO(), systemID, nil)
	CheckESIResponse(err, response)
	return system.SecurityStatus
}

// GetIncursions fetches the latest incursion data from ESI, processes it and returns it to the caller.
func GetIncursions(client *goesi.APIClient) []IncursionData {
	rawData := GetIncursionData(client)
	return ProcessIncursionData(client, rawData)
}
