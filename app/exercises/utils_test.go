package exercises

import "testing"

func TestClearText(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"dot1", args{"test. This text"}, "test this text"},
		{"dot2", args{"test, This text"}, "test this text"},
		{"dot3", args{"test! This text"}, "test this text"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClearText(tt.args.text); got != tt.want {
				t.Errorf("ClearText() = %v, want %v", got, tt.want)
			}
		})
	}
}