package util_test

import (
	"reflect"
	"testing"

	"github.com/ray31245/seo_cluster/pkg/util"
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

func TestMdToHTML(t *testing.T) {
	t.Parallel()

	type args struct {
		md []byte
	}

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "test1",
			args: args{
				md: []byte("# title"),
			},
			want: []byte("<h1 id=\"title\">title</h1>\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := util.MdToHTML(tt.args.md); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MdToHTML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTMLToMd(t *testing.T) {
	type args struct {
		htmlCode string
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
				htmlCode: "<h1 id=\"title\">title</h1>\n",
			},
			want: "# title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := util.HTMLToMd(tt.args.htmlCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("HTMLToMd() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if got != tt.want {
				t.Errorf("HTMLToMd() = %v, want %v", got, tt.want)
			}
		})
	}
}
