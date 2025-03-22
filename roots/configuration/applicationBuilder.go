package configuration

type ApplicationBuilder struct {
    app *Application
}

func NewApplicationBuilder() *ApplicationBuilder {
    return &ApplicationBuilder{
        app: New(),
    }
}

func (b *ApplicationBuilder) Build() *ApplicationBuilder {
    return b.app;
}

func (b *ApplicationBuilder) withConfig(key string, value interface{}) *ApplicationBuilder {
    b.app.Bind(key, value);
    return b;
}

func (b *ApplicationBuilder) Boot() *Application {
    b.app.Boot;
    return b.app;
}
