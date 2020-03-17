package console

import (
	"os"
	"strings"
)

type Router struct {
	path    string
	match   RouterMatchHandlerFunc
	next    []*Router
	resolve RouterResolver
}

func NewRouter(path string) *Router {
	return &Router{path: path}
}
func NewAppRouter() *Router {
	return &Router{
		path:    strings.ToLower(os.Args[0]),
		match:   RouterMatchIgnoreCaseStringHandler(os.Args[0]),
		next:    nil,
		resolve: nil,
	}
}
func (r *Router) SetMatcherFunc(fn RouterMatchHandlerFunc) *Router {
	r.match = fn
	return r
}
func (r *Router) AddNextRouter(router *Router) *Router {
	r.next = append(r.next, router)
	return r
}
func (r *Router) PutNextRouter(router *Router) *Router {
	r.next = append(r.next, router)
	return router
}
func (r *Router) PatchNextRouter(router *Router) *Router {
	for _, next := range r.next {
		if next.path == router.path {
			return next
		}
	}
	r.next = append(r.next, router)
	return router
}

func (r *Router) CreateNextPath(path []string) *Router {
	resultRouter := r
	for _, pattern := range path {
		if next, err := buildNext(pattern); err == nil {
			resultRouter = resultRouter.PatchNextRouter(NewRouter(pattern).SetMatcherFunc(next))
		}
	}
	return resultRouter
}
func (r *Router) CreateNext(path string) *Router {
	return r.CreateNextPath(strings.Split(path, " "))
}
func (r *Router) AddNextPath(path []string) *Router {
	_ = r.CreateNextPath(path)
	return r
}
func (r *Router) AddNext(path string) *Router {
	return r.AddNextPath(strings.Split(path, " "))
}

func (r *Router) SetResolver(resolve RouterResolver) *Router {
	r.resolve = resolve
	return r
}

func (r *Router) SetHandlerFunc(fn RouterHandlerFunc) *Router {
	return r.SetResolver(newRouterHandlerFuncEndpoint(fn))
}
func (r *Router) AddHandlerFuncPath(path []string, fn RouterHandlerFunc) *Router {
	return r.CreateNextPath(path).SetHandlerFunc(fn)
}
func (r *Router) AddHandlerFunc(path string, fn RouterHandlerFunc) *Router {
	return r.CreateNext(path).SetHandlerFunc(fn)
}
