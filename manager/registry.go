package momock_manager

import "errors"

type MockCall struct {
	Method interface{}
	In     []any
	Out    []any
}

type MockCallRegister struct {
	mocks []MockCall
}

func (register *MockCallRegister) Size() int {
	return len(register.mocks)
}

func (register *MockCallRegister) Push(mock MockCall) {
	register.mocks = append(register.mocks, mock)
}

func (register *MockCallRegister) Pop() (MockCall, error) {
	if register.Size() == 0 {
		return MockCall{}, errors.New("empty register")
	}

	firstItem := register.mocks[0]
	register.mocks = register.mocks[1:len(register.mocks)]

	return firstItem, nil
}
