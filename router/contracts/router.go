package contracts

import (
    "gokatan/support"
    "net/http"
)

type Router interface {
    // Route registration methods
    Get(uri, controller, action string, callback *support.Callback) Route
    Post(uri, controller, action string, callback *support.Callback) Route
    Put(uri, controller, action string, callback *support.Callback) Route
    Delete(uri, controller, action string, callback *support.Callback) Route
    Match(methods []string, uri, controller, action string, callback *support.Callback) Route
    Any(uri, controller, action string, callback *support.Callback) Route
    
    // Group management
    Group(prefix string, callback func()) Router
    
    // Middleware management
    AddMiddleware(name string, middleware ...interface{}) Router
    AddMiddlewareGroups(name string, middleware ...interface{}) Router
    
    // Route finding
    FindRoute(name string) Route
    
    // HTTP Handler implementation
    ServeHTTP(w http.ResponseWriter, r *http.Request)
}
