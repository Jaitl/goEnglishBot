package command

import (
	"reflect"
	"testing"
)

func Test_parseTextCommand(t *testing.T) {
	type args struct {
		chatId int
		cmd    string
	}
	tests := []struct {
		name string
		args args
		want Command
	}{
		{"add", args{1, "/add слово"}, &AddCommand{1, "слово"}},
		{"text", args{1, "text"}, &TextCommand{1, "text"}},
		{"list", args{1, "/list"}, &ListCommand{1}},
		{"audio", args{1, "/audio 10"}, &AudioCommand{1, 10}},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := parseTextCommand(tt.args.chatId, tt.args.cmd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTextCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
