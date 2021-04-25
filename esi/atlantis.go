package esi

import (
	"context"

	"github.com/antihax/goesi"
	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
)

const (
	Shortest = "shortest"
	Secure   = "secure"
	Insecure = "insecure"
)

// GetDistanceToSystem returns the number of jumps in the shortest route between two systems.
// This function only considers standard, developer-built stargate routes, via the ESI API.
func GetDistanceToSystem(client *goesi.APIClient, originId, destinationId int32, routeType string) int {
	route, response, err := client.ESI.RoutesApi.GetRouteOriginDestination(
		context.TODO(), destinationId, originId,
		&esi.GetRouteOriginDestinationOpts{
			Flag: optional.NewString(routeType),
		},
	)
	CheckESIResponse(err, response)
	return len(route)
}
