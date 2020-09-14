package executor

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/spyzhov/healthy/internal"
	"github.com/spyzhov/healthy/step"
	"github.com/spyzhov/safe"
)

type Executor struct {
	mu      sync.Mutex
	ctx     context.Context
	version string

	setters     map[string]func() (interface{}, error)
	connections map[string]interface{}
}

func NewExecutor(ctx context.Context, version string) *Executor {
	return &Executor{
		ctx:         ctx,
		version:     version,
		setters:     make(map[string]func() (interface{}, error)),
		connections: make(map[string]interface{}),
	}
}

// Get will return a result of any public method from Executor.
// Method will be found by its name, example: string `get name` will call method `GetName`.
// All necessary arguments will be taken from the args variable and placed innumerable: `GetName(args[0], args[1], ...)`
func Get(e *Executor, name string, args []interface{}) (step.Function, error) {
	// region Method
	methodName := getMethodName(name)
	method := reflect.ValueOf(e).MethodByName(methodName)
	if !method.IsValid() {
		return nil, fmt.Errorf(
			"method not found: %s\nAvailable list:\n\t%s",
			name,
			strings.Join(getMethodNames(e), "\n\t"),
		)
	}
	// endregion
	// region Argv
	argv, err := getMethodArguments(methodName, &method, args)
	if err != nil {
		return nil, err
	}
	// endregion
	// region Call
	result := method.Call(argv)
	if len(result) == 2 {
		if err := result[1].Interface(); !safe.IsNil(err) {
			return nil, fmt.Errorf("%v", err)
		}
		if fn, ok := result[0].Interface().(step.Function); ok {
			return fn, nil
		}
	}
	// endregion
	return nil, fmt.Errorf("method `Executor.%s` has wrong declaration, it should be `func MethodName(args map[string]string) (step.Function, error)`", methodName)
}

func getMethodName(name string) string {
	return strings.ReplaceAll(strings.Title(strings.ReplaceAll(strings.Title(name), "/", " ")), " ", "")
}

func getMethodNames(e *Executor) []string {
	ref := reflect.TypeOf(e)
	result := make([]string, 0)
	for i := 0; i < ref.NumMethod(); i++ {
		result = append(result, ref.Method(i).Name)
	}
	return result
}

func getMethodArguments(name string, method *reflect.Value, args []interface{}) (argv []reflect.Value, err error) {
	t := method.Type()
	if len(args) > t.NumIn() {
		return nil, fmt.Errorf("method %s should have no more than %d arguments", name, t.NumIn())
	}
	argv = make([]reflect.Value, t.NumIn())
	for i := range argv {
		value := reflect.New(t.In(i).Elem()).Interface()
		if len(args) > i {
			if err = internal.Unmarshal([]string{"args"}, &value, args[i]); err != nil {
				return nil, err
			}
		}
		argv[i] = reflect.ValueOf(value)
	}
	return argv, nil
}

// region Executor

// Close is an io.Closer function
func (e *Executor) Close() error {
	for id, connection := range e.connections {
		if closer, ok := connection.(io.Closer); ok {
			defer safe.Close(closer, "Executor:connection:"+id)
		}
	}
	return nil
}

func (e *Executor) addSetter(id string, setter func() (interface{}, error)) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.setters[id] = setter
}

func (e *Executor) getConnection(id string) (interface{}, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if value, ok := e.connections[id]; ok {
		return value, nil
	}
	if setter, ok := e.setters[id]; ok {
		if value, err := setter(); err != nil {
			return nil, err
		} else {
			e.connections[id] = value
			return value, nil
		}
	}
	return nil, fmt.Errorf("connection not found")
}

// endregion
