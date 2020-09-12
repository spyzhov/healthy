package internal

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestReaderCloserCallback_Read(t *testing.T) {
	tests := []struct {
		name    string
		reader  *ReaderCloserCallback
		wantBuf string
		wantN   int
		wantErr bool
	}{
		{
			name:    "nil",
			reader:  new(ReaderCloserCallback),
			wantBuf: "",
			wantN:   0,
			wantErr: false,
		},
		{
			name:    "callback: nil",
			reader:  NewReaderCloserCallback(ioutil.NopCloser(strings.NewReader("simple")), nil),
			wantBuf: "simple",
			wantN:   6,
			wantErr: false,
		},
		{
			name: "simple",
			reader: NewReaderCloserCallback(
				ioutil.NopCloser(strings.NewReader("simple example")),
				func(p []byte, i int) ([]byte, error) {
					result := make([]byte, i)
					copy(result, p)
					return result, nil
				},
			),
			wantBuf: "simple example",
			wantN:   14,
			wantErr: false,
		},
		{
			name: "replace",
			reader: NewReaderCloserCallback(
				ioutil.NopCloser(strings.NewReader("1X22X333X")),
				func(p []byte, n int) ([]byte, error) {
					result := make([]byte, 0, len(p))
					for i := 0; i < n; i++ {
						if p[i] == 'X' {
							result = append(result, 'X')
						}
						result = append(result, p[i])
					}
					return result, nil
				},
			),
			wantBuf: "1XX22XX333XX",
			wantN:   12,
			wantErr: false,
		},
		{
			name: "replace: long",
			reader: NewReaderCloserCallback(
				ioutil.NopCloser(strings.NewReader(strings.Repeat("X", 1024))),
				func(p []byte, n int) ([]byte, error) {
					result := make([]byte, 0, len(p))
					for i := 0; i < n; i++ {
						if p[i] == 'X' {
							result = append(result, 'X')
						}
						result = append(result, p[i])
					}
					return result, nil
				},
			),
			wantBuf: strings.Repeat("X", 2048),
			wantN:   2048,
			wantErr: false,
		},
		{
			name: "replace: long XY",
			reader: NewReaderCloserCallback(
				ioutil.NopCloser(strings.NewReader(strings.Repeat("Y", 1024))),
				func(p []byte, n int) ([]byte, error) {
					result := make([]byte, 0, len(p))
					for i := 0; i < n; i++ {
						if p[i] == 'Y' {
							result = append(result, 'X')
						}
						result = append(result, p[i])
					}
					return result, nil
				},
			),
			wantBuf: strings.Repeat("XY", 1024),
			wantN:   2048,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBuf, err := ioutil.ReadAll(tt.reader)
			gotN := len(gotBuf)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Read() gotN = %v, want %v", gotN, tt.wantN)
			}
			if string(gotBuf) != tt.wantBuf {
				t.Errorf("Read() gotBuf = `%v`, want `%v`", string(gotBuf), tt.wantBuf)
			}
		})
	}
}
