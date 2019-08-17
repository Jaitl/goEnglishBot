package aws

import "testing"

func TestParseUrl(t *testing.T) {
	type args struct {
		s3Url string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{"root", args{"https://s3.eu-west-2.amazonaws.com/test-bucket/file.json"}, "test-bucket", "/file.json"},
		{"sub", args{"https://s3.eu-west-2.amazonaws.com/test-bucket/folder/file.json"}, "test-bucket", "/folder/file.json"},
		{"domainRoot", args{"https://test-bucket.s3.eu-west-2.amazonaws.com/file.mp3"}, "test-bucket", "/file.mp3"},
		{"domainSub", args{"https://test-bucket.s3.eu-west-2.amazonaws.com/folder/file.mp3"}, "test-bucket", "/folder/file.mp3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ParseUrl(tt.args.s3Url)
			if got != tt.want {
				t.Errorf("ParseUrl() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseUrl() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
