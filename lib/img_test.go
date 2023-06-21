package lib

import (
	"reflect"
	"testing"
)

func TestGetTerminal(t *testing.T) {
	width, height := getTerminalSize()
	if reflect.TypeOf(width).String() != "int" || reflect.TypeOf(height).String() != "int" {
		t.Errorf("getTermial() is not returning int")
	}
}

func TestGetImageSize(t *testing.T) {
	type args struct {
		path string
	}
	type want struct {
		width  int
		height int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "default image path",
			args: args{"../img/BolinopsisMikado.jpg"},
			want: want{width: 950, height: 600},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotw, goth := getImageSize(tt.args.path)
			if gotw != tt.want.width || goth != tt.want.height {
				t.Errorf("wrong image size")
			}
		})
	}
}
