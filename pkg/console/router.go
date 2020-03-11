package console

import (
	"context"
)

type Router struct {
	parent       *Router
	args         []argsMatcher
	next         []*Router
	handlerFuncs []RouterHandlerFuncMatcher
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) SetParent(parent *Router) *Router {
	r.parent = parent
	return r
}

func (r *Router) AddNext(router *Router) *Router {
	router.parent = r
	r.next = append(r.next, router)
	return r
}
func (r *Router) AddNextRoute(args string, router *Router) *Router {
	router.SetArgs(args)
	return r.AddNext(router)
}

func (r *Router) CreateNext() *Router {
	router := NewRouter().SetParent(r)
	r.next = append(r.next, router)
	return router
}
func (r *Router) CreateNextRoute(args string) *Router {
	return r.CreateNext().SetArgs(args)
}

func (r *Router) SetArgs(args string) *Router {
	r.args = []argsMatcher{newArgsMatcher(args)}
	return r
}
func (r *Router) AddArgs(args string) *Router {
	r.args = append(r.args, newArgsMatcher(args))
	return r
}

func (r *Router) SetHandlerFunc(h RouterHandlerFunc) *Router {
	r.handlerFuncs = append(r.handlerFuncs, RouterHandlerFuncMatcher{handler: h})
	return r
}
func (r *Router) AddHandlerFunc(args string, h RouterHandlerFunc) *Router {
	r.handlerFuncs = append(r.handlerFuncs, RouterHandlerFuncMatcher{matcher: newArgsMatcher(args), handler: h})
	return r
}

func (r *Router) matchArgs(ctx context.Context, args []string) (bool, context.Context, []string) {
	for _, arg := range r.args {
		if match, newCtx, newArgs := arg.match(ctx, args); match {
			ctx = newCtx
			args = newArgs
		} else {
			return false, ctx, args
		}
	}
	return true, ctx, args
}

func (r *Router) Resolve(ctx context.Context, args []string) (RouterResolve, bool) {
	return resolveRouter(r, ctx, args)
}
