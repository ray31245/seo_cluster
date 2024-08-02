package util

import (
	"testing"
)

func TestEscapeHTMLMarshual(t *testing.T) {
	type args struct {
		art map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				art: map[string]interface{}{
					"Content": "content<p>content</p>",
				},
			},
			want: "{\"Content\":\"content<p>content</p>\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EscapeHTMLMarshual(tt.args.art)
			if (err != nil) != tt.wantErr {
				t.Errorf("EscapeHTMLMarshual() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("EscapeHTMLMarshual() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
