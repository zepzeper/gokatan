package contracts

import (
    "gokatan/support"
    "net/http"
)

type IRouter interface {
    // Route registration methods
    Get(uri, controller, action string, callback *support.Callback) IRoute
    Post(uri, controller, action string, callback *support.Callback) IRoute
    Put(uri, controller, action string, callback *support.Callback) IRoute
    Delete(uri, controller, action string, callback *support.Callback) IRoute
    Match(methods []string, uri, controller, action string, callback *support.Callback) IRoute
    Any(uri, controller, action string, callback *support.Callback) IRoute
    
    // Group management
    Group(prefix string, callback func()) IRouter
    
    // Middleware management
    AddMiddleware(name string, middleware ...interface{}) IRouter
    AddMiddlewareGroups(name string, middleware ...interface{}) IRouter
    
    // Route finding
    FindRoute(name string) IRoute
    
    // HTTP Handler implementation
    ServeHTTP(w http.ResponseWriter, r *http.Request)
}
