package bootstrap

import (
	"fmt"
	"gokatan/roots"
)

func LoadEnvironmentVariables(app *app.Application) error {
    err := app.LoadEnvironment()
    if err != nil {
        fmt.Println("Error loading environment:", err)
    }

    // Retrieve the env.get function
    envGetFunc, exists := app.Resolve("env.get")

    if !exists {
        fmt.Println("Environment getter not found")
    }

    if fn, ok := envGetFunc.(func(string, string) string); ok {
        // Try getting an environment variable with a default value
        testEnvVar := fn("APP_NAME", "default_value")
        fmt.Println("Test Env Var Value:", testEnvVar)
    } else {
        fmt.Printf("Unexpected type: %T\n", envGetFunc)
    }
    return nil
}

func LoadConfiguration(app *app.Application) error {
    // Implementation for loading configuration
    return nil
}

func RegisterProviders(app *app.Application) error {
    // Implementation for registering providers
    return nil
}

func BootProvider(app *app.Application) error {
    // Implementation for booting providers
    return nil
}
