package contracts

import "net/http"

type IKernel interface {
    Bootstrap() error
    Handle(*http.Request) error
    RegisterRoutes()
}

