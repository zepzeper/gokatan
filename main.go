package main

import (
    "gokatan/roots/configuration"
    "net/http"
)

func main() {
    // Create the application and boot it
    app := configuration.NewApplicationBuilder().WithKernel().Boot()

    app.HandleRequest(w http.ResponseWriter, r *http.Request);
    
}
