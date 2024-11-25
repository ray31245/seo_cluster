package origin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ray31245/seo_cluster/pkg/wordpress_api/model"
)

// ListCategory is a function to list category
func ListCategory(ctx context.Context, baseURL string, basicAuth model.BasicAuthentication, args model.ListCategoryArgs) (model.ListCategoryResponse, error) {
	param, err := json.Marshal(args)
	if err != nil {
		return model.ListCategoryResponse{}, fmt.Errorf("marshal error: %w", err)
	}

	paramsMap := map[string]interface{}{}

	err = json.Unmarshal(param, &paramsMap)
	if err != nil {
		return model.ListCategoryResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	// paramsMap["_fields"] = "id"

	route := "categories"

	resBody, err := doRequest(ctx, baseURL, http.MethodGet, route, basicAuth, paramsMap, nil)
	if err != nil {
		return model.ListCategoryResponse{}, fmt.Errorf("list category error: %w", err)
	}

	// log.Println(string(resBody))

	resData := model.ListCategoryResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.ListCategoryResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}
