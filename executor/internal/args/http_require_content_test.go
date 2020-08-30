package args

import "testing"

func TestHttpArgsRequireContent_Match(t *testing.T) {
	type fields struct {
		RequireMatch RequireMatch
		Type         HttpArgsRequireContentType
		Length       *RequireNumeric
		JSON         []RequireJSON
		XML          []RequireXPath
		HTML         []RequireXPath
	}
	type args struct {
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
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
			},
			args: args{
				content: nil,
			},
			wantErr: false,
		},
		{
			name: "blank",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    make(RequireFieldMatch, 0),
					NotRegexp: make(RequireFieldMatchNot, 0),
				},
				Type: "",
				Length: &RequireNumeric{
					In:    make([]float64, 0),
					NotIn: make([]float64, 0),
					Eq:    nil,
					Le:    nil,
					Leq:   nil,
					Ge:    nil,
					Geq:   nil,
					Not:   nil,
				},
				JSON: make([]RequireJSON, 0),
				XML:  make([]RequireXPath, 0),
				HTML: make([]RequireXPath, 0),
			},
			args: args{
				content: make([]byte, 0),
			},
			wantErr: false,
		},
		{
			name: "valid_match",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp: RequireFieldMatch{
						`content`,
					},
					NotRegexp: nil,
				},
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
			},
			args: args{
				content: []byte(`content`),
			},
			wantErr: false,
		},
		{
			name: "invalid_match",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp: RequireFieldMatch{
						`wrong`,
					},
					NotRegexp: nil,
				},
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
			},
			args: args{
				content: []byte(`content`),
			},
			wantErr: true,
		},
		{
			name: "valid_not_match",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp: nil,
					NotRegexp: RequireFieldMatchNot{
						`wrong`,
					},
				},
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
			},
			args: args{
				content: []byte(`content`),
			},
			wantErr: false,
		},
		{
			name: "invalid_not_match",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp: nil,
					NotRegexp: RequireFieldMatchNot{
						`content`,
					},
				},
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
			},
			args: args{
				content: []byte(`content`),
			},
			wantErr: true,
		},
		{
			name: "valid_type",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type:   "JSON",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
			},
			args: args{
				content: []byte(`{"foo": true}`),
			},
			wantErr: false,
		},
		{
			name: "invalid_type",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type:   "JSON",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
			},
			args: args{
				content: []byte(`{"foo"`),
			},
			wantErr: true,
		},
		{
			name: "valid_length",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type: "",
				Length: &RequireNumeric{
					In:    []float64{7},
					NotIn: nil,
					Eq:    nil,
					Le:    nil,
					Leq:   nil,
					Ge:    nil,
					Geq:   nil,
					Not:   nil,
				},
				JSON: nil,
				XML:  nil,
				HTML: nil,
			},
			args: args{
				content: []byte(`content`),
			},
			wantErr: false,
		},
		{
			name: "invalid_length",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type: "",
				Length: &RequireNumeric{
					In:    []float64{1},
					NotIn: nil,
					Eq:    nil,
					Le:    nil,
					Leq:   nil,
					Ge:    nil,
					Geq:   nil,
					Not:   nil,
				},
				JSON: nil,
				XML:  nil,
				HTML: nil,
			},
			args: args{
				content: []byte(`content`),
			},
			wantErr: true,
		},
		{
			name: "valid_json",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type:   "",
				Length: nil,
				JSON: []RequireJSON{
					{
						JSONPath: "$",
						RequireXPath: RequireXPath{
							XPath: "",
							RequireMatch: RequireMatch{
								Regexp: RequireFieldMatch{
									`.+`,
								},
								NotRegexp: nil,
							},
						},
					},
				},
				XML:  nil,
				HTML: nil,
			},
			args: args{
				content: []byte(`{}`),
			},
			wantErr: false,
		},
		{
			name: "invalid_json",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type:   "",
				Length: nil,
				JSON: []RequireJSON{
					{
						JSONPath: "$",
						RequireXPath: RequireXPath{
							XPath: "",
							RequireMatch: RequireMatch{
								Regexp: nil,
								NotRegexp: RequireFieldMatchNot{
									`.+`,
								},
							},
						},
					},
				},
				XML:  nil,
				HTML: nil,
			},
			args: args{
				content: []byte(`{}`),
			},
			wantErr: true,
		},
		{
			name: "valid_xml",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML: []RequireXPath{
					{
						XPath: "//key",
						RequireMatch: RequireMatch{
							Regexp: RequireFieldMatch{
								`.+`,
							},
							NotRegexp: nil,
						},
					},
				},
				HTML: nil,
			},
			args: args{
				content: []byte(`<?xml version="1.0" encoding="UTF-8" ?>
<key>value</key>`),
			},
			wantErr: false,
		},
		{
			name: "invalid_xml",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML: []RequireXPath{
					{
						XPath: "//key",
						RequireMatch: RequireMatch{
							Regexp: nil,
							NotRegexp: RequireFieldMatchNot{
								`.+`,
							},
						},
					},
				},
				HTML: nil,
			},
			args: args{
				content: []byte(`<?xml version="1.0" encoding="UTF-8" ?>
<key>value</key>`),
			},
			wantErr: true,
		},
		{
			name: "valid_html",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML: []RequireXPath{
					{
						XPath: "//key",
						RequireMatch: RequireMatch{
							Regexp: RequireFieldMatch{
								`.+`,
							},
							NotRegexp: nil,
						},
					},
				},
			},
			args: args{
				content: []byte(`<!doctype html>
<html lang="en"><head></head><body><key>value</key></body></html>`),
			},
			wantErr: false,
		},
		{
			name: "invalid_html",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML: []RequireXPath{
					{
						XPath: "//key",
						RequireMatch: RequireMatch{
							Regexp: nil,
							NotRegexp: RequireFieldMatchNot{
								`.+`,
							},
						},
					},
				},
			},
			args: args{
				content: []byte(`<!doctype html>
<html lang="en"><head></head><body><key>value</key></body></html>`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &HttpArgsRequireContent{
				RequireMatch: tt.fields.RequireMatch,
				Type:         tt.fields.Type,
				Length:       tt.fields.Length,
				JSON:         tt.fields.JSON,
				XML:          tt.fields.XML,
				HTML:         tt.fields.HTML,
			}
			if err := a.Match(tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttpArgsRequireContent_Validate(t *testing.T) {
	type fields struct {
		RequireMatch RequireMatch
		Type         HttpArgsRequireContentType
		Length       *RequireNumeric
		JSON         []RequireJSON
		XML          []RequireXPath
		HTML         []RequireXPath
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "nil",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
			},
			wantErr: false,
		},
		{
			name: "blank",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    make(RequireFieldMatch, 0),
					NotRegexp: make(RequireFieldMatchNot, 0),
				},
				Type: "",
				Length: &RequireNumeric{
					In:    make([]float64, 0),
					NotIn: make([]float64, 0),
					Eq:    nil,
					Le:    nil,
					Leq:   nil,
					Ge:    nil,
					Geq:   nil,
					Not:   nil,
				},
				JSON: make([]RequireJSON, 0),
				XML:  make([]RequireXPath, 0),
				HTML: make([]RequireXPath, 0),
			},
			wantErr: false,
		},
		{
			name: "valid_RequireMatch_Regexp",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp: RequireFieldMatch{
						`.`,
					},
					NotRegexp: nil,
				},
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
			},
			wantErr: false,
		},
		{
			name: "invalid_RequireMatch_Regexp",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp: RequireFieldMatch{
						`(`,
					},
					NotRegexp: nil,
				},
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
			},
			wantErr: true,
		},
		{
			name: "valid_RequireMatch_NotRegexp",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp: nil,
					NotRegexp: RequireFieldMatchNot{
						`.`,
					},
				},
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
			},
			wantErr: false,
		},
		{
			name: "invalid_RequireMatch_NotRegexp",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp: nil,
					NotRegexp: RequireFieldMatchNot{
						`(`,
					},
				},
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
			},
			wantErr: true,
		},
		{
			name: "valid_Type",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type:   "JSON",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
			},
			wantErr: false,
		},
		{
			name: "invalid_Type",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type:   "WRONG",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
			},
			wantErr: true,
		},
		{
			name: "valid_JSON",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type:   "",
				Length: nil,
				JSON: []RequireJSON{
					{
						JSONPath: "$",
						RequireXPath: RequireXPath{
							XPath: "",
							RequireMatch: RequireMatch{
								Regexp: RequireFieldMatch{
									`.`,
								},
								NotRegexp: nil,
							},
						},
					},
				},
				XML:  nil,
				HTML: nil,
			},
			wantErr: false,
		},
		{
			name: "invalid_JSON",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type:   "",
				Length: nil,
				JSON: []RequireJSON{
					{
						JSONPath: "$",
						RequireXPath: RequireXPath{
							XPath: "",
							RequireMatch: RequireMatch{
								Regexp: RequireFieldMatch{
									`(`,
								},
								NotRegexp: nil,
							},
						},
					},
				},
				XML:  nil,
				HTML: nil,
			},
			wantErr: true,
		},
		{
			name: "valid_XML",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML: []RequireXPath{
					{
						XPath: "//key",
						RequireMatch: RequireMatch{
							Regexp: RequireFieldMatch{
								`.`,
							},
							NotRegexp: nil,
						},
					},
				},
				HTML: nil,
			},
			wantErr: false,
		},
		{
			name: "invalid_XML",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML: []RequireXPath{
					{
						XPath: "//key",
						RequireMatch: RequireMatch{
							Regexp: RequireFieldMatch{
								`(`,
							},
							NotRegexp: nil,
						},
					},
				},
				HTML: nil,
			},
			wantErr: true,
		},
		{
			name: "valid_HTML",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML: []RequireXPath{
					{
						XPath: "//key",
						RequireMatch: RequireMatch{
							Regexp: RequireFieldMatch{
								`.`,
							},
							NotRegexp: nil,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid_HTML",
			fields: fields{
				RequireMatch: RequireMatch{
					Regexp:    nil,
					NotRegexp: nil,
				},
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML: []RequireXPath{
					{
						XPath: "//key",
						RequireMatch: RequireMatch{
							Regexp: RequireFieldMatch{
								`(`,
							},
							NotRegexp: nil,
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &HttpArgsRequireContent{
				RequireMatch: tt.fields.RequireMatch,
				Type:         tt.fields.Type,
				Length:       tt.fields.Length,
				JSON:         tt.fields.JSON,
				XML:          tt.fields.XML,
				HTML:         tt.fields.HTML,
			}
			if err := a.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
