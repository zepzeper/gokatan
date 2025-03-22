package app;

type Application struct {
    bindings map[string]interface{} // Dict key: "string" => value: "interface aka any type"
    singletons map[string]interface{}
    booted bool
}

func New() *Application {
    return &Application { // &Application returns a pointer
        bindings: make(map[string]interface{}),
        singletons: make(map[string]interface{}),
        booted: false,
    }
}

func (a *Application) Bind(name string, concrete interface{}) {
    a.bindings[name] = concrete
}

func (a *Application) Singleton(name string, concrete interface{}) {
    a.bindings[name] = concrete
}

func (a *Application) Resolve(name string) (interface{}, bool) {
    if instance, exists := a.singletons[name]; exists {
        if fn, ok := instance.(func() interface{}); ok {
            return fn(), true;
        }

        return instance, true;
    }

    if bindings, exists := a.bindings[name]; exists {
        if fn, ok := bindings.(func() interface{}); ok {
            return fn(), true;
        }
        return bindings, true;
    }

    return nil, false;
}

func (a *Application) Boot() {
    if !a.booted {
        a.booted = true;
    }

    return;
}
