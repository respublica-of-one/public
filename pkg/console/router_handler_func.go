package console

import "context"

type RouterHandlerFunc func(ctx context.Context, args []string) error

type RouterHandlerFuncMatcher struct {
	matcher argsMatcher
	handler RouterHandlerFunc
}
