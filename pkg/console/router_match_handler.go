package console

import (
	"context"
	"strings"
)

type RouterMatchHandlerFunc func(c *RouterContext, arg string) (bool, error)
type RouterMatchHandlerBuilderFunc func(pattern string) RouterMatchHandlerFunc

func RouterMatchStringHandler(pattern string) RouterMatchHandlerFunc {
	return func(c *RouterContext, arg string) (bool, error) {
		if arg == pattern {
			return true, nil
		}
		return false, nil
	}
}
func RouterMatchIgnoreCaseStringHandler(pattern string) RouterMatchHandlerFunc {
	return func(c *RouterContext, arg string) (bool, error) {
		if strings.ToLower(arg) == pattern {
			return true, nil
		}
		return false, nil
	}
}
func RouterMatchCtxHandler(pattern string) RouterMatchHandlerFunc {
	return func(c *RouterContext, arg string) (bool, error) {
		c.Ctx = context.WithValue(c.Ctx, pattern, arg)
		return true, nil
	}
}
