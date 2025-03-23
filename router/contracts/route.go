package contracts

import "net/http"

type IRoute interface {
    // Path and method matching
    Matches(uri, method string) bool
    GetUri() string
    GetMethods() []string
    
    // Configuration
    Name(name string) IRoute
    Prefix(prefix string) IRoute
    Domain(domain string) IRoute
    Handler(handler http.HandlerFunc) IRoute
    Middleware(middleware ...interface{}) IRoute
    
    // Getters
    GetName() string
    GetHandler() http.HandlerFunc
    GetMiddleware() []interface{}
}
