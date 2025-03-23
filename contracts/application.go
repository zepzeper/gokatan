package contracts

type IApplication interface {
    Bind(name string, concrete interface{})
    Singleton(name string, concrete interface{})
    Resolve(name string) (interface{}, bool)
    Boot()
    LoadEnvironment() error
    
}
