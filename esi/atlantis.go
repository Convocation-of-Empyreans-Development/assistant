package esi

import (
	"context"

	"github.com/antihax/goesi"
)

// GetDistanceToSystem returns the number of jumps in the shortest route between two systems.
// This function only considers standard, developer-built stargate routes, via the ESI API.
func GetDistanceToSystem(client *goesi.APIClient, originId, destinationId int32) int {
	route, response, err := client.ESI.RoutesApi.GetRouteOriginDestination(
		context.TODO(), destinationId, originId, nil,
	)
	CheckESIResponse(err, response)
	return len(route)
}
