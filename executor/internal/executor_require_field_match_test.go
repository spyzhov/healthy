package internal

import (
	"testing"
)

func TestRequireFieldMatch_Match(t *testing.T) {
	type args struct {
		name  string
		input []byte
	}
	tests := []struct {
		name    string
		match   RequireFieldMatch
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
			match: make(RequireFieldMatch, 0),
			args: args{
				name:  "",
				input: make([]byte, 0),
			},
			wantErr: false,
		},
		{
			name: "match(.*)",
			match: RequireFieldMatch([]string{
				".*",
			}),
			args: args{
				name:  "",
				input: []byte(``),
			},
			wantErr: false,
		},
		{
			name: "match(.+)",
			match: RequireFieldMatch([]string{
				".+",
			}),
			args: args{
				name:  "",
				input: []byte(`string`),
			},
			wantErr: false,
		},
		{
			name: "match(string)",
			match: RequireFieldMatch([]string{
				"string",
			}),
			args: args{
				name:  "",
				input: []byte(`string`),
			},
			wantErr: false,
		},
		{
			name: "match((string))",
			match: RequireFieldMatch([]string{
				"(string)",
			}),
			args: args{
				name:  "",
				input: []byte(`long string`),
			},
			wantErr: false,
		},
		{
			name: "match((another) and (string))",
			match: RequireFieldMatch([]string{
				"(another)",
				"(string)",
			}),
			args: args{
				name:  "",
				input: []byte(`another string`),
			},
			wantErr: false,
		},
		{
			name: "not_match()",
			match: RequireFieldMatch([]string{
				"not_match",
			}),
			args: args{
				name:  "",
				input: []byte(`not match`),
			},
			wantErr: true,
		},
		{
			name: "not_match(one_of)",
			match: RequireFieldMatch([]string{
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

func TestRequireFieldMatch_MatchStrings(t *testing.T) {
	type args struct {
		name  string
		input []string
	}
	tests := []struct {
		name    string
		match   RequireFieldMatch
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
			match: make(RequireFieldMatch, 0),
			args: args{
				name:  "",
				input: make([]string, 0),
			},
			wantErr: false,
		},
		{
			name: "match(.*)",
			match: RequireFieldMatch([]string{
				".*",
			}),
			args: args{
				name:  "",
				input: []string{""},
			},
			wantErr: false,
		},
		{
			name: "match(.+)",
			match: RequireFieldMatch([]string{
				".+",
			}),
			args: args{
				name:  "",
				input: []string{`string`},
			},
			wantErr: false,
		},
		{
			name: "match(string)",
			match: RequireFieldMatch([]string{
				"string",
			}),
			args: args{
				name:  "",
				input: []string{`string`},
			},
			wantErr: false,
		},
		{
			name: "match((string))",
			match: RequireFieldMatch([]string{
				"(string)",
			}),
			args: args{
				name:  "",
				input: []string{`long string`},
			},
			wantErr: false,
		},
		{
			name: "match((another) and (string))",
			match: RequireFieldMatch([]string{
				"(another)",
				"(string)",
			}),
			args: args{
				name:  "",
				input: []string{`another string`},
			},
			wantErr: false,
		},
		{
			name: "not_match()",
			match: RequireFieldMatch([]string{
				"not_match",
			}),
			args: args{
				name:  "",
				input: []string{`not match`},
			},
			wantErr: true,
		},
		{
			name: "not_match(one_of)",
			match: RequireFieldMatch([]string{
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
			match: RequireFieldMatch([]string{
				"match",
				"not_match",
			}),
			args: args{
				name:  "",
				input: []string{`match`, `not_match`},
			},
			wantErr: false,
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

func TestRequireFieldMatch_Validate(t *testing.T) {
	tests := []struct {
		name    string
		match   RequireFieldMatch
		wantErr bool
	}{
		{
			name:    "nil",
			match:   nil,
			wantErr: false,
		},
		{
			name:    "blank",
			match:   make(RequireFieldMatch, 0),
			wantErr: false,
		},
		{
			name: "valid",
			match: RequireFieldMatch([]string{
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
			match: RequireFieldMatch([]string{
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
