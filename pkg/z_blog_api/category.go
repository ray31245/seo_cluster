package zblogapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"goTool/pkg/z_blog_api/model"
)

func (t *Client) listCategory(ctx context.Context) (model.ListCategoryResponse, error) {
	resBody, err := t.requestWithBlock(ctx, http.MethodGet, map[string]interface{}{"mod": "category", "act": "list"}, nil)
	if err != nil {
		return model.ListCategoryResponse{}, fmt.Errorf("list category error: %w", err)
	}

	var resData model.ListCategoryResponse
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.ListCategoryResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}
