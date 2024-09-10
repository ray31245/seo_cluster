package zblogapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	zBlogErr "github.com/ray31245/seo_cluster/pkg/z_blog_api/error"
	"github.com/ray31245/seo_cluster/pkg/z_blog_api/origin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// TestUserName is the test username
	TestUserName = "admin"
	// TestPassword is the test password
	TestPassword = "admin"
	// TestToken is the test token
	TestToken = "this_is_a_test_token"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// assert route path is correct
	if r.URL.Path != "/"+origin.APIPath {
		http.Error(w, `{"code":404,"message":"not found"}`, http.StatusNotFound)

		return
	}

	// assert method is correct
	if r.Method != http.MethodPost {
		http.Error(w, `{"code":405,"message":"method not allowed"}`, http.StatusMethodNotAllowed)

		return
	}

	// assert parameter is correct
	if r.URL.Query().Get("act") != "login" {
		http.Error(w, `{"code":419,"message":"illegal access"}`, zBlogErr.StatusIllegalAccess)

		return
	}

	if r.URL.Query().Get("mod") != "member" {
		http.Error(w, `{"code":419,"message":"illegal access"}`, zBlogErr.StatusIllegalAccess)

		return
	}

	// assert login info is correct
	var reqBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, `{"code":400,"message":"bad request"}`, http.StatusBadRequest)

		return
	} else if reqBody["username"] != TestUserName || reqBody["password"] != TestPassword {
		http.Error(w, `{"code":401,"message":"unauthorized"}`, http.StatusUnauthorized)
	}

	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(fmt.Sprintf(`{"code":200,"data":{"Token":"%s"}}`, TestToken)))
	if err != nil {
		http.Error(w, `{"code":500,"message":"internal server error"}`, http.StatusInternalServerError)
	}
}

func TestClient_Login(t *testing.T) {
	t.Parallel()

	svr := httptest.NewServer(http.HandlerFunc(LoginHandler))

	t.Cleanup(func() { svr.Close() })

	type fields struct {
		baseURL  string
		token    string
		userName string
		password string
	}

	type args struct{}

	tests := []struct {
		name        string
		fields      fields
		args        args
		expectToken string
		wantErr     bool
		expectErrs  []error
	}{
		{
			name: "TestClient_Login",
			fields: fields{
				baseURL:  svr.URL,
				userName: TestUserName,
				password: TestPassword,
			},
			args:        args{},
			expectToken: TestToken,
		},
		{
			name: "TestClient_Login_Error_Not_Found",
			fields: fields{
				baseURL:  svr.URL + "/not_found",
				userName: TestUserName,
				password: TestPassword,
			},
			args:    args{},
			wantErr: true,
			expectErrs: []error{
				zBlogErr.ErrHTTPNotFound,
			},
		},
		{
			name: "TestClient_Login_Error_Incorrect_Login_Info",
			fields: fields{
				baseURL:  svr.URL,
				userName: "no_user",
				password: "no_password",
			},
			args:    args{},
			wantErr: true,
			expectErrs: []error{
				zBlogErr.ErrHTTPUnauthorized,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)
			require := require.New(t)
			tr := &Client{
				baseURL:  tt.fields.baseURL,
				token:    tt.fields.token,
				userName: tt.fields.userName,
				password: tt.fields.password,
				lock:     &sync.Mutex{},
			}

			err := tr.Login(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Login() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if len(tt.expectErrs) > 0 {
				for _, expectErr := range tt.expectErrs {
					require.ErrorIs(err, expectErr)
				}
			} else {
				assert.Equal(tt.expectToken, tr.token)
			}
		})
	}
}
