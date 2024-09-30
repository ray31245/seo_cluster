package commentbot

import (
	"log"
	"testing"
	"time"

	"github.com/ray31245/seo_cluster/pkg/util"
	zModel "github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
)

func Test_computeGap(t *testing.T) {
	type args struct {
		article zModel.Article
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test1",
			args: args{
				article: zModel.Article{
					PostTime: util.UnixTime{Time: time.Now().Add(-time.Hour * 30)},
					CommNums: util.StringNumber(9),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := computeGap(tt.args.article); got != tt.want {
				log.Printf("computeGap() = %v, want %v", got, tt.want)
			}
		})
	}
}
