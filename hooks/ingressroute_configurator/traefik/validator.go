package traefik

import (
	"fmt"
	"strings"
)

//ValidateMatchRule validates and returns true if both URL and the header value matches
func ValidateMatchRule(url string , headerValue string , route Route ) bool {
	headerMatchString := fmt.Sprintf("`%s`,`%s`",HEADER_KEY,headerValue)
    return strings.Contains(route.Match,url) && strings.Contains(route.Match,headerMatchString)
}

//IsRulePresent validates if there is already an entry for routes in the Ingress route object
func IsRulePresent(url string , headerValue string , routeSpec IngressRouteSpec) bool {
	for _,route := range routeSpec.Routes {
		if ValidateMatchRule(url , headerValue, route) {
			return true
		}
	}
	return false
}




