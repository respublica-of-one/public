package console

import (
	"context"
	"fmt"
	"strings"
)

type argsMatcher struct {
	handlers []matchHandler
}

func newArgsMatcher(args string) argsMatcher {
	var m argsMatcher
	for _, arg := range strings.Split(args, " ") {
		m.handlers = append(m.handlers, newMatchHandler(arg))
	}
	return m
}
func (m argsMatcher) match(sourceCtx context.Context, args []string) (bool, context.Context, []string) {
	if len(m.handlers) == 0 {
		return true, sourceCtx, args
	}
	ctx := sourceCtx
	match := false
	i := 0
	for index, handler := range m.handlers {
		i = index
		arg := ""
		if len(args) > index {
			arg = args[index]
		}
		match, ctx = handler.handleMatch(ctx, arg)
		if !match {
			return false, sourceCtx, args
		}
	}
	return true, ctx, args[i:]
}
func (m argsMatcher) display() string {
	result := make([]string, 0, len(m.handlers))
	for _, handler := range m.handlers {
		result = append(result, handler.display())
	}
	return strings.Join(result, " ")
}

type matchHandler interface {
	handleMatch(ctx context.Context, arg string) (bool, context.Context)
	display() string
}

func newMatchHandler(arg string) matchHandler {
	if strings.HasPrefix(arg, contextItemPrefix) {
		return newContextItemMatchHandler(arg)
	}
	return newPathItemMatchHandler(arg)
}

const pathItemVariantsSeparator string = "|"

type pathItemMatchHandler []string

func newPathItemMatchHandler(arg string) pathItemMatchHandler {
	return strings.Split(arg, pathItemVariantsSeparator)
}
func (m pathItemMatchHandler) handleMatch(ctx context.Context, arg string) (bool, context.Context) {
	for _, item := range m {
		if item == arg {
			return true, ctx
		}
	}
	return false, ctx
}
func (m pathItemMatchHandler) display() string {
	return fmt.Sprintf("<%s>", strings.Join(m, pathItemVariantsSeparator))
}

const contextItemPrefix = "$"

type contextItemMatchHandler string

func newContextItemMatchHandler(arg string) contextItemMatchHandler {
	return contextItemMatchHandler(arg[1:])
}
func (c contextItemMatchHandler) handleMatch(ctx context.Context, arg string) (bool, context.Context) {
	return true, context.WithValue(ctx, c, arg)
}
func (c contextItemMatchHandler) display() string {
	return fmt.Sprintf("<$%s>", c)
}
