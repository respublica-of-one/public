package console

import "context"

type RouterEndpoint interface {
	Execute() error
}

type RouterHandlerFunc func(ctx context.Context, args []string) error

type routerHandlerFuncEndpoint struct {
	ctx  context.Context
	args []string
	fn   RouterHandlerFunc
}

func newRouterHandlerFuncEndpoint(fn RouterHandlerFunc) routerHandlerFuncEndpoint {
	return routerHandlerFuncEndpoint{
		fn: fn,
	}
}
func (e routerHandlerFuncEndpoint) Resolve(c *RouterContext) (RouterEndpoint, error) {
	e.ctx = c.Ctx
	e.args = c.Args
	return e, nil
}
func (e routerHandlerFuncEndpoint) Execute() error {
	return e.fn(e.ctx, e.args)
}
