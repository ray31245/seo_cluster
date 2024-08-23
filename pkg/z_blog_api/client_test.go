package zblogapi_test

import (
	"context"
	"log"
	"reflect"
	"testing"

	zAPI "github.com/ray31245/seo_cluster/pkg/z_blog_api"
	zBlogErr "github.com/ray31245/seo_cluster/pkg/z_blog_api/error"
	"github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
	"github.com/ray31245/seo_cluster/pkg/z_blog_api/origin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_ListMember(t *testing.T) {
	t.Parallel()

	srv := newFullMockServer()
	faultSrv := newMockFaultServer(origin.ModMember, origin.ActList)

	t.Cleanup(func() {
		srv.Close()
		faultSrv.Close()
	})

	type fields struct {
		baseURL  string
		userName string
		password string
	}

	type args struct{}

	tests := []struct {
		name       string
		fields     fields
		args       args
		want       []model.Member
		wantErr    bool
		expectErrs []error
	}{
		{
			name: "normal",
			fields: fields{
				baseURL:  srv.URL,
				userName: zAPI.TestUserName,
				password: zAPI.TestPassword,
			},
			want: mockMembers,
		},
		{
			name: "fault",
			fields: fields{
				baseURL:  faultSrv.URL,
				userName: zAPI.TestUserName,
				password: zAPI.TestPassword,
			},
			wantErr: true,
			expectErrs: []error{
				zBlogErr.ErrHTTPInternal,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)
			require := require.New(t)

			tr, err := zAPI.NewClient(context.Background(), tt.fields.baseURL, tt.fields.userName, tt.fields.password)
			require.NoError(err)

			got, err := tr.ListMember(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListMember() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err == nil {
				assert.Equal(tt.want, got)
			} else {
				for _, expectErr := range tt.expectErrs {
					assert.ErrorIs(err, expectErr)
				}
			}
		})
	}
}

func TestClient_PostArticle(t *testing.T) {
	t.Parallel()

	srv := newFullMockServer()
	faultSrv := newMockFaultServer(origin.ModPost, origin.ActPost)

	t.Cleanup(func() {
		srv.Close()
		faultSrv.Close()
	})

	type fields struct {
		baseURL  string
		userName string
		password string
	}

	type args struct {
		art model.PostArticleRequest
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		expectErrs []error
	}{
		{
			name: "normal",
			fields: fields{
				baseURL:  srv.URL,
				userName: zAPI.TestUserName,
				password: zAPI.TestPassword,
			},
		},
		{
			name: "fault",
			fields: fields{
				baseURL:  faultSrv.URL,
				userName: zAPI.TestUserName,
				password: zAPI.TestPassword,
			},
			wantErr: true,
			expectErrs: []error{
				zBlogErr.ErrHTTPInternal,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)
			require := require.New(t)

			tr, err := zAPI.NewClient(context.Background(), tt.fields.baseURL, tt.fields.userName, tt.fields.password)
			require.NoError(err)

			if err := tr.PostArticle(context.Background(), tt.args.art); (err != nil) != tt.wantErr {
				t.Errorf("Client.PostArticle() error = %v, wantErr %v", err, tt.wantErr)
			} else if len(tt.expectErrs) > 0 {
				for _, expectErr := range tt.expectErrs {
					log.Println(err)
					assert.ErrorIs(err, expectErr)
				}
			}
		})
	}
}

func TestClient_ListArticle(t *testing.T) {
	t.Parallel()

	srv := newFullMockServer()
	faultSrv := newMockFaultServer(origin.ModPost, origin.ActList)

	t.Cleanup(func() {
		srv.Close()
		faultSrv.Close()
	})

	type fields struct {
		baseURL  string
		userName string
		password string
	}

	type args struct {
		req model.ListArticleRequest
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		want       []model.Article
		wantErr    bool
		expectErrs []error
	}{
		{
			name: "normal",
			fields: fields{
				baseURL:  srv.URL,
				userName: zAPI.TestUserName,
				password: zAPI.TestPassword,
			},
			want: mockArticles,
		},
		{
			name: "filter",
			fields: fields{
				baseURL:  srv.URL,
				userName: zAPI.TestUserName,
				password: zAPI.TestPassword,
			},
			args: args{
				req: model.ListArticleRequest{
					CateID: 1,
				},
			},
			want: []model.Article{mockArticles[0]},
		},
		{
			name: "fault",
			fields: fields{
				baseURL:  faultSrv.URL,
				userName: zAPI.TestUserName,
				password: zAPI.TestPassword,
			},
			wantErr: true,
			expectErrs: []error{
				zBlogErr.ErrHTTPInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// assert := assert.New(t)
			require := require.New(t)
			tr, err := zAPI.NewClient(context.Background(), tt.fields.baseURL, tt.fields.userName, tt.fields.password)
			require.NoError(err)

			got, err := tr.ListArticle(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListArticle() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if len(tt.expectErrs) > 0 {
				for _, expectErr := range tt.expectErrs {
					require.ErrorIs(err, expectErr)
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListArticle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetCountOfArticle(t *testing.T) {
	t.Parallel()

	srv := newFullMockServer()
	faultSrv := newMockFaultServer(origin.ModPost, origin.ActList)

	t.Cleanup(func() {
		srv.Close()
		faultSrv.Close()
	})

	type fields struct {
		baseURL  string
		userName string
		password string
	}

	type args struct {
		req model.ListArticleRequest
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		want       int
		wantErr    bool
		exceptErrs []error
	}{
		{
			name: "normal",
			fields: fields{
				baseURL:  srv.URL,
				userName: zAPI.TestUserName,
				password: zAPI.TestPassword,
			},
			want: len(mockArticles),
		},
		{
			name: "fault",
			fields: fields{
				baseURL:  faultSrv.URL,
				userName: zAPI.TestUserName,
				password: zAPI.TestPassword,
			},
			wantErr: true,
			exceptErrs: []error{
				zBlogErr.ErrHTTPInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// assert := assert.New(t)
			require := require.New(t)
			tr, err := zAPI.NewClient(context.Background(), tt.fields.baseURL, tt.fields.userName, tt.fields.password)
			require.NoError(err)

			got, err := tr.GetCountOfArticle(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetCountOfArticle() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if len(tt.exceptErrs) > 0 {
				for _, expectErr := range tt.exceptErrs {
					require.ErrorIs(err, expectErr)
				}
			}

			if got != tt.want {
				t.Errorf("Client.GetCountOfArticle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_DeleteArticle(t *testing.T) {
	t.Parallel()

	srv := newFullMockServer()
	faultSrv := newMockFaultServer(origin.ModPost, origin.ActDelete)

	t.Cleanup(func() {
		srv.Close()
		faultSrv.Close()
	})

	type fields struct {
		baseURL  string
		userName string
		password string
	}

	type args struct {
		id string
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		expectErrs []error
	}{
		{
			name: "normal",
			fields: fields{
				baseURL:  srv.URL,
				userName: zAPI.TestUserName,
				password: zAPI.TestPassword,
			},
			args: args{
				id: "1",
			},
		},
		{
			name: "fault",
			fields: fields{
				baseURL:  faultSrv.URL,
				userName: zAPI.TestUserName,
				password: zAPI.TestPassword,
			},
			wantErr: true,
			expectErrs: []error{
				zBlogErr.ErrHTTPInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// assert := assert.New(t)
			require := require.New(t)
			tr, err := zAPI.NewClient(context.Background(), tt.fields.baseURL, tt.fields.userName, tt.fields.password)
			require.NoError(err)

			if err := tr.DeleteArticle(context.Background(), tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteArticle() error = %v, wantErr %v", err, tt.wantErr)
			} else if len(tt.expectErrs) > 0 {
				for _, expectErr := range tt.expectErrs {
					require.ErrorIs(err, expectErr)
				}
			}
		})
	}
}

func TestClient_ListCategory(t *testing.T) {
	t.Parallel()

	srv := newFullMockServer()
	faultSrv := newMockFaultServer(origin.ModCategory, origin.ActList)

	t.Cleanup(func() {
		srv.Close()
		faultSrv.Close()
	})

	type fields struct {
		baseURL  string
		userName string
		password string
	}

	type args struct{}

	tests := []struct {
		name       string
		fields     fields
		args       args
		want       []model.Category
		wantErr    bool
		expectErrs []error
	}{
		{
			name: "normal",
			fields: fields{
				baseURL:  srv.URL,
				userName: zAPI.TestUserName,
				password: zAPI.TestPassword,
			},
			want: mockCategories,
		},
		{
			name: "fault",
			fields: fields{
				baseURL:  faultSrv.URL,
				userName: zAPI.TestUserName,
				password: zAPI.TestPassword,
			},
			wantErr: true,
			expectErrs: []error{
				zBlogErr.ErrHTTPInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// assert := assert.New(t)
			require := require.New(t)
			tr, err := zAPI.NewClient(context.Background(), tt.fields.baseURL, tt.fields.userName, tt.fields.password)
			require.NoError(err)

			got, err := tr.ListCategory(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListCategory() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if len(tt.expectErrs) > 0 {
				for _, expectErr := range tt.expectErrs {
					require.ErrorIs(err, expectErr)
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}
