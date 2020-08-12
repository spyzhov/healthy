package internal

import "testing"

func TestHttpArgsRequireHeader_Match(t *testing.T) {
	type fields struct {
		Exists    []string
		NotExists []string
		Regexp    map[string]RequireFieldMatch
		NotRegexp map[string]RequireFieldMatchNot
		Eq        map[string][]string
	}
	type args struct {
		header map[string][]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "nil",
			fields: fields{
				Exists:    nil,
				NotExists: nil,
				Regexp:    nil,
				NotRegexp: nil,
				Eq:        nil,
			},
			args: args{
				header: nil,
			},
			wantErr: false,
		},
		{
			name: "blank",
			fields: fields{
				Exists:    make([]string, 0),
				NotExists: make([]string, 0),
				Regexp:    make(map[string]RequireFieldMatch),
				NotRegexp: make(map[string]RequireFieldMatchNot),
				Eq:        make(map[string][]string),
			},
			args: args{
				header: make(map[string][]string),
			},
			wantErr: false,
		},
		{
			name: "exists(valid)",
			fields: fields{
				Exists:    []string{"key"},
				NotExists: nil,
				Regexp:    nil,
				NotRegexp: nil,
				Eq:        nil,
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: false,
		},
		{
			name: "exists(valid2)",
			fields: fields{
				Exists:    []string{"key", "foo"},
				NotExists: nil,
				Regexp:    nil,
				NotRegexp: nil,
				Eq:        nil,
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: false,
		},
		{
			name: "exists(invalid)",
			fields: fields{
				Exists:    []string{"key", "bar"},
				NotExists: nil,
				Regexp:    nil,
				NotRegexp: nil,
				Eq:        nil,
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: true,
		},
		{
			name: "exists(invalid)",
			fields: fields{
				Exists:    []string{"invalid"},
				NotExists: nil,
				Regexp:    nil,
				NotRegexp: nil,
				Eq:        nil,
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: true,
		},
		{
			name: "not_exists(valid)",
			fields: fields{
				Exists:    nil,
				NotExists: []string{"bar"},
				Regexp:    nil,
				NotRegexp: nil,
				Eq:        nil,
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: false,
		},
		{
			name: "not_exists(valid2)",
			fields: fields{
				Exists:    nil,
				NotExists: []string{"bar", "baz"},
				Regexp:    nil,
				NotRegexp: nil,
				Eq:        nil,
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: false,
		},
		{
			name: "not_exists(invalid)",
			fields: fields{
				Exists:    nil,
				NotExists: []string{"foo", "baz"},
				Regexp:    nil,
				NotRegexp: nil,
				Eq:        nil,
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: true,
		},
		{
			name: "not_exists(invalid2)",
			fields: fields{
				Exists:    nil,
				NotExists: []string{"foo", "key"},
				Regexp:    nil,
				NotRegexp: nil,
				Eq:        nil,
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: true,
		},
		{
			name: "regex(valid)",
			fields: fields{
				Exists:    nil,
				NotExists: nil,
				Regexp: map[string]RequireFieldMatch{
					"key": {`value`},
				},
				NotRegexp: nil,
				Eq:        nil,
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: false,
		},
		{
			name: "regex(invalid)",
			fields: fields{
				Exists:    nil,
				NotExists: nil,
				Regexp: map[string]RequireFieldMatch{
					"key": {`wrong value`},
				},
				NotRegexp: nil,
				Eq:        nil,
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: true,
		},
		{
			name: "regex(not_exists)",
			fields: fields{
				Exists:    nil,
				NotExists: nil,
				Regexp: map[string]RequireFieldMatch{
					"wrong": {`value`},
				},
				NotRegexp: nil,
				Eq:        nil,
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: true,
		},
		{
			name: "not_regex(valid)",
			fields: fields{
				Exists:    nil,
				NotExists: nil,
				Regexp:    nil,
				NotRegexp: map[string]RequireFieldMatchNot{
					"key": {`bar`},
				},
				Eq: nil,
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: false,
		},
		{
			name: "not_regex(invalid)",
			fields: fields{
				Exists:    nil,
				NotExists: nil,
				Regexp:    nil,
				NotRegexp: map[string]RequireFieldMatchNot{
					"key": {`.+`},
				},
				Eq: nil,
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: true,
		},
		{
			name: "not_regex(not_exists)",
			fields: fields{
				Exists:    nil,
				NotExists: nil,
				Regexp:    nil,
				NotRegexp: map[string]RequireFieldMatchNot{
					"wrong": {`field`},
				},
				Eq: nil,
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: true,
		},
		{
			name: "eq(valid)",
			fields: fields{
				Exists:    nil,
				NotExists: nil,
				Regexp:    nil,
				NotRegexp: nil,
				Eq: map[string][]string{
					"key": {"value"},
				},
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: false,
		},
		{
			name: "eq(valid2)",
			fields: fields{
				Exists:    nil,
				NotExists: nil,
				Regexp:    nil,
				NotRegexp: nil,
				Eq: map[string][]string{
					"key": {"value"},
					"foo": {"bar", "baz"},
				},
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: false,
		},
		{
			name: "eq(invalid)",
			fields: fields{
				Exists:    nil,
				NotExists: nil,
				Regexp:    nil,
				NotRegexp: nil,
				Eq: map[string][]string{
					"key": {"wrong"},
				},
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: true,
		},
		{
			name: "eq(invalid2)",
			fields: fields{
				Exists:    nil,
				NotExists: nil,
				Regexp:    nil,
				NotRegexp: nil,
				Eq: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "biz"},
				},
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: true,
		},
		{
			name: "eq(not_exists)",
			fields: fields{
				Exists:    nil,
				NotExists: nil,
				Regexp:    nil,
				NotRegexp: nil,
				Eq: map[string][]string{
					"Wrong": {"field"},
				},
			},
			args: args{
				header: map[string][]string{
					"Key": {"value"},
					"Foo": {"bar", "baz"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &HttpArgsRequireHeader{
				Exists:    tt.fields.Exists,
				NotExists: tt.fields.NotExists,
				Regexp:    tt.fields.Regexp,
				NotRegexp: tt.fields.NotRegexp,
				Eq:        tt.fields.Eq,
			}
			if err := a.Match(tt.args.header); (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttpArgsRequireHeader_Validate(t *testing.T) {
	type fields struct {
		Exists    []string
		NotExists []string
		Regexp    map[string]RequireFieldMatch
		NotRegexp map[string]RequireFieldMatchNot
		Eq        map[string][]string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "nil",
			fields: fields{
				Exists:    nil,
				NotExists: nil,
				Regexp:    nil,
				NotRegexp: nil,
				Eq:        nil,
			},
			wantErr: false,
		},
		{
			name: "blank",
			fields: fields{
				Exists:    make([]string, 0),
				NotExists: make([]string, 0),
				Regexp:    make(map[string]RequireFieldMatch),
				NotRegexp: make(map[string]RequireFieldMatchNot),
				Eq:        make(map[string][]string),
			},
			wantErr: false,
		},
		{
			name: "valid_exists",
			fields: fields{
				Exists:    []string{`any`, ``},
				NotExists: nil,
				Regexp:    nil,
				NotRegexp: nil,
				Eq:        nil,
			},
			wantErr: false,
		},
		{
			name: "valid_not_exists",
			fields: fields{
				Exists:    nil,
				NotExists: []string{`any`, ``},
				Regexp:    nil,
				NotRegexp: nil,
				Eq:        nil,
			},
			wantErr: false,
		},
		{
			name: "valid_regex",
			fields: fields{
				Exists:    nil,
				NotExists: nil,
				Regexp: map[string]RequireFieldMatch{
					`key`: {`.`},
				},
				NotRegexp: nil,
				Eq:        nil,
			},
			wantErr: false,
		},
		{
			name: "invalid_regex",
			fields: fields{
				Exists:    nil,
				NotExists: nil,
				Regexp: map[string]RequireFieldMatch{
					`key`: {`(`},
				},
				NotRegexp: nil,
				Eq:        nil,
			},
			wantErr: true,
		},
		{
			name: "valid_not_regex",
			fields: fields{
				Exists:    nil,
				NotExists: nil,
				Regexp:    nil,
				NotRegexp: map[string]RequireFieldMatchNot{
					`key`: {`.`},
				},
				Eq: nil,
			},
			wantErr: false,
		},
		{
			name: "invalid_not_regex",
			fields: fields{
				Exists:    nil,
				NotExists: nil,
				Regexp:    nil,
				NotRegexp: map[string]RequireFieldMatchNot{
					`key`: {`(`},
				},
				Eq: nil,
			},
			wantErr: true,
		},
		{
			name: "valid_eq",
			fields: fields{
				Exists:    nil,
				NotExists: nil,
				Regexp:    nil,
				NotRegexp: nil,
				Eq:        map[string][]string{`key`: {``, `value`}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &HttpArgsRequireHeader{
				Exists:    tt.fields.Exists,
				NotExists: tt.fields.NotExists,
				Regexp:    tt.fields.Regexp,
				NotRegexp: tt.fields.NotRegexp,
				Eq:        tt.fields.Eq,
			}
			if err := a.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
