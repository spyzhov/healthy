package executor

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/spyzhov/healthy/step"
)

func call(fn step.Function) []interface{} {
	res, err := fn()
	return []interface{}{res, err}
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
