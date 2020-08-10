package executor

import "testing"

func TestHttpArgsRequireHeader_Validate(t *testing.T) {
	type fields struct {
		Exists    []string
		NotExists []string
		Eq        map[string][]string
		Regexp    map[string][]string
	}
	tests := []struct {
		name    string
		fields  fields
		header  map[string][]string
		wantErr bool
	}{
		{
			name:    "Blank",
			fields:  fields{},
			header:  make(map[string][]string, 0),
			wantErr: false,
		},
		// fixme
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &HttpArgsRequireHeader{
				Exists:    tt.fields.Exists,
				NotExists: tt.fields.NotExists,
				Eq:        tt.fields.Eq,
				Regexp:    tt.fields.Regexp,
			}
			if err := a.Match(tt.header); (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
