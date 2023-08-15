package app

import (
	"context"
	"github.com/lishimeng/app-starter/factory"
)

func (h *application) applyComponents(components []func(ctx context.Context) (err error)) (err error) {

	for _, componentFunc := range components {
		err = componentFunc(factory.GetCtx())
		if err != nil {
			break
		}
	}
	return
}

func (h *ApplicationBuilder) ComponentBefore(component func(context.Context) (err error)) *ApplicationBuilder {

	if component != nil {
		h.componentsBeforeWebServer = append(h.componentsBeforeWebServer, component)
	}
	return h
}

func (h *ApplicationBuilder) ComponentAfter(component func(context.Context) (err error)) *ApplicationBuilder {

	if component != nil {
		h.componentsAfterWebServer = append(h.componentsAfterWebServer, component)
	}
	return h
}
