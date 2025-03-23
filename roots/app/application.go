package app

import (
	"gokatan/config"
	"gokatan/roots/contracts"
	"net/http"
	netHttp "net/http"
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

func (a *Application) HandleRequest(w http.ResponseWriter, r *http.Request) {
    kernelInterface, exists := a.Resolve("http.kernel");

    if !exists {
        netHttp.Error(w, "Kernel not found", netHttp.StatusInternalServerError);
        return;
    }

    kernel, ok := kernelInterface.(contracts.IKernel)
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

func (a *Application) RunServer() error {
	portInterface, exists := a.Resolve("env.get")
	port := "8000" 
	
	if exists {
		if fn, ok := portInterface.(func(string, string) string); ok {
			port = fn("APP_PORT", "8000")
		}
	}
	
	return http.ListenAndServe(":"+port, a.router)
}
