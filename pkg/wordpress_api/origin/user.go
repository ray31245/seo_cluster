package origin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ray31245/seo_cluster/pkg/wordpress_api/model"
)

func RetrieveUserMe(ctx context.Context, baseURL string, basicAuth model.BasicAuthentication, apiContext model.ApiContext) (model.RetrieveUserMeResponse, error) {
	paramsMap := map[string]interface{}{}
	if apiContext != "" {
		paramsMap["context"] = apiContext
	}

	route := "users/me"

	// paramsMap["_fields"] = "id,name,slug,email,roles,avatar_urls"

	resBody, err := doRequest(ctx, baseURL, http.MethodGet, route, basicAuth, paramsMap, nil)
	if err != nil {
		return model.RetrieveUserMeResponse{}, fmt.Errorf("retrieve user me error: %w", err)
	}

	resData := model.RetrieveUserMeResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.RetrieveUserMeResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}
