package playground

import (
	"context"
	"fmt"
	"github.com/panpito/momock/manager"
	"log"
	"testing"
)

type MockSomeInterface struct {
	t           *testing.T
	MockManager *momock_manager.MockManager
}

func NewMockSomeInterface(t *testing.T) *MockSomeInterface {
	return &MockSomeInterface{t: t, MockManager: momock_manager.NewMockManager(t)}
}

func (mock *MockSomeInterface) Do(ctx context.Context, x string) (context.Context, error) {
	callerData := momock_manager.CallerData{MethodName: mock.MockManager.WhatsMyName(), InputsLength: 2, OutputsLength: 2, Inputs: map[int]any{0: ctx, 1: x}}
	out := mock.MockManager.Verify(callerData)
	log.Print(out)
	return0, _ := out[0].(context.Context)
	return1, _ := out[1].(error)
	return return0, return1
}

type MockAnotherOne struct {
	t           *testing.T
	MockManager *momock_manager.MockManager
}

func NewMockAnotherOne(t *testing.T) *MockAnotherOne {
	return &MockAnotherOne{t: t, MockManager: momock_manager.NewMockManager(t)}
}

func (mock *MockAnotherOne) Badu(x string) {
	callerData := momock_manager.CallerData{MethodName: mock.MockManager.WhatsMyName(), InputsLength: 1, OutputsLength: 0, Inputs: map[int]any{0: x}}
	out := mock.MockManager.Verify(callerData)
	log.Print(out)
	return
}
