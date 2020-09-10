package internal

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"runtime/debug"
	"testing"
	"time"
)

type TestFoo struct {
	Foo int
}

type TestBar struct {
	Bar string `name:"bAr"`
}

type TestFooBar struct {
	TestFoo
	TestBar
}

type TestFooBarBaz struct {
	TestFoo
	TestBar
	Baz float64 `json:"baz"`
}

type TestCustom struct {
	Inner1 struct {
		Inner2 struct {
			Value int
		}
	}
}

type TestDuration struct {
	time.Duration
}

func (d *TestDuration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("invalid duration: %s", string(b))
	}
}

func TestUnmarshal(t *testing.T) {
	var (
		path    = []string{"root"}
		iface1  interface{}
		iface2  interface{} = 1
		slice1  []int
		error1  = fmt.Errorf("example")
		pString *string
		struct1 = struct{}{}
		foo1    TestFoo
		ptr     = func(val interface{}) *interface{} {
			return &val
		}
		ptrStr = func(val string) *string {
			return &val
		}
		pptrStr = func(val *string) **string {
			return &val
		}
	)
	tests := []struct {
		name     string
		target   interface{}
		value    interface{}
		expected interface{}
		wantErr  bool
	}{
		// region Base types
		// region simple
		{
			name:     "simple: int",
			target:   int(0),
			value:    int(123),
			expected: int(123),
			wantErr:  false,
		},
		{
			name:     "simple: uint",
			target:   uint(0),
			value:    uint(123),
			expected: uint(123),
			wantErr:  false,
		},
		{
			name:     "simple: float",
			target:   float64(0),
			value:    float64(123),
			expected: float64(123),
			wantErr:  false,
		},
		{
			name:     "simple: string",
			target:   "",
			value:    "simple",
			expected: "simple",
			wantErr:  false,
		},
		{
			name:     "simple: bool",
			target:   true,
			value:    false,
			expected: false,
			wantErr:  false,
		},
		// endregion
		// region convert
		// region float
		{
			name:     "convert: int -> float32",
			target:   float32(0),
			value:    int(123),
			expected: float32(123),
			wantErr:  false,
		},
		{
			name:     "convert: string -> float32",
			target:   float32(0),
			value:    "123",
			expected: float32(123),
			wantErr:  false,
		},
		{
			name:     "error convert: string -> float32",
			target:   float32(0),
			value:    "one-two-three",
			expected: nil,
			wantErr:  true,
		},
		// endregion
		// region bool
		{
			name:     "convert: string -> bool",
			target:   false,
			value:    "true",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "error convert: string -> bool",
			target:   false,
			value:    "unknown",
			expected: nil,
			wantErr:  true,
		},
		// endregion
		// region int
		{
			name:     "convert: string -> int",
			target:   int(0),
			value:    "123",
			expected: int(123),
			wantErr:  false,
		},
		{
			name:     "convert: float -> int",
			target:   int(0),
			value:    float64(123),
			expected: int(123),
			wantErr:  false,
		},
		{
			name:     "error convert: string -> int",
			target:   int(0),
			value:    "true",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "error convert: struct -> int",
			target:   int(0),
			value:    struct{}{},
			expected: nil,
			wantErr:  true,
		},
		// endregion
		// region uint
		{
			name:     "convert: string -> uint",
			target:   uint(0),
			value:    "123",
			expected: uint(123),
			wantErr:  false,
		},
		{
			name:     "convert: float -> uint",
			target:   uint(0),
			value:    float64(123),
			expected: uint(123),
			wantErr:  false,
		},
		{
			name:     "error convert: string -> uint",
			target:   uint(0),
			value:    "true",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "error convert: struct -> uint",
			target:   uint(0),
			value:    struct{}{},
			expected: nil,
			wantErr:  true,
		},
		// endregion
		// endregion
		// region interface
		{
			name:     "simple: interface{}",
			target:   iface1,
			value:    iface2,
			expected: iface2,
			wantErr:  false,
		},
		{
			name:     "simple: int -> interface{}",
			target:   iface1,
			value:    int(1),
			expected: int(1),
			wantErr:  false,
		},
		{
			name:     "simple: string -> interface{}",
			target:   iface1,
			value:    "int(1)",
			expected: "int(1)",
			wantErr:  false,
		},
		{
			name:     "simple: slice -> interface{}",
			target:   iface1,
			value:    slice1,
			expected: slice1,
			wantErr:  false,
		},
		{
			name:     "simple: array -> interface{}",
			target:   iface1,
			value:    [3]int{},
			expected: [3]int{},
			wantErr:  false,
		},
		// endregion
		// endregion
		// region Slices
		{
			name:     "simple: []int",
			target:   slice1,
			value:    []int{1, 2, 3},
			expected: []int{1, 2, 3},
			wantErr:  false,
		},
		{
			name:     "convert: []float64 -> []int",
			target:   make([]int, 0),
			value:    []float64{1, 2, 3},
			expected: []int{1, 2, 3},
			wantErr:  false,
		},
		{
			name:     "convert: []interface{} -> []int",
			target:   make([]int, 0),
			value:    []interface{}{1, "2", 3.0},
			expected: []int{1, 2, 3},
			wantErr:  false,
		},
		{
			name:     "nullable: []int",
			target:   make([]int, 0),
			value:    nil,
			expected: make([]int, 0),
			wantErr:  false,
		},
		{
			name:     "error convert: []error -> []int",
			target:   make([]int, 0),
			value:    []error{error1},
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "error convert: map[string]int -> []int",
			target:   make([]int, 0),
			value:    make(map[string]int),
			expected: nil,
			wantErr:  true,
		},
		// endregion
		// region Maps
		{
			name:     "simple: map[string]int",
			target:   make(map[string]int),
			value:    map[string]int{"foo": 1, "bar": 2},
			expected: map[string]int{"foo": 1, "bar": 2},
			wantErr:  false,
		},
		{
			name:     "convert: map[int32]float64 -> map[float64]int",
			target:   make(map[float64]int),
			value:    map[int32]float64{1: 2, 3: 4},
			expected: map[float64]int{1: 2, 3: 4},
			wantErr:  false,
		},
		{
			name:     "convert: map[interface{}]interface{} -> map[int]float64",
			target:   make(map[int]float64),
			value:    map[interface{}]interface{}{1: "2", "3": 4.0},
			expected: map[int]float64{1: 2, 3: 4},
			wantErr:  false,
		},
		{
			name:     "nullable: map[string]int",
			target:   make(map[string]int),
			value:    nil,
			expected: make(map[string]int),
			wantErr:  false,
		},
		{
			name:     "error convert: map[error]int -> map[string]int",
			target:   make(map[string]int),
			value:    map[error]int{error1: 1},
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "error convert: map[string]error -> map[string]int",
			target:   make(map[string]int),
			value:    map[string]error{"foo": error1},
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "error convert: []int -> map[string]int",
			target:   make(map[string]int),
			value:    make([]int, 0),
			expected: nil,
			wantErr:  true,
		},
		// endregion
		// region Arrays
		{
			name:     "simple: [3]int",
			target:   [3]int{},
			value:    [3]int{1, 2, 3},
			expected: [3]int{1, 2, 3},
			wantErr:  false,
		},
		{
			name:     "convert: [3]float64 -> [3]int",
			target:   [3]int{},
			value:    [3]float64{1, 2, 3},
			expected: [3]int{1, 2, 3},
			wantErr:  false,
		},
		{
			name:     "convert: [3]interface{} -> [3]int",
			target:   [3]int{},
			value:    [3]interface{}{1, "2", 3.0},
			expected: [3]int{1, 2, 3},
			wantErr:  false,
		},
		{
			name:     "nullable: [3]int",
			target:   [3]int{},
			value:    nil,
			expected: [3]int{},
			wantErr:  false,
		},
		{
			name:     "error convert: [3]error -> [3]int",
			target:   [3]int{},
			value:    []error{},
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "error convert: map[string]int -> [3]int",
			target:   [3]int{},
			value:    make(map[string]int),
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "error convert: [6]int -> [3]int",
			target:   [3]int{},
			value:    [6]int{},
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "error convert: [3]interface{} -> [3]int",
			target:   [3]int{},
			value:    [3]interface{}{1, "2.4", 3.0},
			expected: nil,
			wantErr:  true,
		},
		// endregion
		// region Ptr
		{
			name:     "simple: *string",
			target:   pString,
			value:    ptr("foo"),
			expected: ptrStr("foo"),
			wantErr:  false,
		},
		{
			name:   "simple: *TestFoo",
			target: new(TestFoo),
			value:  map[string]interface{}{"foo": 1},
			expected: &TestFoo{
				Foo: 1,
			},
			wantErr: false,
		},
		{
			name:     "convert simple: string -> *string",
			target:   pString,
			value:    "bar",
			expected: ptrStr("bar"),
			wantErr:  false,
		},
		{
			name:     "convert simple: string -> **string",
			target:   &pString,
			value:    "baz",
			expected: pptrStr(ptrStr("baz")),
			wantErr:  false,
		},
		{
			name:     "error convert: bool -> *string",
			target:   pString,
			value:    true,
			expected: nil,
			wantErr:  true,
		},
		// endregion
		// region Struct
		{
			name:     "nullable: struct{}",
			target:   struct1,
			value:    nil,
			expected: struct1,
			wantErr:  false,
		},
		{
			name:     "nullable: struct(TestFoo)",
			target:   foo1,
			value:    nil,
			expected: foo1,
			wantErr:  false,
		},
		{
			name:   "simple: map -> struct(TestFoo)",
			target: TestFoo{},
			value:  map[string]int{"foo": 123},
			expected: TestFoo{
				Foo: 123,
			},
			wantErr: false,
		},
		{
			name:   "simple: map -> struct(TestBar)",
			target: TestBar{},
			value:  map[string]string{"bAr": "123"},
			expected: TestBar{
				Bar: "123",
			},
			wantErr: false,
		},
		{
			name:   "simple: map -> struct(TestFooBar)",
			target: TestFooBar{},
			value:  map[string]interface{}{"foo": 123, "bAr": "123"},
			expected: TestFooBar{
				TestFoo{
					Foo: 123,
				},
				TestBar{
					Bar: "123",
				},
			},
			wantErr: false,
		},
		{
			name:   "simple: map -> struct(TestFooBarBaz)",
			target: TestFooBarBaz{},
			value:  map[string]interface{}{"foo": 123, "bAr": "123", "baz": "456"},
			expected: TestFooBarBaz{
				TestFoo: TestFoo{
					Foo: 123,
				},
				TestBar: TestBar{
					Bar: "123",
				},
				Baz: 456,
			},
			wantErr: false,
		},
		{
			name:   "simple: map -> struct(TestCustom)",
			target: TestCustom{},
			value:  map[string]interface{}{"inner_1": map[string]interface{}{"inner_2": map[string]interface{}{"value": "1"}}},
			expected: TestCustom{
				Inner1: struct{ Inner2 struct{ Value int } }{Inner2: struct{ Value int }{Value: 1}},
			},
			wantErr: false,
		},
		{
			name:     "error simple: bool -> struct(TestCustom)",
			target:   TestCustom{},
			value:    false,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "error simple: map[int]any -> struct(TestCustom)",
			target:   TestCustom{},
			value:    map[int]int{1: 1},
			expected: nil,
			wantErr:  true,
		},
		// endregion
		// region Unmarshaler
		{
			name:   "unmarshaler: string -> Duration",
			target: new(TestDuration),
			value:  "1s",
			expected: &TestDuration{
				Duration: time.Second,
			},
			wantErr: false,
		},
		{
			name:     "error unmarshaler: bool -> Duration",
			target:   new(TestDuration),
			value:    false,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "error unmarshaler: infinity -> Duration",
			target:   new(TestDuration),
			value:    math.Inf(1),
			expected: nil,
			wantErr:  true,
		},
		// endregion
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if rec := recover(); rec != nil {
					t.Errorf("recovered:\n%v", rec)
					debug.PrintStack()
				}
			}()
			if err := Unmarshal(path, &tt.target, tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil && !tt.wantErr {
				if !reflect.DeepEqual(tt.target, tt.expected) {
					t.Errorf("DeepEqual():\n  actual(%T): %v\nexpected(%T): %v", tt.target, tt.target, tt.expected, tt.expected)
				}
			}
		})
	}
}

func TestTypeToString(t *testing.T) {
	tests := []struct {
		name string
		tp   reflect.Type
		want string
	}{
		{
			name: "int",
			tp:   reflect.TypeOf(int(0)),
			want: "int",
		},
		{
			name: "[]int",
			tp:   reflect.TypeOf(make([]int, 0)),
			want: "[]int",
		},
		{
			name: "[3]int",
			tp:   reflect.TypeOf([3]int{}),
			want: "[3]int",
		},
		{
			name: "map[string]int",
			tp:   reflect.TypeOf(make(map[string]int)),
			want: "map[string]int",
		},
		{
			name: "map[string][]int",
			tp:   reflect.TypeOf(make(map[string][]int)),
			want: "map[string][]int",
		},
		{
			name: "map[string]*[]*int",
			tp:   reflect.TypeOf(make(map[string]*[]*int)),
			want: "map[string]*[]*int",
		},
		{
			name: "*string",
			tp:   reflect.TypeOf(new(string)),
			want: "*string",
		},
		{
			name: "*TestCustom",
			tp:   reflect.TypeOf(new(TestCustom)),
			want: "*TestCustom",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TypeToString(tt.tp); got != tt.want {
				t.Errorf("TypeToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getValue(t *testing.T) {
	ptr := func(val interface{}) *interface{} {
		return &val
	}
	tests := []struct {
		name  string
		value interface{}
		want  interface{}
	}{
		{
			name:  "*string",
			value: ptr("str"),
			want:  "str",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getValue(tt.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
