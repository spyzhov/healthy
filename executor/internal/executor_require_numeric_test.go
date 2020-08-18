package internal

import "testing"

func TestRequireNumeric_Match(t *testing.T) {
	ref := func(i float64) *float64 {
		return &i
	}
	type fields struct {
		In    []float64
		NotIn []float64
		Eq    *float64
		Le    *float64
		Leq   *float64
		Ge    *float64
		Geq   *float64
		Not   *float64
	}
	type args struct {
		name    string
		numeric float64
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
				In:    nil,
				NotIn: nil,
				Eq:    nil,
				Le:    nil,
				Leq:   nil,
				Ge:    nil,
				Geq:   nil,
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: false,
		},
		{
			name: "in_array(true)",
			fields: fields{
				In:    []float64{0, 1, 2},
				NotIn: nil,
				Eq:    nil,
				Le:    nil,
				Leq:   nil,
				Ge:    nil,
				Geq:   nil,
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: false,
		},
		{
			name: "in_array(false)",
			fields: fields{
				In:    []float64{3, 1, 2},
				NotIn: nil,
				Eq:    nil,
				Le:    nil,
				Leq:   nil,
				Ge:    nil,
				Geq:   nil,
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: true,
		},
		{
			name: "not_in_array(true)",
			fields: fields{
				In:    nil,
				NotIn: []float64{3, 1, 2},
				Eq:    nil,
				Le:    nil,
				Leq:   nil,
				Ge:    nil,
				Geq:   nil,
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: false,
		},
		{
			name: "not_in_array(false)",
			fields: fields{
				In:    nil,
				NotIn: []float64{0, 1, 2},
				Eq:    nil,
				Le:    nil,
				Leq:   nil,
				Ge:    nil,
				Geq:   nil,
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: true,
		},
		{
			name: "eq(true)",
			fields: fields{
				In:    nil,
				NotIn: nil,
				Eq:    ref(0),
				Le:    nil,
				Leq:   nil,
				Ge:    nil,
				Geq:   nil,
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: false,
		},
		{
			name: "eq(false)",
			fields: fields{
				In:    nil,
				NotIn: nil,
				Eq:    ref(1),
				Le:    nil,
				Leq:   nil,
				Ge:    nil,
				Geq:   nil,
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: true,
		},
		{
			name: "le(true)",
			fields: fields{
				In:    nil,
				NotIn: nil,
				Eq:    nil,
				Le:    ref(1),
				Leq:   nil,
				Ge:    nil,
				Geq:   nil,
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: false,
		},
		{
			name: "le(false)",
			fields: fields{
				In:    nil,
				NotIn: nil,
				Eq:    nil,
				Le:    ref(0),
				Leq:   nil,
				Ge:    nil,
				Geq:   nil,
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: true,
		},
		{
			name: "leq(true)",
			fields: fields{
				In:    nil,
				NotIn: nil,
				Eq:    nil,
				Le:    nil,
				Leq:   ref(1),
				Ge:    nil,
				Geq:   nil,
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: false,
		},
		{
			name: "leq(true)",
			fields: fields{
				In:    nil,
				NotIn: nil,
				Eq:    nil,
				Le:    nil,
				Leq:   ref(0),
				Ge:    nil,
				Geq:   nil,
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: false,
		},
		{
			name: "leq(false)",
			fields: fields{
				In:    nil,
				NotIn: nil,
				Eq:    nil,
				Le:    nil,
				Leq:   ref(-1),
				Ge:    nil,
				Geq:   nil,
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: true,
		},
		{
			name: "ge(true)",
			fields: fields{
				In:    nil,
				NotIn: nil,
				Eq:    nil,
				Le:    nil,
				Leq:   nil,
				Ge:    ref(-1),
				Geq:   nil,
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: false,
		},
		{
			name: "ge(false)",
			fields: fields{
				In:    nil,
				NotIn: nil,
				Eq:    nil,
				Le:    nil,
				Leq:   nil,
				Ge:    ref(0),
				Geq:   nil,
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: true,
		},
		{
			name: "geq(true)",
			fields: fields{
				In:    nil,
				NotIn: nil,
				Eq:    nil,
				Le:    nil,
				Leq:   nil,
				Ge:    nil,
				Geq:   ref(-1),
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: false,
		},
		{
			name: "geq(true)",
			fields: fields{
				In:    nil,
				NotIn: nil,
				Eq:    nil,
				Le:    nil,
				Leq:   nil,
				Ge:    nil,
				Geq:   ref(0),
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: false,
		},
		{
			name: "geq(false)",
			fields: fields{
				In:    nil,
				NotIn: nil,
				Eq:    nil,
				Le:    nil,
				Leq:   nil,
				Ge:    nil,
				Geq:   ref(1),
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: true,
		},
		{
			name: "not(true)",
			fields: fields{
				In:    nil,
				NotIn: nil,
				Eq:    nil,
				Le:    nil,
				Leq:   nil,
				Ge:    nil,
				Geq:   nil,
				Not:   ref(1),
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: false,
		},
		{
			name: "not(false)",
			fields: fields{
				In:    nil,
				NotIn: nil,
				Eq:    nil,
				Le:    nil,
				Leq:   nil,
				Ge:    nil,
				Geq:   nil,
				Not:   ref(0),
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: true,
		},
		{
			name: "ge_le(true)",
			fields: fields{
				In:    nil,
				NotIn: nil,
				Eq:    nil,
				Le:    ref(1),
				Leq:   nil,
				Ge:    ref(-1),
				Geq:   nil,
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: false,
		},
		{
			name: "ge_le(false)",
			fields: fields{
				In:    nil,
				NotIn: nil,
				Eq:    nil,
				Le:    ref(-1),
				Leq:   nil,
				Ge:    ref(-3),
				Geq:   nil,
				Not:   nil,
			},
			args: args{
				name:    "",
				numeric: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &RequireNumeric{
				In:    tt.fields.In,
				NotIn: tt.fields.NotIn,
				Eq:    tt.fields.Eq,
				Le:    tt.fields.Le,
				Leq:   tt.fields.Leq,
				Ge:    tt.fields.Ge,
				Geq:   tt.fields.Geq,
				Not:   tt.fields.Not,
			}
			if err := a.Match(tt.args.name, tt.args.numeric); (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
