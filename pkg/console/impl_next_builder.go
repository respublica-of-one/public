package console

import (
	"sort"
	"strings"
)

type RouterNextPatternMatcher interface {
	HandlePatternMatch(pattern string) (string, bool)
}

type RouterNextMatchSuffix string

func (suffix RouterNextMatchSuffix) HandlePatternMatch(pattern string) (string, bool) {
	if strings.HasSuffix(pattern, string(suffix)) {
		return strings.Replace(pattern, string(suffix), "", 1), true
	}
	return pattern, false
}

type RouterNextMatchPrefix string

func (prefix RouterNextMatchPrefix) HandlePatternMatch(pattern string) (string, bool) {
	if strings.HasPrefix(pattern, string(prefix)) {
		return pattern[:len(pattern)-len(prefix)], true
	}
	return pattern, false
}

type RouterNextMatchBool bool

func (b RouterNextMatchBool) HandlePatternMatch(pattern string) (string, bool) {
	return pattern, bool(b)
}

type RouterNextBuilderFunc func(route string) RouterNextAddCreator
type routerNextBuilderFuncSearchWrapper struct {
	weight  int
	fn      RouterNextBuilderFunc
	pattern RouterNextPatternMatcher
}

func RegisterRouterNextBuilderFunc(pattern RouterNextPatternMatcher, weight int, builderFunc RouterNextBuilderFunc) {
	nextBuilder = append(nextBuilder, routerNextBuilderFuncSearchWrapper{
		pattern: pattern,
		weight:  weight,
		fn:      builderFunc,
	})
	sort.Slice(nextBuilder, func(i, j int) bool {
		return nextBuilder[i].weight > nextBuilder[j].weight
	})
}
func getNextBuilder(route string) (RouterNextBuilderFunc, string, bool) {
	for _, builder := range nextBuilder {
		if val, match := builder.pattern.HandlePatternMatch(route); match {
			return builder.fn, val, true
		}
	}
	return nil, route, false
}
func buildNext(route string) (RouterNextAddCreator, bool) {
	if builder, val, found := getNextBuilder(route); found {
		return builder(val), true
	}
	return nil, false
}

var nextBuilder []routerNextBuilderFuncSearchWrapper

func init() {

}
