package esi

import (
	"context"
	"log"
	"net/http"
	"runtime"

	"github.com/antihax/goesi"
)

// CheckESIResponse checks whether an error was returned by the ESI API, and panics if this is the case.
func CheckESIResponse(err error, response *http.Response) {
	if err != nil || response.StatusCode != http.StatusOK {
		caller, file, line, ok := runtime.Caller(1)
		if ok {
			log.Printf("%v:%v (%v) | %v -> %v", file, line, caller, response.Status, err.Error())
		}
	}
}

// CreateESIClient creates a new ESI client.
func CreateESIClient(userAgent string) *goesi.APIClient {
	client := goesi.NewAPIClient(&http.Client{}, userAgent)
	return client
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

// SystemNameToId converts the name of a solar system to its system ID using the ESI API.
func SystemNameToId(client *goesi.APIClient, name string) (id int32) {
	ids, response, err := client.ESI.UniverseApi.PostUniverseIds(context.TODO(), []string{name}, nil)
	CheckESIResponse(err, response)
	return ids.Systems[0].Id
}

// GetSecurityStatus gets the security status of a given system ID from the ESI API.
func GetSecurityStatus(client *goesi.APIClient, systemID int32) float32 {
	system, response, err := client.ESI.UniverseApi.GetUniverseSystemsSystemId(context.TODO(), systemID, nil)
	CheckESIResponse(err, response)
	return system.SecurityStatus
}
