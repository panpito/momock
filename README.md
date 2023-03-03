# Momock

## About

Generating Go mocks for testing.

### Features / Roadmap

| Feature                                             | Status             |
|-----------------------------------------------------|--------------------|
| Creating mock file                                  | :white_check_mark: |
| Choosing mock file location                         | :x:                |
| Choosing mock file package name                     | :x:                |
| Choosing mock struct name                           | :x:                |
| Assigning build tag                                 | :soon:             |
| Choosing assertion framework                        | :x:                |
| Matching statically typed arguments                 | :white_check_mark: |
| Returning statically typed outputs                  | :white_check_mark: |
| Matching multiple invocations                       | :white_check_mark: |
| Erroring if expectations were set but not triggered | :soon:             |
| Erroring if expectations were not set but triggered | :white_check_mark: |

## Usage

Place that generate directive at the top of file containing the interface you want to generate a mock for:
```go
//go:generate go run github.com/panpito/momock

package playground

import (
	"context"
)

type SomeInterface interface {
	Do(ctx context.Context, x string) (context.Context, error)
}
```

Use native `generate` go command:
```shell
go generate ./...
```
It will generate a mock file (`*_mock.go`) alongside your interface file. Please do not modify it.

Let's say you have a service using the previous interface:
```go
package playground

import "context"

type MyService struct {
	TheInterface SomeInterface
}

func (receiver MyService) UseInterface(ctx context.Context) error {
	_, err := receiver.TheInterface.Do(ctx, "hello")

	return err
}
```

Happy path:
```go
func Test_Sucess(t *testing.T) {
	// given
	ctx := context.TODO()

	mock := playground.NewMockSomeInterface(t)
	mock.MockManager.SetExpectations([]momock.MockCall{
		{
			Method: mock.Do,
			In:     []any{ctx, "hello"},
			Out:    []any{nil, nil},
		},
	})

	myStruct := playground.MyService{
		TheInterface: mock,
	}

	// when
	result := myStruct.UseInterface(ctx)

	// then
	if result != nil {
		t.Fatalf("was expecting nil")
	}
}
```

Forgetting to setup the mock expectation, will generate that error:
```go
=== RUN   Test_ForgotToSetExpectations
    manager.go:45: Was not expecting any mock calls
```

Not matching the arguments, will generate that error:
```go
=== RUN   Test_WrongArgumentsToTheMock
    manager.go:70: 
        Argument: 	1
        Got: 		GOODBYE
        Wanted: 	hello
```

Providing the wrong number of arguments, will generate that error:
```go
=== RUN   Test_WrongNumberOfArguments
    manager.go:57: Wrong number of inputs
```

## Installation

At the root of the project, extend your `tools.go` file:
```go
//go:build tools

package your_project

import (
	_ "github.com/panpito/momock"
)
```

And then `go mod download`.

## Contributing

Please free to fork, create a branch and open a pull request.

## License

This is under MIT license.

## Contact

Please contact:
[Twitter](https://twitter.com/Panpit0)
