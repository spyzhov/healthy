package args

import "testing"

func TestHttpArgsRequireStatus_Match(t *testing.T) {
	type fields struct {
		RequireNumeric RequireNumeric
	}
	type args struct {
		status int
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
				RequireNumeric: RequireNumeric{
					In:    nil,
					NotIn: nil,
					Eq:    nil,
					Le:    nil,
					Leq:   nil,
					Ge:    nil,
					Geq:   nil,
					Not:   nil,
				},
			},
			args: args{
				status: 0,
			},
			wantErr: false,
		},
		{
			name: "blank",
			fields: fields{
				RequireNumeric: RequireNumeric{
					In:    make([]float64, 0),
					NotIn: make([]float64, 0),
					Eq:    nil,
					Le:    nil,
					Leq:   nil,
					Ge:    nil,
					Geq:   nil,
					Not:   nil,
				},
			},
			args: args{
				status: 0,
			},
			wantErr: false,
		},
		{
			name: "valid",
			fields: fields{
				RequireNumeric: RequireNumeric{
					In:    []float64{200},
					NotIn: make([]float64, 0),
					Eq:    nil,
					Le:    nil,
					Leq:   nil,
					Ge:    nil,
					Geq:   nil,
					Not:   nil,
				},
			},
			args: args{
				status: 200,
			},
			wantErr: false,
		},
		{
			name: "invalid",
			fields: fields{
				RequireNumeric: RequireNumeric{
					In:    []float64{200},
					NotIn: make([]float64, 0),
					Eq:    nil,
					Le:    nil,
					Leq:   nil,
					Ge:    nil,
					Geq:   nil,
					Not:   nil,
				},
			},
			args: args{
				status: 201,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &HttpArgsRequireStatus{
				RequireNumeric: tt.fields.RequireNumeric,
			}
			if err := a.Match(tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
