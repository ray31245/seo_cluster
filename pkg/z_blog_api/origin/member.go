package origin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
)

func Login(ctx context.Context, baseURL string, token string, userName string, password string) (model.LoginResponse, error) {
	data := map[string]string{}
	data["username"] = userName
	data["password"] = password

	bytesData, err := json.Marshal(data)
	if err != nil {
		return model.LoginResponse{}, fmt.Errorf("marshal error: %w", err)
	}

	resBody, err := doRequest(ctx, baseURL, http.MethodPost, token, map[string]interface{}{"mod": "member", "act": "login"}, bytesData)
	if err != nil {
		return model.LoginResponse{}, fmt.Errorf("login error: %w", err)
	}

	var resData model.LoginResponse
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.LoginResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}

func ListMember(ctx context.Context, baseURL string, token string) (model.ListMemberResponse, error) {
	resBody, err := doRequest(ctx, baseURL, http.MethodGet, token, map[string]interface{}{"mod": "member", "act": "list"}, nil)
	if err != nil {
		return model.ListMemberResponse{}, fmt.Errorf("list member error: %w", err)
	}

	var resData model.ListMemberResponse
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.ListMemberResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}
