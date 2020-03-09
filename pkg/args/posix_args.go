package args

import (
	"context"
	"reflect"
)

func FromPosix(ctx context.Context, args []string, val interface{}) error {
	t := reflect.TypeOf(val)

	return nil
}
