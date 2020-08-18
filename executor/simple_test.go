package executor

import (
	"context"
	"reflect"
	"testing"
	"time"

	. "github.com/spyzhov/healthy/executor/internal"
	"github.com/spyzhov/healthy/step"
)

func TestExecutor_Simple(t *testing.T) {
	type args struct {
		args *SimpleArgs
	}
	tests := []struct {
		name     string
		executor *Executor
		args     args
		want     step.Function
		wantErr  bool
	}{
		{
			name:     "nil",
			executor: NewExecutor(context.Background()),
			args: args{
				args: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:     "blank",
			executor: NewExecutor(context.Background()),
			args: args{
				args: &SimpleArgs{
					Sleep:   Duration{},
					Status:  "",
					Message: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:     "valid(success)",
			executor: NewExecutor(context.Background()),
			args: args{
				args: &SimpleArgs{
					Sleep:   Duration{},
					Status:  "success",
					Message: "success message",
				},
			},
			want: func() (*step.Result, error) {
				return step.NewResultSuccess("success message"), nil
			},
			wantErr: false,
		},
		{
			name:     "valid(warning)",
			executor: NewExecutor(context.Background()),
			args: args{
				args: &SimpleArgs{
					Sleep:   Duration{},
					Status:  "warning",
					Message: "warning message",
				},
			},
			want: func() (*step.Result, error) {
				return step.NewResultWarning("warning message"), nil
			},
			wantErr: false,
		},
		{
			name:     "valid(error)",
			executor: NewExecutor(context.Background()),
			args: args{
				args: &SimpleArgs{
					Sleep:   Duration{},
					Status:  "error",
					Message: "error message",
				},
			},
			want: func() (*step.Result, error) {
				return step.NewResultError("error message"), nil
			},
			wantErr: false,
		},
		{
			name:     "valid(success)_sleep",
			executor: NewExecutor(context.Background()),
			args: args{
				args: &SimpleArgs{
					Sleep: Duration{
						Duration: time.Nanosecond,
					},
					Status:  "success",
					Message: "",
				},
			},
			want: func() (*step.Result, error) {
				return step.NewResultSuccess("success"), nil
			},
			wantErr: false,
		},
		{
			name:     "invalid(status)",
			executor: NewExecutor(context.Background()),
			args: args{
				args: &SimpleArgs{
					Sleep:   Duration{},
					Status:  "status",
					Message: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.executor.Simple(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Simple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(call(got), call(tt.want)) {
				t.Errorf("Simple() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleArgs_Validate(t *testing.T) {
	tests := []struct {
		name    string
		a       *SimpleArgs
		wantErr bool
	}{
		{
			name:    "nil",
			a:       nil,
			wantErr: true,
		},
		{
			name: "blank",
			a: &SimpleArgs{
				Sleep:   Duration{},
				Status:  "",
				Message: "",
			},
			wantErr: true,
		},
		{
			name: "valid(success)",
			a: &SimpleArgs{
				Sleep:   Duration{},
				Status:  "success",
				Message: "",
			},
			wantErr: false,
		},
		{
			name: "valid(warning)",
			a: &SimpleArgs{
				Sleep:   Duration{},
				Status:  "warning",
				Message: "warning",
			},
			wantErr: false,
		},
		{
			name: "valid(error)",
			a: &SimpleArgs{
				Sleep:   Duration{Duration: 1},
				Status:  "error",
				Message: "error",
			},
			wantErr: false,
		},
		{
			name: "invalid(sleep)",
			a: &SimpleArgs{
				Sleep:   Duration{Duration: -1},
				Status:  "success",
				Message: "",
			},
			wantErr: true,
		},
		{
			name: "invalid(status)",
			a: &SimpleArgs{
				Sleep:   Duration{},
				Status:  "status",
				Message: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
