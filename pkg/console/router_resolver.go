package console

import "context"

type RouterResolve struct {
	router *Router
	Ctx    context.Context
	Args   []string
	Fn     RouterHandlerFunc
}

func resolveRouter(router *Router, ctx context.Context, args []string) (RouterResolve, bool) {
	r := RouterResolve{
		router: router,
		Ctx:    ctx,
		Args:   args,
	}
	if match, ctx, args := r.router.matchArgs(ctx, args); match {
		r.Ctx = ctx
		r.Args = args
	} else {
		return r, false
	}
	for _, next := range r.router.next {
		if resolver, resolved := resolveRouter(next, r.Ctx, r.Args); resolved {
			return resolver, true
		}
	}
	for _, fn := range r.router.handlerFuncs {
		if match, ctx, args := fn.matcher.match(r.Ctx, r.Args); match {
			r.Ctx = ctx
			r.Args = args
			r.Fn = fn.handler
			return r, true
		}
	}
	return r, false
}
