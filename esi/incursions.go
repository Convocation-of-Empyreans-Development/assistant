package esi

import (
	"context"
	"fmt"
	"github.com/antihax/goesi"
	"github.com/antihax/goesi/esi"
	"net/http"
)

// IncursionData holds processed data on the state of a currently active incursion.
type IncursionData struct {
	Constellation        string
	Faction              string
	HasBoss              bool
	InfestedSolarSystems []string
	Influence            float32
	StagingSolarSystem   string
	State                string
	Type                 string
	SecurityStatus       float32
}

// Checks whether an error was returned by the ESI API, and panics if this is the case.
func CheckESIResponse(err error, response *http.Response) {
	if err != nil || response.StatusCode != http.StatusOK {
		panic(err)
	}
}

// Create a new ESI client.
func CreateESIClient() *goesi.APIClient {
	client := goesi.NewAPIClient(&http.Client{}, "MALRO Incursions Monitor")
	return client
}

// Get raw incursion data from the ESI API.
func GetIncursionData(client *goesi.APIClient) []esi.GetIncursions200Ok {
	incursions, response, err := client.ESI.IncursionsApi.GetIncursions(context.TODO(), nil)
	CheckESIResponse(err, response)
	return incursions
}

// Process the raw incursion data returned by the API into a human-readable form.
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

// Converts an int32 ID to the corresponding name using the ESI API.
func IdToName(client *goesi.APIClient, id int32) string {
	names, response, err := client.ESI.UniverseApi.PostUniverseNames(context.TODO(), []int32{id}, nil)
	CheckESIResponse(err, response)
	return names[0].Name
}

// Converts a list of int32 IDs to their corresponding names using the ESI API.
func IdsToNames(client *goesi.APIClient, ids []int32) (names []string) {
	apiNames, response, err := client.ESI.UniverseApi.PostUniverseNames(context.TODO(), ids, nil)
	CheckESIResponse(err, response)
	for _, item := range apiNames {
		names = append(names, item.Name)
	}
	return names
}

func GetSecurityStatus(client *goesi.APIClient, systemID int32) float32 {
	system, response, err := client.ESI.UniverseApi.GetUniverseSystemsSystemId(context.TODO(), systemID, nil)
	CheckESIResponse(err, response)
	return system.SecurityStatus
}

func GetIncursions() []IncursionData {
	client := CreateESIClient()
	rawData := GetIncursionData(client)
	return ProcessIncursionData(client, rawData)
}

// main() is only used for testing the incursion data fetch & process.
func main() {
	client := CreateESIClient()
	rawData := GetIncursionData(client)
	fmt.Printf("%+v\n", ProcessIncursionData(client, rawData))
}
