package internal

import "testing"

func TestHttpArgsRequireContentType_Is(t *testing.T) {
	type args struct {
		value HttpArgsRequireContentType
	}
	tests := []struct {
		name string
		v    HttpArgsRequireContentType
		args args
		want bool
	}{
		{
			name: "blank",
			v:    "",
			args: args{
				value: "",
			},
			want: true,
		},
		{
			name: "eq",
			v:    HttpArgsRequireContentTypeHTML,
			args: args{
				value: HttpArgsRequireContentTypeHTML,
			},
			want: true,
		},
		{
			name: "neq",
			v:    HttpArgsRequireContentTypeHTML,
			args: args{
				value: HttpArgsRequireContentTypeJSON,
			},
			want: false,
		},
		{
			name: "one_blank2",
			v:    "",
			args: args{
				value: HttpArgsRequireContentTypeHTML,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Is(tt.args.value); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpArgsRequireContentType_Match(t *testing.T) {
	type args struct {
		content []byte
	}
	tests := []struct {
		name    string
		v       HttpArgsRequireContentType
		args    args
		wantErr bool
	}{
		{
			name: "nil",
			v:    "",
			args: args{
				content: nil,
			},
			wantErr: false,
		},
		{
			name: "blank",
			v:    "",
			args: args{
				content: make([]byte, 0),
			},
			wantErr: false,
		},
		{
			name: "valid_xml",
			v:    HttpArgsRequireContentTypeXML,
			args: args{
				content: []byte(`<?xml version="1.0" encoding="UTF-8" ?>
		<key>value</key>`),
			},
			wantErr: false,
		},
		{
			name: "invalid_xml",
			v:    HttpArgsRequireContentTypeXML,
			args: args{
				content: []byte(`<?xml version="1.1" encoding="UTF-8" ?>
		<key>value</key>`),
			},
			wantErr: true,
		},
		{
			name: "valid_html",
			v:    HttpArgsRequireContentTypeHTML,
			args: args{
				content: []byte(`<!doctype html>
		<html lang="en"><head></head><body><key>value</key></body></html>`),
			},
			wantErr: false,
		},
		{
			name: "invalid_html",
			v:    HttpArgsRequireContentTypeHTML,
			args: args{
				content: []byte(`<foo><bar></foo></baz>`),
			},
			wantErr: false, // can't find invalid example
		},
		{
			name: "valid_json",
			v:    HttpArgsRequireContentTypeJSON,
			args: args{
				content: []byte(`{"valid": true}`),
			},
			wantErr: false,
		},
		{
			name: "invalid_json",
			v:    HttpArgsRequireContentTypeJSON,
			args: args{
				content: []byte(`{"'`),
			},
			wantErr: true,
		},
		{
			name: "valid_yaml",
			v:    HttpArgsRequireContentTypeYAML,
			args: args{
				content: []byte(`---
valid: yaml
`),
			},
			wantErr: false,
		},
		{
			name: "invalid_yaml",
			v:    HttpArgsRequireContentTypeYAML,
			args: args{
				content: []byte(`{"'`),
			},
			wantErr: true,
		},
		{
			name: "invalid_type",
			v:    "wrong",
			args: args{
				content: []byte(``),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.Match(tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttpArgsRequireContentType_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       HttpArgsRequireContentType
		wantErr bool
	}{
		{
			name:    "blank",
			v:       "",
			wantErr: false,
		},
		{
			name:    "valid_lower_html",
			v:       "html",
			wantErr: false,
		},
		{
			name:    "valid_upper_json",
			v:       "JSON",
			wantErr: false,
		},
		{
			name:    "valid_yaml",
			v:       "YAML",
			wantErr: false,
		},
		{
			name:    "valid_xml",
			v:       HttpArgsRequireContentTypeXML,
			wantErr: false,
		},
		{
			name:    "invalid",
			v:       "invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
