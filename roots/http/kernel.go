package http

import (
    "gokatan/contracts"
    "net/http"
    "time"
)

type Kernel struct {
    requestStartTime time.Time
    app contracts.IApplication
    router contracts.IRouter
    bootstrappers []func(contracts.IApplication) error
}

func NewKernel(application contracts.IApplication) *Kernel {
    return &Kernel{
        app: application,
        bootstrappers: []func(contracts.IApplication) error{
        },
    }
}

func (k *Kernel) Bootstrap() error {
    for _, bootstrapper := range k.bootstrappers {
        err := bootstrapper(k.app)
        if err != nil {
            return err
        }
    }
    return nil
}

func (k *Kernel) AddBootstrapper(bootstrapper func(contracts.IApplication) error) {
    k.bootstrappers = append(k.bootstrappers, bootstrapper)
}

func (k *Kernel) Handle(r *http.Request) error {
    k.requestStartTime = time.Now()

    err := k.Bootstrap()
    if err != nil {
        return err
    }

    k.app.Singleton("request", r)

    // Fixed the return type issue
    _, err = k.sendRequestThroughPipeline(r)
    return err
}

func (k *Kernel) sendRequestThroughPipeline(r *http.Request) (interface{}, error) {
    // Implement your request pipeline logic here
    // This is a placeholder implementation
    return nil, nil
}
