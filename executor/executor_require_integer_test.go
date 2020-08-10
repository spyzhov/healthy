package executor

import "testing"

func TestRequireInteger_Validate(t *testing.T) {
	type fields struct {
		In    []int
		NotIn []int
		Eq    *int
		Le    *int
		Leq   *int
		Ge    *int
		Geq   *int
		Not   *int
	}
	type args struct {
		name    string
		integer int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Blank",
			fields: fields{},
			args: args{
				name:    "",
				integer: 0,
			},
			wantErr: false,
		},
		// fixme
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &RequireInteger{
				In:    tt.fields.In,
				NotIn: tt.fields.NotIn,
				Eq:    tt.fields.Eq,
				Le:    tt.fields.Le,
				Leq:   tt.fields.Leq,
				Ge:    tt.fields.Ge,
				Geq:   tt.fields.Geq,
				Not:   tt.fields.Not,
			}
			if err := a.Match(tt.args.name, tt.args.integer); (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
