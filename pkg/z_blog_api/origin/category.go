package origin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
)

func ListCategory(ctx context.Context, baseURL string, token string) (model.ListCategoryResponse, error) {
	resBody, err := doRequest(ctx, baseURL, http.MethodGet, token, map[string]interface{}{ParamMod: ModCategory, ParamAct: ActList}, nil)
	if err != nil {
		return model.ListCategoryResponse{}, fmt.Errorf("list category error: %w", err)
	}

	var resData model.ListCategoryResponse
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.ListCategoryResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}
