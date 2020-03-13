package console

import "context"

type RouterHandlerFunc func(ctx context.Context, args []string) error

type RouterEndpoint struct{}

func (r RouterEndpoint) Resolved() bool {
	return true
}
func (r RouterEndpoint) Finished() bool {
	return true
}

type RouterHandlerFuncWrapper struct {
	RouterEndpoint
	ctx  context.Context
	args []string
	fn   RouterHandlerFunc
}

func (r *RouterHandlerFuncWrapper) Execute() error {
	return r.fn(r.ctx, r.args)
}
