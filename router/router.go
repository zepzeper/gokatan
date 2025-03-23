package router

import (
    "fmt"
    "gokatan/support"
    "net/http"
    "slices"
    "strings"
)

var VALID_METHODS = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}

type Router struct {
    middleware map[string]interface{}
    middlewareGroups map[string][]interface{}
    middlewareAlias map[string]interface{}
    routes *RouteCollection
    
}

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

func NewRouter() *Router {
    return &Router{
        middleware:       make(map[string]interface{}),
        middlewareGroups: make(map[string][]interface{}),
        middlewareAlias:  make(map[string]interface{}),
        routes:           NewRouteCollection(),
    }
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    route := r.routes.Match(req.URL.Path, req.Method)
    if route == nil {
        http.NotFound(w, req)
        return
    }
    
    if route.handler != nil {
        route.handler(w, req)
    } else {
        // Handle controller/action routing
        fmt.Fprintf(w, "Route matched: %s %s -> %s.%s", 
            route.method, route.uri, route.controller, route.action)
    }
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

func (r *Router) Get(uri string, controller string, action string, callback *support.Callback) *Route {
    resolve(callback);
    return r.add([]string{"GET"}, uri, controller, action, nil)
}

func (r *Router) Post(uri string, controller string, action string, callback *support.Callback) *Route {
    resolve(callback);
    return r.add([]string{"POST"}, uri, controller, action, nil)
}

func (r *Router) Put(uri string, controller string, action string, callback *support.Callback) *Route {
    resolve(callback);
    return r.add([]string{"PUT"}, uri, controller, action, nil)
}

func (r *Router) Delete(uri string, controller string, action string, callback *support.Callback) *Route {
    resolve(callback);
    return r.add([]string{"DELETE"}, uri, controller, action, nil)
}

func (r *Router) Match(methods []string, uri string, controller string, action string, callback *support.Callback) *Route {
    resolve(callback);
    return r.add(methods, uri, controller, action, nil)
}

func (r *Router) Any(uri string, controller string, action string, callback *support.Callback) *Route {
    resolve(callback);
    return r.add(VALID_METHODS, uri, controller, action, nil)
}

func (route *Route) Handler(handler http.HandlerFunc) *Route {
    route.handler = handler;
    return route;
}

func (route *Route) Middleware(middleware ...interface{}) *Route {
    route.middleware = append(route.middleware, middleware...)
    return route
}

func (r *Router) AddMiddleware(name string, middleware ...interface{}) *Router {
    r.middleware[name] = middleware;
    return r
}

func (r *Router) AddMiddlewareGroups(name string, middleware ...interface{}) *Router {
    r.middlewareGroups[name] = middleware;
    return r
}

func (r *Router) AddMiddlewareAlias(name string, middleware ...interface{}) *Router {
    r.middlewareAlias[name] = middleware;
    return r
}

func (r *Router) add(methods []string, uri string, controller string, action string, callback *support.Callback) *Route {
    // Create a new route for each method
    var lastRoute *Route
    
    for _, method := range methods {
        if !isValidMethod(method) {
            fmt.Println("Invalid HTTP method:", method)
            continue
        }
        
        newRoute := Route{
            method:     strings.ToUpper(method),
            uri:        uri,
            controller: controller,
            action:     action,
            middleware: []interface{}{},
        }

        // Add the route to the collection
        r.routes.add(newRoute)
    }
    
    if callback != nil {
        (*callback)()
    }
    
    return lastRoute
}


func resolve(callback *support.Callback) {
    if callback != nil {
        (*callback)();
    }
}

func (route *Route) getName() string {
    return route.name
}

func (route *Route) setCollection(collection *RouteCollection) {
    route.collection = collection
}

// pathMatches checks if a request path matches a route pattern
func pathMatches(pattern, path string) (bool, map[string]string) {
    // Split paths into segments
    patternSegments := strings.Split(strings.Trim(pattern, "/"), "/")
    pathSegments := strings.Split(strings.Trim(path, "/"), "/")
    
    // Quick length check
    if len(patternSegments) != len(pathSegments) {
        return false, nil
    }
    
    // Check each segment
    params := make(map[string]string)
    for i, segment := range patternSegments {
        // Parameter segment (e.g., ":id")
        if strings.HasPrefix(segment, ":") {
            paramName := segment[1:]
            params[paramName] = pathSegments[i]
            continue
        }
        
        // Static segment - must match exactly
        if segment != pathSegments[i] {
            return false, nil
        }
    }
    
    return true, params
}

func isValidMethod(method string) bool {
    return slices.Contains(VALID_METHODS, method);

    // Original, above suggested by go lsp
    // for _, validMethod := range VALID_METHODS {
    //     if method == validMethod {
    //         return true;
    //     }
    // }
    //
    // return false;
}
