package phrase

import "testing"

func TestClear(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"he`s", args{"he`s "}, "he's"},
		{"he‘s", args{" he‘s"}, "he's"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clear(tt.args.text); got != tt.want {
				t.Errorf("Clear() = %v, want %v", got, tt.want)
			}
		})
	}
}