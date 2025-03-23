package contracts

type IRouteCollection interface {
    // Route registration
    Add(route Route) Route
    
    // Route lookup
    Match(uri, method string) Route
    GetByName(name string) Route
    
    // Group management
    Group(attributes map[string]interface{}, callback func()) IRouteCollection
    
    // Getters
    GetRoutes() map[string]Route
    GetNamedRoutes() map[string]Route
}
