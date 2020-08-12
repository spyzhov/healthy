package internal

import "testing"

func TestHttpArgsRequireContentJSON_Match(t *testing.T) {
	type fields struct {
		JSONPath     string
		RequireXPath RequireXPath
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
				JSONPath: "",
				RequireXPath: RequireXPath{
					XPath: "",
					RequireMatch: RequireMatch{
						Regexp:    nil,
						NotRegexp: nil,
					},
				},
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
				JSONPath: "",
				RequireXPath: RequireXPath{
					XPath: "",
					RequireMatch: RequireMatch{
						Regexp:    make(RequireFieldMatch, 0),
						NotRegexp: make(RequireFieldMatchNot, 0),
					},
				},
			},
			args: args{
				name:    "",
				content: make([]byte, 0),
			},
			wantErr: false,
		},
		{
			name: "valid_jsonpath_match",
			fields: fields{
				JSONPath: "$.key",
				RequireXPath: RequireXPath{
					XPath: "",
					RequireMatch: RequireMatch{
						Regexp: RequireFieldMatch{
							".+value.+",
						},
						NotRegexp: nil,
					},
				},
			},
			args: args{
				name:    "",
				content: []byte(`{"key": "value"}`),
			},
			wantErr: false,
		},
		{
			name: "invalid_jsonpath_match",
			fields: fields{
				JSONPath: "$.key",
				RequireXPath: RequireXPath{
					XPath: "",
					RequireMatch: RequireMatch{
						Regexp: RequireFieldMatch{
							".+wrong.+",
						},
						NotRegexp: nil,
					},
				},
			},
			args: args{
				name:    "",
				content: []byte(`{"key": "value"}`),
			},
			wantErr: true,
		},
		{
			name: "valid_jsonpath_not_match",
			fields: fields{
				JSONPath: "$.key",
				RequireXPath: RequireXPath{
					XPath: "",
					RequireMatch: RequireMatch{
						Regexp: nil,
						NotRegexp: RequireFieldMatchNot{
							".+wrong.+",
						},
					},
				},
			},
			args: args{
				name:    "",
				content: []byte(`{"key": "value"}`),
			},
			wantErr: false,
		},
		{
			name: "invalid_jsonpath_not_match",
			fields: fields{
				JSONPath: "$.key",
				RequireXPath: RequireXPath{
					XPath: "",
					RequireMatch: RequireMatch{
						Regexp: nil,
						NotRegexp: RequireFieldMatchNot{
							".+value.+",
						},
					},
				},
			},
			args: args{
				name:    "",
				content: []byte(`{"key": "value"}`),
			},
			wantErr: true,
		},
		{
			name: "valid_xpath_match",
			fields: fields{
				JSONPath: "",
				RequireXPath: RequireXPath{
					XPath: "//key",
					RequireMatch: RequireMatch{
						Regexp: RequireFieldMatch{
							"value.*",
						},
						NotRegexp: nil,
					},
				},
			},
			args: args{
				name:    "",
				content: []byte(`{"key": "value"}`),
			},
			wantErr: false,
		},
		{
			name: "invalid_xpath_match",
			fields: fields{
				JSONPath: "",
				RequireXPath: RequireXPath{
					XPath: "//key",
					RequireMatch: RequireMatch{
						Regexp: RequireFieldMatch{
							".+wrong.+",
						},
						NotRegexp: nil,
					},
				},
			},
			args: args{
				name:    "",
				content: []byte(`{"key": "value"}`),
			},
			wantErr: true,
		},
		{
			name: "valid_xpath_not_match",
			fields: fields{
				JSONPath: "",
				RequireXPath: RequireXPath{
					XPath: "//key",
					RequireMatch: RequireMatch{
						Regexp: nil,
						NotRegexp: RequireFieldMatchNot{
							".+wrong.+",
						},
					},
				},
			},
			args: args{
				name:    "",
				content: []byte(`{"key": "value"}`),
			},
			wantErr: false,
		},
		{
			name: "invalid_xpath_not_match",
			fields: fields{
				JSONPath: "",
				RequireXPath: RequireXPath{
					XPath: "//key",
					RequireMatch: RequireMatch{
						Regexp: nil,
						NotRegexp: RequireFieldMatchNot{
							"value.*",
						},
					},
				},
			},
			args: args{
				name:    "",
				content: []byte(`{"key": "value"}`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := HttpArgsRequireContentJSON{
				JSONPath:     tt.fields.JSONPath,
				RequireXPath: tt.fields.RequireXPath,
			}
			if err := a.Match(tt.args.name, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttpArgsRequireContentJSON_Validate(t *testing.T) {
	type fields struct {
		JSONPath     string
		RequireXPath RequireXPath
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
				RequireXPath: RequireXPath{
					XPath: "",
					RequireMatch: RequireMatch{
						Regexp:    nil,
						NotRegexp: nil,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "blank",
			fields: fields{
				JSONPath: "",
				RequireXPath: RequireXPath{
					XPath: "",
					RequireMatch: RequireMatch{
						Regexp:    make(RequireFieldMatch, 0),
						NotRegexp: make(RequireFieldMatchNot, 0),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid",
			fields: fields{
				JSONPath: "$",
				RequireXPath: RequireXPath{
					XPath: "//path",
					RequireMatch: RequireMatch{
						Regexp: RequireFieldMatch{
							`.+`,
						},
						NotRegexp: RequireFieldMatchNot{
							`.+`,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid_1",
			fields: fields{
				JSONPath: "$",
				RequireXPath: RequireXPath{
					XPath: "//path",
					RequireMatch: RequireMatch{
						Regexp: RequireFieldMatch{
							`(`,
						},
						NotRegexp: RequireFieldMatchNot{
							`.+`,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid_2",
			fields: fields{
				JSONPath: "$",
				RequireXPath: RequireXPath{
					XPath: "//path",
					RequireMatch: RequireMatch{
						Regexp: RequireFieldMatch{
							`.+`,
						},
						NotRegexp: RequireFieldMatchNot{
							`(`,
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := HttpArgsRequireContentJSON{
				JSONPath:     tt.fields.JSONPath,
				RequireXPath: tt.fields.RequireXPath,
			}
			if err := a.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
