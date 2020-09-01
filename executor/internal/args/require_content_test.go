package args

import "testing"

func TestHttpArgsRequireContent_Match(t *testing.T) {
	type fields struct {
		Type   HttpArgsRequireContentType
		Length *RequireNumeric
		JSON   []RequireJSON
		XML    []RequireXPath
		HTML   []RequireXPath
		Text   []RequireMatch
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
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
				Text:   nil,
			},
			args: args{
				content: nil,
			},
			wantErr: false,
		},
		{
			name: "blank",
			fields: fields{
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
				Text: make([]RequireMatch, 0),
			},
			args: args{
				content: make([]byte, 0),
			},
			wantErr: false,
		},
		{
			name: "valid_match",
			fields: fields{
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
				Text: []RequireMatch{
					{
						Regexp: RequireFieldMatch{
							`content`,
						},
						NotRegexp: nil,
					},
				},
			},
			args: args{
				content: []byte(`content`),
			},
			wantErr: false,
		},
		{
			name: "invalid_match",
			fields: fields{
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
				Text: []RequireMatch{
					{
						Regexp: RequireFieldMatch{
							`wrong`,
						},
						NotRegexp: nil,
					},
				},
			},
			args: args{
				content: []byte(`content`),
			},
			wantErr: true,
		},
		{
			name: "valid_not_match",
			fields: fields{
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
				Text: []RequireMatch{
					{
						Regexp: nil,
						NotRegexp: RequireFieldMatchNot{
							`wrong`,
						},
					},
				},
			},
			args: args{
				content: []byte(`content`),
			},
			wantErr: false,
		},
		{
			name: "invalid_not_match",
			fields: fields{
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
				Text: []RequireMatch{
					{
						Regexp: nil,
						NotRegexp: RequireFieldMatchNot{
							`content`,
						},
					},
				},
			},
			args: args{
				content: []byte(`content`),
			},
			wantErr: true,
		},
		{
			name: "valid_type",
			fields: fields{
				Type:   "JSON",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
				Text:   nil,
			},
			args: args{
				content: []byte(`{"foo": true}`),
			},
			wantErr: false,
		},
		{
			name: "invalid_type",
			fields: fields{
				Type:   "JSON",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
				Text:   nil,
			},
			args: args{
				content: []byte(`{"foo"`),
			},
			wantErr: true,
		},
		{
			name: "valid_length",
			fields: fields{
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
				Text: nil,
			},
			args: args{
				content: []byte(`content`),
			},
			wantErr: false,
		},
		{
			name: "invalid_length",
			fields: fields{
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
				Text: nil,
			},
			args: args{
				content: []byte(`content`),
			},
			wantErr: true,
		},
		{
			name: "valid_json",
			fields: fields{
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
				Text: nil,
			},
			args: args{
				content: []byte(`{}`),
			},
			wantErr: false,
		},
		{
			name: "invalid_json",
			fields: fields{
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
				Text: nil,
			},
			args: args{
				content: []byte(`{}`),
			},
			wantErr: true,
		},
		{
			name: "valid_xml",
			fields: fields{
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
				Text: nil,
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
				Text: nil,
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
				Text: nil,
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
				Text: nil,
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
			a := &RequireContent{
				Type:   tt.fields.Type,
				Length: tt.fields.Length,
				JSON:   tt.fields.JSON,
				XML:    tt.fields.XML,
				HTML:   tt.fields.HTML,
				Text:   tt.fields.Text,
			}
			if err := a.Match("content", tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttpArgsRequireContent_Validate(t *testing.T) {
	type fields struct {
		Type   HttpArgsRequireContentType
		Length *RequireNumeric
		JSON   []RequireJSON
		XML    []RequireXPath
		HTML   []RequireXPath
		Text   []RequireMatch
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "nil",
			fields: fields{
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
				Text:   nil,
			},
			wantErr: false,
		},
		{
			name: "blank",
			fields: fields{
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
				Text: make([]RequireMatch, 0),
			},
			wantErr: false,
		},
		{
			name: "valid_RequireMatch_Regexp",
			fields: fields{
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
				Text: []RequireMatch{
					{
						Regexp: RequireFieldMatch{
							`.`,
						},
						NotRegexp: nil,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid_RequireMatch_Regexp",
			fields: fields{
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
				Text: []RequireMatch{
					{
						Regexp: RequireFieldMatch{
							`(`,
						},
						NotRegexp: nil,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid_RequireMatch_NotRegexp",
			fields: fields{
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
				Text: []RequireMatch{
					{
						Regexp: nil,
						NotRegexp: RequireFieldMatchNot{
							`.`,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid_RequireMatch_NotRegexp",
			fields: fields{
				Type:   "",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
				Text: []RequireMatch{
					{
						Regexp: nil,
						NotRegexp: RequireFieldMatchNot{
							`(`,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid_Type",
			fields: fields{
				Type:   "JSON",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
				Text:   nil,
			},
			wantErr: false,
		},
		{
			name: "invalid_Type",
			fields: fields{
				Type:   "WRONG",
				Length: nil,
				JSON:   nil,
				XML:    nil,
				HTML:   nil,
				Text:   nil,
			},
			wantErr: true,
		},
		{
			name: "valid_JSON",
			fields: fields{
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
				Text: nil,
			},
			wantErr: false,
		},
		{
			name: "invalid_JSON",
			fields: fields{
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
				Text: nil,
			},
			wantErr: true,
		},
		{
			name: "valid_XML",
			fields: fields{
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
				Text: nil,
			},
			wantErr: false,
		},
		{
			name: "invalid_XML",
			fields: fields{
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
				Text: nil,
			},
			wantErr: true,
		},
		{
			name: "valid_HTML",
			fields: fields{
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
				Text: nil,
			},
			wantErr: false,
		},
		{
			name: "invalid_HTML",
			fields: fields{
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
				Text: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &RequireContent{
				Type:   tt.fields.Type,
				Length: tt.fields.Length,
				JSON:   tt.fields.JSON,
				XML:    tt.fields.XML,
				HTML:   tt.fields.HTML,
				Text:   tt.fields.Text,
			}
			if err := a.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
