package util_test

import (
	"testing"

	"github.com/ray31245/seo_cluster/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListImageSrcFromHtml(t *testing.T) {
	t.Parallel()

	type args struct {
		body []byte
	}

	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				body: []byte(`<p>test</p><img src="https://picsum.photos/200/300" /><img src="https://picsum.photos/500/500" />`),
			},
			want: []string{"https://picsum.photos/200/300", "https://picsum.photos/500/500"},
		},
		{
			name: "no image",
			args: args{
				body: []byte(`<p>test</p>`),
			},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)
			require := require.New(t)

			got, err := util.ListImageSrcFromHtml(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListImageSrc() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			// for _, v := range got {
			// 	log.Println(v)
			// }
			require.Equal(len(tt.want), len(got))
			assert.ElementsMatch(got, tt.want)
		})
	}
}

func TestGenImageListEncodeDiv(t *testing.T) {
	t.Parallel()

	type args struct {
		body []byte
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				body: []byte(`<p>test</p><img src="https://picsum.photos/200/300" /><img src="https://picsum.photos/500/500" /></body></html>`),
			},
			want: string([]byte(`<div class="encodeImageList">PGltZyBzcmM9Imh0dHBzOi8vcGljc3VtLnBob3Rvcy8yMDAvMzAwIiAvPjxpbWcgc3JjPSJodHRwczovL3BpY3N1bS5waG90b3MvNTAwLzUwMCIgLz4=</div>`)),
		},
		{
			name: "no image",
			args: args{
				body: []byte(`<p>test</p>`),
			},
			want: string([]byte(`<div class="encodeImageList"></div>`)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := util.GenImageListEncodeDiv(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenImageListEncodeDiv() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if got != tt.want {
				t.Errorf("GenImageListEncodeDiv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeImageListDivFromHTMl(t *testing.T) {
	t.Parallel()

	type args struct {
		body []byte
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				body: []byte(`<p>test</p><div class="encodeImageList">PGltZyBzcmM9Imh0dHBzOi8vcGljc3VtLnBob3Rvcy8yMDAvMzAwIiAvPjxpbWcgc3JjPSJodHRwczovL3BpY3N1bS5waG90b3MvNTAwLzUwMCIgLz4=</div>`),
			},
			want: `<p>test</p><img src="https://picsum.photos/200/300"/><img src="https://picsum.photos/500/500"/>`,
		},
		{
			name: "no image",
			args: args{
				body: []byte(`<p>test</p><div class="encodeImageList"></div>`),
			},
			want: "<p>test</p>",
		},
		{
			name: "image list div not found",
			args: args{
				body: []byte(`<p>test</p>`),
			},
			want: "<p>test</p>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := util.DecodeImageListDivFromHTMl(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeImageListDivFromHTMl() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if got != tt.want {
				t.Errorf("DecodeImageListDivFromHTMl() = %v, want %v", got, tt.want)
			}
		})
	}
}
