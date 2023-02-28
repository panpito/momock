package momock_manager

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type MockManager struct {
	t            *testing.T
	mockRegister MockCallRegister
}

func NewMockManager(t *testing.T) *MockManager {
	return &MockManager{t: t}
}

func (manager *MockManager) SetExpectations(mockCalls []MockCall) {
	for _, mockCall := range mockCalls {
		manager.mockRegister.Push(mockCall)
	}
}

func (manager *MockManager) TearDown() {
	if manager.mockRegister.Size() != 0 {
		manager.t.Fatalf("Was expecting mock calls")
	}
}

func (manager *MockManager) WhatsMyName() string {
	pc, _, _, _ := runtime.Caller(1)
	ss := strings.Split(runtime.FuncForPC(pc).Name(), ".")

	if len(ss) != 0 {
		return ss[len(ss)-1]
	}

	return ""
}

func (manager *MockManager) Verify(callerData CallerData) []any {
	if manager.mockRegister.Size() == 0 {
		manager.t.Fatalf("Was not expecting any mock calls")
	}

	mock, err := manager.mockRegister.Pop()
	if err != nil {
		manager.t.Fatalf("Could not get mock: %v", err)
	}

	if mockMethod(runtime.FuncForPC(reflect.ValueOf(mock.Method).Pointer()).Name()) != callerData.MethodName {
		manager.t.Fatalf("\nWas not expecting calls on: \n%s\nbut:\n%s", callerData.MethodName, mockMethod)
	}
	if len(mock.In) != callerData.InputsLength {
		manager.t.Fatalf("Wrong number of inputs")
	}
	if len(callerData.Inputs) != callerData.InputsLength {
		manager.t.Fatalf("Malformed mock for inputs")
	}
	if len(mock.Out) != callerData.OutputsLength {
		manager.t.Fatalf("Wrong number of outputs")
	}

	// Inputs
	// TODO implement any
	for inIdx, val := range callerData.Inputs {
		if arg := mock.In[inIdx]; !reflect.DeepEqual(arg, val) {
			manager.t.Errorf("\nArgument: \t%d\nGot: \t\t%v\nWanted: \t%v", inIdx, arg, val)
		}
	}

	// Outputs
	// TODO implement defaut return
	return mock.Out
}

func mockMethod(method string) string {
	ss := strings.Split(method, ".")
	if len(ss) == 0 {
		return ""
	}

	return strings.TrimSuffix(ss[len(ss)-1], "-fm")
}

type CallerData struct {
	MethodName    string
	InputsLength  int
	Inputs        map[int]any
	OutputsLength int
}
