package config

import (
    "github.com/joho/godotenv"
    "os"
)

type EnvLoader struct {
    path string
    loaded bool
}

func NewEnvLoader(path string) *EnvLoader {
    if path == "" {
        path = ".env";
    }

    return &EnvLoader {
        path: path,
        loaded: false,
    }
}

// loads environment variables from .env file
func (e *EnvLoader) Load() error {
    if e.loaded {
        return nil
    }

    err := godotenv.Load(e.path);
    if err != nil && !os.IsNotExist(err) {
        return err;
    }

    e.loaded = true;
    return nil;
}

func Get(key string, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value;
    }

    return defaultValue;
}

// func GetBool(key string, defaultValue bool) bool {
//     if value, exists := os.LookupEnv(key); exists {
//         return value;
//     }
//
//     return defaultValue;
// }
