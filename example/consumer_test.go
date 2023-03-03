//go:build example

package playground_test

import (
	"context"
	playground "github.com/panpito/momock/example"
	"github.com/panpito/momock/momock"
	"testing"
)

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

func Test_WrongArgumentsToTheMock(t *testing.T) {
	// given
	ctx := context.TODO()

	mock := playground.NewMockSomeInterface(t)
	mock.MockManager.SetExpectations([]momock.MockCall{
		{
			Method: mock.Do,
			In:     []any{ctx, "GOODBYE"},
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

func Test_WrongNumberOfArguments(t *testing.T) {
	// given
	ctx := context.TODO()

	mock := playground.NewMockSomeInterface(t)
	mock.MockManager.SetExpectations([]momock.MockCall{
		{
			Method: mock.Do,
			In:     []any{ctx, "hello", "EXTRA HELLO"},
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

func Test_ForgotToSetExpectations(t *testing.T) {
	// given
	ctx := context.TODO()

	mock := playground.NewMockSomeInterface(t)

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

func Test_ExpectationsWereNotMet(t *testing.T) {
	// given
	ctx := context.TODO()

	mock := playground.NewMockSomeInterface(t)
	mock.MockManager.SetExpectations([]momock.MockCall{
		{
			Method: mock.Do,
			In:     []any{ctx, "hello"},
			Out:    []any{nil, nil},
		},
		{
			Method: mock.Do,
			In:     []any{ctx, "GOODNIGHT"},
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
