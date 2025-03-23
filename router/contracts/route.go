package contracts

import "net/http"

type Route interface {
    // Path and method matching
    Matches(uri, method string) bool
    GetUri() string
    GetMethods() []string
    
    // Configuration
    Name(name string) Route
    Prefix(prefix string) Route
    Domain(domain string) Route
    Handler(handler http.HandlerFunc) Route
    Middleware(middleware ...interface{}) Route
    
    // Getters
    GetName() string
    GetHandler() http.HandlerFunc
    GetMiddleware() []interface{}
}
