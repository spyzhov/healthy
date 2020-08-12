package internal

import (
	"reflect"
	"testing"
	"time"
)

func TestDuration_MarshalJSON(t *testing.T) {
	type fields struct {
		Duration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "nil",
			fields: fields{
				Duration: 0,
			},
			want:    []byte(`"0s"`),
			wantErr: false,
		},
		{
			name: "1s",
			fields: fields{
				Duration: 1 * time.Second,
			},
			want:    []byte(`"1s"`),
			wantErr: false,
		},
		{
			name: "1m",
			fields: fields{
				Duration: 1 * time.Minute,
			},
			want:    []byte(`"1m0s"`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Duration{
				Duration: tt.fields.Duration,
			}
			got, err := d.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Duration
		wantErr bool
	}{
		{
			name:    "nil",
			args:    args{b: []byte(``)},
			want:    Duration{},
			wantErr: true,
		},
		{
			name:    "valid(1s)",
			args:    args{b: []byte(`"1s"`)},
			want:    Duration{Duration: time.Second},
			wantErr: false,
		},
		{
			name:    "valid(1)",
			args:    args{b: []byte(`1`)},
			want:    Duration{Duration: time.Nanosecond},
			wantErr: false,
		},
		{
			name:    "invalid(1v)",
			args:    args{b: []byte(`"1v"`)},
			want:    Duration{},
			wantErr: true,
		},
		{
			name:    "invalid(false)",
			args:    args{b: []byte(`false`)},
			want:    Duration{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Duration{}
			if err := d.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			} else if d.Duration != tt.want.Duration {
				t.Errorf("Duration not eq %v != %v", d.Duration, tt.want.Duration)
			}
		})
	}
}

func TestDuration_Validate(t *testing.T) {
	type fields struct {
		Duration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "nil",
			fields: fields{
				Duration: 0,
			},
			wantErr: false,
		},
		{
			name: "valid",
			fields: fields{
				Duration: 1,
			},
			wantErr: false,
		},
		{
			name: "invalid",
			fields: fields{
				Duration: -1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Duration{
				Duration: tt.fields.Duration,
			}
			if err := d.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
