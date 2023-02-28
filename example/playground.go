//go:generate go run $GOPATH/src/momock/main.go

package playground

import "context"

type SomeInterface interface {
	Do(ctx context.Context, x string) (context.Context, error)
}

type AnotherOne interface {
	Badu(x string)
}
