package console

import "strings"

type RouterNextMatchFunc func(pattern string) (string, bool)

func RouterNextMatchBool(val bool) RouterNextMatchFunc {
	return func(pattern string) (string, bool) {
		return pattern, val
	}
}
func RouterNextMatchString(val string) RouterNextMatchFunc {
	return func(pattern string) (string, bool) {
		if val == pattern {
			return pattern, true
		}
		return pattern, false
	}
}
func RouterNextMatchPrefix(prefix string) RouterNextMatchFunc {
	return func(pattern string) (string, bool) {
		if strings.HasPrefix(pattern, prefix) {
			return strings.Replace(pattern, prefix, "", 1), true
		}
		return pattern, false
	}
}
