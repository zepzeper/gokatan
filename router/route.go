package router

import (
	"net/http"
	"strings"
)

type Route struct {
    method string
    uri string
    controller string
    action string
    name string
    collection *RouteCollection
    domain string
    prefix string
    handler http.HandlerFunc
    middleware []interface{}
}

func (route *Route) Handler(handler http.HandlerFunc) *Route {
    route.handler = handler;
    return route;
}

func (route *Route) Middleware(middleware ...interface{}) *Route {
    route.middleware = append(route.middleware, middleware...)
    return route
}

// Name sets the route's name
func (route *Route) Name(name string) *Route {
    route.name = name
    
    // Update the route collection if needed
    if route.collection != nil {
        identifier := generateRouteIdentifier(*route)
        route.collection.updateNamedRoute(identifier, *route)
    }
    
    return route
}

func (route *Route) Matches(uri string, method string) bool {
    if route.method != method {
        return false
    }

    if route.uri == uri {
        return true
    }

    matches, _ := pathMatches(route.uri, uri);
    return matches;
}

func (route *Route) Prefix(prefix string) *Route {
    prefix = strings.Trim(prefix, "/")
    if prefix != "" {
        route.uri = prefix + "/" + strings.TrimPrefix(route.uri, "/")
    }
    return route
}

func (route *Route) Domain(domain string) *Route {
    route.domain = domain
    return route
}

func (route *Route) GetMethods() []string {
    return []string{route.method}
}

func (route *Route) GetUri() string {
    return route.uri
}

func (route *Route) getName() string {
    return route.name
}

func (route *Route) setCollection(collection *RouteCollection) {
    route.collection = collection
}
