package args

import "testing"

func TestRequireJSONPath_Match(t *testing.T) {
	type fields struct {
		JSONPath     string
		RequireMatch RequireMatch
	}
	type args struct {
		name    string
		content []byte
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
				JSONPath:     "",
				RequireMatch: RequireMatch{},
			},
			args: args{
				name:    "",
				content: nil,
			},
			wantErr: false,
		},
		{
			name: "blank",
			fields: fields{
				JSONPath:     "",
				RequireMatch: RequireMatch{},
			},
			args: args{
				name:    "",
				content: make([]byte, 0),
			},
			wantErr: false,
		},
		{
			name: "blank_with_path",
			fields: fields{
				JSONPath:     "$",
				RequireMatch: RequireMatch{},
			},
			args: args{
				name:    "",
				content: make([]byte, 0),
			},
			wantErr: true,
		},
		{
			name: "any_root(true)",
			fields: fields{
				JSONPath: "$",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch([]string{`.+`}),
					NotRegexp: nil,
				},
			},
			args: args{
				name:    "",
				content: []byte(`null`),
			},
			wantErr: false,
		},
		{
			name: "root_equal(false)",
			fields: fields{
				JSONPath: "$",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch([]string{`wrong`}),
					NotRegexp: nil,
				},
			},
			args: args{
				name:    "",
				content: []byte(`null`),
			},
			wantErr: true,
		},
		{
			name: "not_any_root",
			fields: fields{
				JSONPath: "$",
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: RequireFieldMatchNot([]string{`.+`}),
				},
			},
			args: args{
				name:    "",
				content: []byte(`null`),
			},
			wantErr: true,
		},
		{
			name: "root_not_equal(true)",
			fields: fields{
				JSONPath: "$",
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: RequireFieldMatchNot([]string{`wrong`}),
				},
			},
			args: args{
				name:    "",
				content: []byte(`null`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &RequireJSONPath{
				JSONPath:     tt.fields.JSONPath,
				RequireMatch: tt.fields.RequireMatch,
			}
			if err := a.Match(tt.args.name, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequireJSONPath_Validate(t *testing.T) {
	type fields struct {
		JSONPath     string
		RequireMatch RequireMatch
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "nil",
			fields: fields{
				JSONPath: "",
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "blank",
			fields: fields{
				JSONPath: "",
				RequireMatch: RequireMatch{
					Regexp:    make(RequireFieldMatch, 0),
					NotRegexp: make(RequireFieldMatchNot, 0),
				},
			},
			wantErr: false,
		},
		{
			name: "match_error",
			fields: fields{
				JSONPath: "",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch([]string{`(`}),
					NotRegexp: make(RequireFieldMatchNot, 0),
				},
			},
			wantErr: true,
		},
		{
			name: "not_match_error",
			fields: fields{
				JSONPath: "",
				RequireMatch: RequireMatch{
					Regexp:    make(RequireFieldMatch, 0),
					NotRegexp: RequireFieldMatchNot([]string{`(`}),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &RequireJSONPath{
				JSONPath:     tt.fields.JSONPath,
				RequireMatch: tt.fields.RequireMatch,
			}
			if err := a.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
