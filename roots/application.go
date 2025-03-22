package app

import (
	"fmt"
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
    err := a.loadEnvironment()
    if err != nil {
        fmt.Println("Error loading environment:", err)
    }

    // Retrieve the env.get function
    envGetFunc, exists := a.Resolve("env.get")

    if !exists {
        fmt.Println("Environment getter not found")
        return
    }

    if fn, ok := envGetFunc.(func(string, string) string); ok {
        // Try getting an environment variable with a default value
        testEnvVar := fn("APP_NAME", "default_value")
        fmt.Println("Test Env Var Value:", testEnvVar)
    } else {
        fmt.Printf("Unexpected type: %T\n", envGetFunc)
    }
}

func (a *Application) loadEnvironment() error {
    loader := config.NewEnvLoader(".env");
    err := loader.Load();

    a.Singleton("env.get", func() interface{} {
        return config.Get
    });

    return err;
}
