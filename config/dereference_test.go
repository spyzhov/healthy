package config

import (
	"os"
	"reflect"
	"testing"
)

func Test_dereference(t *testing.T) {
	values := map[string]string{
		"__TEST__":  "test value",
		"__CHECK__": "check value",
	}
	for key, value := range values {
		err := os.Setenv(key, value)
		if err != nil {
			t.Errorf("Setenv(): %s", err)
		}
	}

	tests := []struct {
		name    string
		content []byte
		want    []byte
	}{
		{
			name:    "Clear",
			content: []byte(`clear value`),
			want:    []byte(`clear value`),
		},
		{
			name:    "__TEST__",
			content: []byte(`value with env(__TEST__)`),
			want:    []byte(`value with test value`),
		},
		{
			name:    "__TEST__/__CHECK__",
			content: []byte(`value with env(__TEST__) and __CHECK__`),
			want:    []byte(`value with test value and __CHECK__`),
		},
		{
			name:    "__TEST__/env()",
			content: []byte(`value with env(__TEST__) and env()`),
			want:    []byte(`value with test value and env()`),
		},
		{
			name:    "__TEST__/env(Wrong value ...)",
			content: []byte(`value with env(__TEST__) and env(Wrong value ...)`),
			want:    []byte(`value with test value and `),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dereference(tt.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dereference() = %v, want %v", got, tt.want)
			}
		})
	}
}
