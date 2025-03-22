package app

import (
    netHttp "net/http"
    "gokatan/config"
)

type Application struct {
    bindings   map[string]interface{} // Dict key: "string" => value: "interface aka any type"
    singletons map[string]interface{}
    booted     bool
}

func New() *Application {
    return &Application{ // &Application returns a pointer
        bindings:   make(map[string]interface{}),
        singletons: make(map[string]interface{}),
        booted:     false,
    }
}

func (a *Application) Bind(name string, concrete interface{}) {
    a.bindings[name] = concrete
}

func (a *Application) Singleton(name string, concrete interface{}) {
    a.singletons[name] = concrete
}

func (a *Application) Resolve(name string) (interface{}, bool) {
    if instance, exists := a.singletons[name]; exists {
        if fn, ok := instance.(func() interface{}); ok {
            return fn(), true
        }

        return instance, true
    }

    if bindings, exists := a.bindings[name]; exists {
        if fn, ok := bindings.(func() interface{}); ok {
            return fn(), true
        }
        return bindings, true
    }

    return nil, false
}

func (a *Application) Boot() {
    if !a.booted {
        a.booted = true
    }
}

func (a *Application) LoadEnvironment() error {
    loader := config.NewEnvLoader(".env");
    err := loader.Load();

    a.Singleton("env.get", func() interface{} {
        return config.Get
    });

    return err;
}

func handleRequest(a *Application, w netHttp.ResponseWriter, r *netHttp.Request) {

    kernelInterface, exists := a.Resolve("http.kernel");

    if !exists {
        netHttp.Error(w, "Kernel not found", netHttp.StatusInternalServerError);
        return;
    }

    kernel, ok := kernelInterface.(*http.Kernel)
    if !ok {
        netHttp.Error(w, "Invalid kernel type", netHttp.StatusInternalServerError)
        return
    }

    err := kernel.Handle(r)

    if err != nil {
        netHttp.Error(w, "Internal Server Error", netHttp.StatusInternalServerError)
        return
    }
}
