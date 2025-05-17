package r4

import (
	"reflect"
	"testing"
)

func Test_parseEntry(t *testing.T) {
	type args struct {
		l []byte
	}
	tests := []struct {
		name        string
		args        args
		wantStation []byte
		wantTemp    float64
	}{
		{
			name: "",
			args: args{
				l: []byte("Da Lat;19.2"),
			},
			wantStation: []byte("Da Lat"),
			wantTemp:    19.2,
		},
		{
			name: "",
			args: args{
				l: []byte("Saint Petersburg;-5.5"),
			},
			wantStation: []byte("Saint Petersburg"),
			wantTemp:    -5.5,
		},
		{
			name: "",
			args: args{
				l: []byte("Kuwait City;31.6"),
			},
			wantStation: []byte("Kuwait City"),
			wantTemp:    31.6,
		},
		{
			name: "",
			args: args{
				l: []byte("Dikson;-13.9"),
			},
			wantStation: []byte("Dikson"),
			wantTemp:    -13.9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStation, gotTemp := parseEntry(tt.args.l)
			if !reflect.DeepEqual(gotStation, tt.wantStation) {
				t.Errorf("parseEntry() gotStation = %v, want %v", gotStation, tt.wantStation)
			}
			if gotTemp != tt.wantTemp {
				t.Errorf("parseEntry() gotTemp = %v, want %v", gotTemp, tt.wantTemp)
			}
		})
	}
}

func Test(t *testing.T) {
	f := func(s []byte) {
		s[0] = byte('a')
	}
	s := []byte("bcd")
	println(string(s))
	f(s)
	println(string(s))
}
