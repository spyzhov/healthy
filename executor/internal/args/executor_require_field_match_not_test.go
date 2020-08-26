package args

import "testing"

func TestRequireFieldMatchNot_Match(t *testing.T) {
	type args struct {
		name  string
		input []byte
	}
	tests := []struct {
		name    string
		match   RequireFieldMatchNot
		args    args
		wantErr bool
	}{
		{
			name:  "nil",
			match: nil,
			args: args{
				name:  "",
				input: nil,
			},
			wantErr: false,
		},
		{
			name:  "blank",
			match: make(RequireFieldMatchNot, 0),
			args: args{
				name:  "",
				input: make([]byte, 0),
			},
			wantErr: false,
		},
		{
			name: "match(.*)",
			match: RequireFieldMatchNot([]string{
				".*",
			}),
			args: args{
				name:  "",
				input: []byte(``),
			},
			wantErr: true,
		},
		{
			name: "match(.+)",
			match: RequireFieldMatchNot([]string{
				".+",
			}),
			args: args{
				name:  "",
				input: []byte(`string`),
			},
			wantErr: true,
		},
		{
			name: "match(string)",
			match: RequireFieldMatchNot([]string{
				"string",
			}),
			args: args{
				name:  "",
				input: []byte(`string`),
			},
			wantErr: true,
		},
		{
			name: "match((string))",
			match: RequireFieldMatchNot([]string{
				"(string)",
			}),
			args: args{
				name:  "",
				input: []byte(`long string`),
			},
			wantErr: true,
		},
		{
			name: "match((another) and (string))",
			match: RequireFieldMatchNot([]string{
				"(another)",
				"(string)",
			}),
			args: args{
				name:  "",
				input: []byte(`another string`),
			},
			wantErr: true,
		},
		{
			name: "not_match()",
			match: RequireFieldMatchNot([]string{
				"not_match",
			}),
			args: args{
				name:  "",
				input: []byte(`not match`),
			},
			wantErr: false,
		},
		{
			name: "not_match(one_of)",
			match: RequireFieldMatchNot([]string{
				"match",
				"not_match",
			}),
			args: args{
				name:  "",
				input: []byte(`match`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.match.Match(tt.args.name, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequireFieldMatchNot_MatchStrings(t *testing.T) {
	type args struct {
		name  string
		input []string
	}
	tests := []struct {
		name    string
		match   RequireFieldMatchNot
		args    args
		wantErr bool
	}{
		{
			name:  "nil",
			match: nil,
			args: args{
				name:  "",
				input: nil,
			},
			wantErr: false,
		},
		{
			name:  "blank",
			match: make(RequireFieldMatchNot, 0),
			args: args{
				name:  "",
				input: make([]string, 0),
			},
			wantErr: false,
		},
		{
			name: "match(.*)",
			match: RequireFieldMatchNot([]string{
				".*",
			}),
			args: args{
				name:  "",
				input: []string{""},
			},
			wantErr: true,
		},
		{
			name: "match(.+)",
			match: RequireFieldMatchNot([]string{
				".+",
			}),
			args: args{
				name:  "",
				input: []string{`string`},
			},
			wantErr: true,
		},
		{
			name: "match(string)",
			match: RequireFieldMatchNot([]string{
				"string",
			}),
			args: args{
				name:  "",
				input: []string{`string`},
			},
			wantErr: true,
		},
		{
			name: "match((string))",
			match: RequireFieldMatchNot([]string{
				"(string)",
			}),
			args: args{
				name:  "",
				input: []string{`long string`},
			},
			wantErr: true,
		},
		{
			name: "match((another) and (string))",
			match: RequireFieldMatchNot([]string{
				"(another)",
				"(string)",
			}),
			args: args{
				name:  "",
				input: []string{`another string`},
			},
			wantErr: true,
		},
		{
			name: "not_match()",
			match: RequireFieldMatchNot([]string{
				"not_match",
			}),
			args: args{
				name:  "",
				input: []string{`not match`},
			},
			wantErr: false,
		},
		{
			name: "not_match(one_of)",
			match: RequireFieldMatchNot([]string{
				"match",
				"not_match",
			}),
			args: args{
				name:  "",
				input: []string{`match`},
			},
			wantErr: true,
		},
		{
			name: "match(one by one)",
			match: RequireFieldMatchNot([]string{
				"match",
				"not_match",
			}),
			args: args{
				name:  "",
				input: []string{`match`, `not_match`},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.match.MatchStrings(tt.args.name, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("MatchStrings() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequireFieldMatchNot_Validate(t *testing.T) {
	tests := []struct {
		name    string
		match   RequireFieldMatchNot
		wantErr bool
	}{
		{
			name:    "nil",
			match:   nil,
			wantErr: false,
		},
		{
			name:    "blank",
			match:   make(RequireFieldMatchNot, 0),
			wantErr: false,
		},
		{
			name: "valid",
			match: RequireFieldMatchNot([]string{
				`.*`,
				`.+`,
				`string`,
				`(string)`,
				`(?m)(string)`,
			}),
			wantErr: false,
		},
		{
			name: "not_valid",
			match: RequireFieldMatchNot([]string{
				`.*`,
				`(.+`,
			}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.match.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
