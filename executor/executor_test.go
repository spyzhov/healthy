package executor

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/spyzhov/healthy/step"
	"github.com/spyzhov/safe"
)

func call(fn step.Function) []interface{} {
	if safe.IsNil(fn) {
		return []interface{}{
			nil,
			fmt.Errorf("function is nil"),
		}
	}
	res, err := fn()
	return []interface{}{res, err}
}

func validate(t *testing.T, got step.Function, want step.Function) {
	if safe.IsNil(got) && safe.IsNil(want) {
		return
	}
	if safe.IsNil(got) && !safe.IsNil(want) {
		t.Error("got nil, want result")
		return
	}
	if !safe.IsNil(got) && safe.IsNil(want) {
		t.Error("got result, want nil")
		return
	}
	aRes, aErr := got()
	wRes, wErr := want()
	if aErr == nil && wErr != nil {
		t.Error("Function() error:got nil, want result")
		return
	}
	if aErr != nil && wErr == nil {
		t.Error("Function() error:got result, want nil")
		return
	}
	if aErr != nil && wErr != nil {
		if wErr.Error() != aErr.Error() {
			t.Errorf("Function() error: got = %v, want %v", wErr, aErr)
		}
	}
	if aRes == nil && safe.IsNil(wRes) {
		return
	}
	if aRes == nil && wRes != nil {
		t.Error("Function() result: got nil, want result")
		return
	}
	if aRes != nil && wRes == nil {
		t.Error("Function() result: got result, want nil")
		return
	}
	if aRes != nil && wRes != nil {
		if aRes.Status != wRes.Status {
			t.Errorf("Function() status: got = %v, want %v", aRes.Status, wRes.Status)
		}
		if aRes.Message != wRes.Message {
			t.Errorf("Function() message: got = %v, want %v", aRes.Message, wRes.Message)
		}
	}
}

func TestGet(t *testing.T) {
	executor := NewExecutor(context.Background())
	tests := []struct {
		e       *Executor
		name    string
		args    []interface{}
		want    step.Function
		wantErr bool
	}{
		// region Error
		{
			e:       executor,
			name:    "wrong_name",
			args:    nil,
			want:    nil,
			wantErr: true,
		},
		// endregion
		// region Simple
		{
			e:    executor,
			name: "simple",
			args: []interface{}{
				map[string]interface{}{
					"status":  "success",
					"message": "OK",
				},
			},
			want: func() (*step.Result, error) {
				return step.NewResultSuccess("OK"), nil
			},
			wantErr: false,
		},
		{
			e:    executor,
			name: "simple",
			args: []interface{}{
				map[string]interface{}{
					"status":  "error",
					"message": "custom error",
				},
			},
			want: func() (*step.Result, error) {
				return step.NewResultError("custom error"), nil
			},
			wantErr: false,
		},
		{
			e:    executor,
			name: "simple",
			args: []interface{}{
				map[string]interface{}{
					"status":  "warning",
					"message": "custom warning",
				},
			},
			want: func() (*step.Result, error) {
				return step.NewResultWarning("custom warning"), nil
			},
			wantErr: false,
		},
		// endregion Simple
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%02d_%s", i, tt.name), func(t *testing.T) {
			got, err := Get(tt.e, tt.name, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got == nil {
				t.Errorf("Get() got = %v, want <func>", got)
			} else if !reflect.DeepEqual(call(got), call(tt.want)) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
