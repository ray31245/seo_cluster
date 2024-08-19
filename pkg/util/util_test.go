package util_test

import (
	"testing"

	"goTool/pkg/util"
)

func TestEscapeHTMLMarshal(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			got, err := util.EscapeHTMLMarshal(tt.args.art)
			if (err != nil) != tt.wantErr {
				t.Errorf("EscapeHTMLMarshal() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if string(got) != tt.want {
				t.Errorf("EscapeHTMLMarshal() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
