package args

import "testing"

func TestRequireMatch_Match(t *testing.T) {
	type fields struct {
		Regexp    RequireFieldMatch
		NotRegexp RequireFieldMatchNot
	}
	type args struct {
		name  string
		input []byte
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
				Regexp:    nil,
				NotRegexp: nil,
			},
			args: args{
				name:  "",
				input: nil,
			},
			wantErr: false,
		},
		{
			name: "blank",
			fields: fields{
				Regexp:    make(RequireFieldMatch, 0),
				NotRegexp: make(RequireFieldMatchNot, 0),
			},
			args: args{
				name:  "",
				input: make([]byte, 0),
			},
			wantErr: false,
		},
		{
			name: "match_any(true)",
			fields: fields{
				Regexp:    RequireFieldMatch([]string{`.+`}),
				NotRegexp: make(RequireFieldMatchNot, 0),
			},
			args: args{
				name:  "",
				input: []byte(`string`),
			},
			wantErr: false,
		},
		{
			name: "match_any(false)",
			fields: fields{
				Regexp:    RequireFieldMatch([]string{`.+`}),
				NotRegexp: make(RequireFieldMatchNot, 0),
			},
			args: args{
				name:  "",
				input: []byte(``),
			},
			wantErr: true,
		},
		{
			name: "not_match_any(false)",
			fields: fields{
				Regexp:    make(RequireFieldMatch, 0),
				NotRegexp: RequireFieldMatchNot([]string{`.+`}),
			},
			args: args{
				name:  "",
				input: []byte(`string`),
			},
			wantErr: true,
		},
		{
			name: "not_match(true)",
			fields: fields{
				Regexp:    make(RequireFieldMatch, 0),
				NotRegexp: RequireFieldMatchNot([]string{`string`}),
			},
			args: args{
				name:  "",
				input: []byte(`wrong`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &RequireMatch{
				Regexp:    tt.fields.Regexp,
				NotRegexp: tt.fields.NotRegexp,
			}
			if err := a.Match(tt.args.name, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequireMatch_MatchStrings(t *testing.T) {
	type fields struct {
		Regexp    RequireFieldMatch
		NotRegexp RequireFieldMatchNot
	}
	type args struct {
		name  string
		input []string
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
				Regexp:    nil,
				NotRegexp: nil,
			},
			args: args{
				name:  "",
				input: nil,
			},
			wantErr: false,
		},
		{
			name: "blank",
			fields: fields{
				Regexp:    make(RequireFieldMatch, 0),
				NotRegexp: make(RequireFieldMatchNot, 0),
			},
			args: args{
				name:  "",
				input: make([]string, 0),
			},
			wantErr: false,
		},
		{
			name: "match_any(true)",
			fields: fields{
				Regexp:    RequireFieldMatch([]string{`.+`}),
				NotRegexp: make(RequireFieldMatchNot, 0),
			},
			args: args{
				name:  "",
				input: []string{`string`},
			},
			wantErr: false,
		},
		{
			name: "match_any(false)",
			fields: fields{
				Regexp:    RequireFieldMatch([]string{`.+`}),
				NotRegexp: make(RequireFieldMatchNot, 0),
			},
			args: args{
				name:  "",
				input: []string{""},
			},
			wantErr: true,
		},
		{
			name: "not_match_any(false)",
			fields: fields{
				Regexp:    make(RequireFieldMatch, 0),
				NotRegexp: RequireFieldMatchNot([]string{`.+`}),
			},
			args: args{
				name:  "",
				input: []string{`string`},
			},
			wantErr: true,
		},
		{
			name: "not_match(true)",
			fields: fields{
				Regexp:    make(RequireFieldMatch, 0),
				NotRegexp: RequireFieldMatchNot([]string{`string`}),
			},
			args: args{
				name:  "",
				input: []string{`wrong`},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &RequireMatch{
				Regexp:    tt.fields.Regexp,
				NotRegexp: tt.fields.NotRegexp,
			}
			if err := a.MatchStrings(tt.args.name, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("MatchStrings() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequireMatch_Validate(t *testing.T) {
	type fields struct {
		Regexp    RequireFieldMatch
		NotRegexp RequireFieldMatchNot
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "nil",
			fields: fields{
				Regexp:    nil,
				NotRegexp: nil,
			},
			wantErr: false,
		},
		{
			name: "blank",
			fields: fields{
				Regexp:    RequireFieldMatch{},
				NotRegexp: RequireFieldMatchNot{},
			},
			wantErr: false,
		},
		{
			name: "valid",
			fields: fields{
				Regexp:    RequireFieldMatch{`.+`},
				NotRegexp: RequireFieldMatchNot{`.+`},
			},
			wantErr: false,
		},
		{
			name: "not_valid_match",
			fields: fields{
				Regexp:    RequireFieldMatch{`(`},
				NotRegexp: RequireFieldMatchNot{},
			},
			wantErr: true,
		},
		{
			name: "not_valid_not_match",
			fields: fields{
				Regexp:    RequireFieldMatch{},
				NotRegexp: RequireFieldMatchNot{`(`},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &RequireMatch{
				Regexp:    tt.fields.Regexp,
				NotRegexp: tt.fields.NotRegexp,
			}
			if err := a.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
