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
		{"me", args{1, "/me"}, &MeCommand{1}},
		{"remove", args{1, "/remove 10"}, &RemoveCommand{1, 10}},
		{"voice", args{1, "/voice 10"}, &VoiceCommand{1, 10}},
		{"puzzleAudio", args{1, "/puzzleAudio 10"}, &PuzzleAudioCommand{1, 10}},
		{"puzzleTrans", args{1, "/puzzleTrans 10"}, &PuzzleTransCommand{1, 10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := parseTextCommand(tt.args.chatId, tt.args.cmd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTextCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseNumberCommand(t *testing.T) {
	type args struct {
		userId    int
		incNumber int
	}
	tests := []struct {
		name string
		args args
		want Command
	}{
		{"card", args{1, 10}, &NumberCommand{1, 10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := parseNumberCommand(tt.args.userId, tt.args.incNumber)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseNumberCommand() got = %v, want %v", got, tt.want)
			}
		})
	}
}
