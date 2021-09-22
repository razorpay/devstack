package traefik

import (
	"errors"
	"fmt"
)

// InsertMatchRule inserts the match rule to the ingressroute spec
func InsertMatchRule(url string , headerValue string ,serviceName string , port int32, routeSpec *IngressRouteSpec) {
	matchString := fmt.Sprintf("Host(`%s`) && Headers (`%s`,`%s`)",url,HEADER_KEY,headerValue)
	services := []Service{{LoadBalancerSpec{Name: serviceName, Port: port}}}
	route := Route{Match: matchString,Kind: KIND,Services: services }
	routeSpec.Routes = append(routeSpec.Routes, route)
}

// DeleteMatchRule deletes the match rule from the ingressroute spec
func DeleteMatchRule(url string , headerValue string , routeSpec *IngressRouteSpec) error {
	deleteIndex := -1
	for i,route := range routeSpec.Routes {
		if ValidateMatchRule(url,headerValue,route) {
			deleteIndex = i
		}
	}
	if deleteIndex == -1 {
		return errors.New("Index not found in the Match rule")
	}
	routeSpec.Routes = append(routeSpec.Routes[:deleteIndex], routeSpec.Routes[deleteIndex+1:]...)
	return nil
}
