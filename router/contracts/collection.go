package contracts

type RouteCollection interface {
    // Route registration
    add(route Route) Route
    
    // Route lookup
    Match(uri, method string) Route
    GetByName(name string) Route
    
    // Group management
    Group(attributes map[string]interface{}, callback func()) RouteCollection
    
    // Getters
    GetRoutes() map[string]Route
    GetNamedRoutes() map[string]Route
}
