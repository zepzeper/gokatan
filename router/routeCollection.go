package router

import (
	"crypto/md5"
	"fmt"
	"gokatan/router/contracts"
	"strings"
)

var _ contracts.IRouteCollection = (*RouteCollection)(nil)

type RouteCollection struct {
    routes map[string]Route
    namedRoutes map[string]Route
    groups []map[string]interface{}
    currentGroup map[string]interface{}
}

func NewRouteCollection() *RouteCollection {
    return &RouteCollection{
        routes: make(map[string]Route),
        namedRoutes: make(map[string]Route),
        groups: []map[string]interface{}{},
    }
}

func (rc *RouteCollection) Match(uri string, method string) contracts.IRoute {
    for _, route := range rc.routes {

        if route.Matches(uri, method) {
            routeCopy := route
            return &routeCopy
        }
    }
    return nil
}

func (rc *RouteCollection) GetByName(name string) contracts.IRoute {
    if route, exists := rc.namedRoutes[name]; exists {
        routeCopy := route
        return &routeCopy
    }
    return nil
}

func (rc *RouteCollection) Group(attributes map[string]interface{}, callback func()) contracts.IRouteCollection {
    previousGroup := rc.currentGroup

    if previousGroup != nil {
        rc.currentGroup = rc.mergeGroups(previousGroup, attributes)
    } else {
        rc.currentGroup = attributes
    }

    rc.groups = append(rc.groups, rc.currentGroup)

    callback()

    rc.groups = rc.groups[:len(rc.groups) - 1]
    rc.currentGroup = previousGroup

    return rc;
}

func (rc *RouteCollection) Add(route contracts.IRoute) contracts.IRoute {
    routeImpl, ok := route.(*Route)
    if !ok {
        // Create a new Route with the information from the interface
        newRoute := Route{
            method: route.GetMethods()[0],
            uri: route.GetUri(),
            name: route.GetName(),
            handler: route.GetHandler(),
            middleware: route.GetMiddleware(),
        }
        routeImpl = &newRoute
    }
    
    routeImpl.setCollection(rc)
    
    identifier := generateRouteIdentifier(*routeImpl)

    if rc.currentGroup != nil {
        rc.applyGroupSettingsToRoute(routeImpl)
    }

    rc.routes[identifier] = *routeImpl
    
    name := routeImpl.GetName()

    if name != "" {
        rc.addNamedRoute(*routeImpl)
    }

    storedRoute := rc.routes[identifier]
    return &storedRoute
}

func (rc *RouteCollection) GetRoutes() map[string]contracts.IRoute {
    result := make(map[string]contracts.IRoute)
    for k, v := range rc.routes {
        routeCopy := v
        result[k] = &routeCopy
    }
    return result
}

func (rc *RouteCollection) GetNamedRoutes() map[string]contracts.IRoute {
    result := make(map[string]contracts.IRoute)
    for k, v := range rc.namedRoutes {
        routeCopy := v
        result[k] = &routeCopy
    }
    return result
}

func (rc *RouteCollection) mergeGroups(previous map[string]interface{}, new map[string]interface{}) map[string]interface{} {
    merged := make(map[string]interface{})

    for k,v := range previous {
        merged[k] = v
    }

    for k, v := range new {
        if k == "prefix" {
            // Special handling for prefix
            prevPrefix, prevOk := previous["prefix"].(string)
            newPrefix, newOk := v.(string)

            if prevOk && newOk {
                // Join prefixes with a slash
                merged[k] = fmt.Sprintf("%s/%s", strings.Trim(prevPrefix, "/"), strings.Trim(newPrefix, "/"))
            } else if newOk {
                merged[k] = newPrefix
            }
        } else if k == "middleware" {
            // Special handling for middleware
            prevMiddleware, prevOk := previous["middleware"].([]interface{})
            newMiddleware, newOk := v.([]interface{})

            if prevOk && newOk {
                // Combine middleware arrays
                merged[k] = append(prevMiddleware, newMiddleware...)
            } else if newOk {
                merged[k] = newMiddleware
            }
        } else {
            // Default case: just override
            merged[k] = v
        }
    }

    return merged
}
 
func (rc *RouteCollection) addNamedRoute(route Route) {
    name := route.GetName()

    if _, exists := rc.namedRoutes[name]; exists {
        panic(fmt.Sprintf("route name '%s' has already been taken", name))
    }

    rc.namedRoutes[name] = route
}

func (rc *RouteCollection) updateNamedRoute(identifier string, route Route) {
    name := route.GetName()

    if name != "" {
        for k := range rc.namedRoutes {
            oldRoute := rc.namedRoutes[k]
            if generateRouteIdentifier(oldRoute) == identifier {
                rc.namedRoutes[name] = route
            }
        }
    }
}

func (rc *RouteCollection) applyGroupSettingsToRoute(route *Route) {
    if middleware, ok := rc.currentGroup["middleware"]; ok {
        if middlewareArray, ok := middleware.([]interface{}); ok {
            route.Middleware(middlewareArray...)
        }
    }

    if prefix, ok := rc.currentGroup["prefix"].(string); ok {
        route.Prefix(prefix)
    }

    if namePrefix, ok := rc.currentGroup["as"].(string); ok {
        currentName := route.GetName()
        if currentName != "" {
            route.Name(namePrefix + currentName)
        }
    }

    if domain, ok := rc.currentGroup["domain"].(string); ok {
        route.Domain(domain)
    }
}

func generateRouteIdentifier(route Route) string {
    hash := md5.Sum([]byte(route.method + route.uri))
    return fmt.Sprintf("%x", hash) // Convert to hex string
}
