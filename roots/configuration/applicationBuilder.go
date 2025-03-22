package configuration

import "gokatan/roots";

type ApplicationBuilder struct {
    app *app.Application
}

func NewApplicationBuilder() *ApplicationBuilder {
    return &ApplicationBuilder{
        app: app.New(),
    }
}

func (b *ApplicationBuilder) Build() *app.Application {
    return b.app;
}

func (b *ApplicationBuilder) WithConfig(key string, value interface{}) *ApplicationBuilder {
    b.app.Bind(key, value);
    return b;
}

func (b *ApplicationBuilder) WithKernel(key string, value interface{}) *app.Application {
    b.app.Singleton(key, value);
    return b.app;
}


func (b *ApplicationBuilder) Boot() *app.Application {
    b.app.Boot();
    return b.app;
}
