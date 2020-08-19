package args

import (
	"testing"
)

func TestRequireXPath_Match(t *testing.T) {
	type fields struct {
		XPath        string
		RequireMatch RequireMatch
	}
	type args struct {
		name    string
		_type   HttpArgsRequireContentType
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
				XPath: "",
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
			},
			args: args{
				name:    "",
				_type:   "",
				content: nil,
			},
			wantErr: false,
		},
		{
			name: "blank",
			fields: fields{
				XPath: "",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{},
					NotRegexp: RequireFieldMatchNot{},
				},
			},
			args: args{
				name:    "",
				_type:   "",
				content: make([]byte, 0),
			},
			wantErr: false,
		},
		{
			name: "match_json(true)",
			fields: fields{
				XPath: "//key",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{`(value)`},
					NotRegexp: RequireFieldMatchNot{},
				},
			},
			args: args{
				name:    "",
				_type:   "JSON",
				content: []byte(`{"key":"value"}`),
			},
			wantErr: false,
		},
		{
			name: "match_json(false)",
			fields: fields{
				XPath: "//key",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{`(wrong)`},
					NotRegexp: RequireFieldMatchNot{},
				},
			},
			args: args{
				name:    "",
				_type:   "JSON",
				content: []byte(`{"key":"value"}`),
			},
			wantErr: true,
		},
		{
			name: "not_match_json(true)",
			fields: fields{
				XPath: "//key",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{},
					NotRegexp: RequireFieldMatchNot{`wrong`},
				},
			},
			args: args{
				name:    "",
				_type:   "JSON",
				content: []byte(`{"key":"value"}`),
			},
			wantErr: false,
		},
		{
			name: "not_match_json(false)",
			fields: fields{
				XPath: "//key",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{},
					NotRegexp: RequireFieldMatchNot{`(value)`},
				},
			},
			args: args{
				name:    "",
				_type:   "JSON",
				content: []byte(`{"key":"value"}`),
			},
			wantErr: true,
		},
		{
			name: "match_html(true)",
			fields: fields{
				XPath: "//key",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{`(value)`},
					NotRegexp: RequireFieldMatchNot{},
				},
			},
			args: args{
				name:  "",
				_type: "HTML",
				content: []byte(`<!doctype html>
<html lang="en"><head></head><body><key>value</key></body></html>`),
			},
			wantErr: false,
		},
		{
			name: "match_html(false)",
			fields: fields{
				XPath: "//key",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{`(wrong)`},
					NotRegexp: RequireFieldMatchNot{},
				},
			},
			args: args{
				name:  "",
				_type: "HTML",
				content: []byte(`<!doctype html>
<html lang="en"><head></head><body><key>value</key></body></html>`),
			},
			wantErr: true,
		},
		{
			name: "not_match_html(true)",
			fields: fields{
				XPath: "//key",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{},
					NotRegexp: RequireFieldMatchNot{`wrong`},
				},
			},
			args: args{
				name:  "",
				_type: "HTML",
				content: []byte(`<!doctype html>
<html lang="en"><head></head><body><key>value</key></body></html>`),
			},
			wantErr: false,
		},
		{
			name: "not_match_html(false)",
			fields: fields{
				XPath: "//key",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{},
					NotRegexp: RequireFieldMatchNot{`(value)`},
				},
			},
			args: args{
				name:  "",
				_type: "HTML",
				content: []byte(`<!doctype html>
<html lang="en"><head></head><body><key>value</key></body></html>`),
			},
			wantErr: true,
		},
		{
			name: "match_xml(true)",
			fields: fields{
				XPath: "//key",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{`(value)`},
					NotRegexp: RequireFieldMatchNot{},
				},
			},
			args: args{
				name:  "",
				_type: "XML",
				content: []byte(`<?xml version="1.0" encoding="UTF-8" ?>
<key>value</key>`),
			},
			wantErr: false,
		},
		{
			name: "match_xml(false)",
			fields: fields{
				XPath: "//key",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{`(wrong)`},
					NotRegexp: RequireFieldMatchNot{},
				},
			},
			args: args{
				name:  "",
				_type: "XML",
				content: []byte(`<?xml version="1.0" encoding="UTF-8" ?>
<key>value</key>`),
			},
			wantErr: true,
		},
		{
			name: "not_match_xml(true)",
			fields: fields{
				XPath: "//key",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{},
					NotRegexp: RequireFieldMatchNot{`wrong`},
				},
			},
			args: args{
				name:  "",
				_type: "XML",
				content: []byte(`<?xml version="1.0" encoding="UTF-8" ?>
<key>value</key>`),
			},
			wantErr: false,
		},
		{
			name: "not_match_xml(false)",
			fields: fields{
				XPath: "//key",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{},
					NotRegexp: RequireFieldMatchNot{`(value)`},
				},
			},
			args: args{
				name:  "",
				_type: "XML",
				content: []byte(`<?xml version="1.0" encoding="UTF-8" ?>
<key>value</key>`),
			},
			wantErr: true,
		},
		{
			name: "wrong_type(false)",
			fields: fields{
				XPath: "//key",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{},
					NotRegexp: RequireFieldMatchNot{},
				},
			},
			args: args{
				name:  "",
				_type: "WRONG",
				content: []byte(`<?xml version="1.0" encoding="UTF-8" ?>
<key>value</key>`),
			},
			wantErr: true,
		},
		{
			name: "xpath_error(false)",
			fields: fields{
				XPath: "//key",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{},
					NotRegexp: RequireFieldMatchNot{},
				},
			},
			args: args{
				name:  "",
				_type: "XML",
				content: []byte(`<?xml version="1.1" encoding="UTF-8" ?>
<key>value</key>`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &RequireXPath{
				XPath:        tt.fields.XPath,
				RequireMatch: tt.fields.RequireMatch,
			}
			if err := a.Match(tt.args.name, tt.args._type, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequireXPath_Validate(t *testing.T) {
	type fields struct {
		XPath        string
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
				XPath: "",
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
				XPath: "",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{},
					NotRegexp: RequireFieldMatchNot{},
				},
			},
			wantErr: false,
		},
		{
			name: "match(true)",
			fields: fields{
				XPath: "",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{`.+`},
					NotRegexp: RequireFieldMatchNot{},
				},
			},
			wantErr: false,
		},
		{
			name: "not_match(true)",
			fields: fields{
				XPath: "",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{},
					NotRegexp: RequireFieldMatchNot{`.+`},
				},
			},
			wantErr: false,
		},
		{
			name: "match(false)",
			fields: fields{
				XPath: "",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{},
					NotRegexp: RequireFieldMatchNot{`(`},
				},
			},
			wantErr: true,
		},
		{
			name: "not_match(false)",
			fields: fields{
				XPath: "",
				RequireMatch: RequireMatch{
					Regexp:    RequireFieldMatch{},
					NotRegexp: RequireFieldMatchNot{`(`},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &RequireXPath{
				XPath:        tt.fields.XPath,
				RequireMatch: tt.fields.RequireMatch,
			}
			if err := a.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
