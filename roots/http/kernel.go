package http

import (
	"gokatan/roots"
	"gokatan/roots/bootstrap"
	"net/http"
	"time"
)

type Kernel struct {
    requestStartTime time.Time
    app *app.Application
    bootstrappers []func(*app.Application) error
}

func NewKernel(application *app.Application) *Kernel {
    return &Kernel{
        app: application,
        bootstrappers: []func(*app.Application) error{
            bootstrap.LoadEnvironmentVariables,
            bootstrap.LoadConfiguration,
            bootstrap.RegisterProviders,
            bootstrap.BootProvider,
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

func (k *Kernel) AddBootstrapper(bootstrapper func(*app.Application) error) {
    k.bootstrappers = append(k.bootstrappers, bootstrapper)
}

func (k* Kernel) Handle(r *http.Request) error {
    k.requestStartTime = time.Now();

    k.Bootstrap();

    k.app.Singleton("request", r);

    resp, err := k.sendRequestThroughPipeline(r);

    if err != nil {
        return err
    }

    return resp

}

func (k *Kernel) sendRequestThroughPipeline(r *http.Request) (interface{}, error) {
    // Implement your request pipeline logic here
    // This is a placeholder implementation
    return nil, nil
}
