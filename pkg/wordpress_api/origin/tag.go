package origin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ray31245/seo_cluster/pkg/wordpress_api/model"
)

// ListTag is a function to list tag
func ListTag(ctx context.Context, baseURL string, basicAuth model.BasicAuthentication, args model.ListTagArgs) (model.ListTagResponse, model.PageSchema, error) {
	param, err := json.Marshal(args)
	if err != nil {
		return model.ListTagResponse{}, model.PageSchema{}, fmt.Errorf("marshal error: %w", err)
	}

	paramsMap := map[string]interface{}{}

	err = json.Unmarshal(param, &paramsMap)
	if err != nil {
		return model.ListTagResponse{}, model.PageSchema{}, fmt.Errorf("unmarshal error: %w", err)
	}

	// paramsMap["_fields"] = "id"

	route := "tags"

	resBody, header, err := doRequest(ctx, baseURL, http.MethodGet, route, basicAuth, paramsMap, nil)
	if err != nil {
		return model.ListTagResponse{}, model.PageSchema{}, fmt.Errorf("list tag error: %w", err)
	}

	// log.Println(string(resBody))

	resData := model.ListTagResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.ListTagResponse{}, model.PageSchema{}, fmt.Errorf("unmarshal error: %w", err)
	}

	total := 0
	if header.Header.Get("X-WP-Total") != "" {
		total, err = strconv.Atoi(header.Header.Get("X-WP-Total"))
		if err != nil {
			return model.ListTagResponse{}, model.PageSchema{}, fmt.Errorf("strconv error: %w", err)
		}
	}

	totalPages := 0
	if header.Header.Get("X-WP-TotalPages") != "" {
		totalPages, err = strconv.Atoi(header.Header.Get("X-WP-TotalPages"))
		if err != nil {
			return model.ListTagResponse{}, model.PageSchema{}, fmt.Errorf("strconv error: %w", err)
		}
	}

	page := model.PageSchema{
		Total:      total,
		TotalPages: totalPages,
	}

	return resData, page, nil
}

// CreateTag is a function to create tag
func CreateTag(ctx context.Context, BaseUrl string, basicAuth model.BasicAuthentication, args model.CreateTagArgs) (model.CreateTagResponse, error) {
	bytesData, err := json.Marshal(args)
	if err != nil {
		return model.CreateTagResponse{}, fmt.Errorf("marshal error: %w", err)
	}

	route := "tags"

	resBody, _, err := doRequest(ctx, BaseUrl, http.MethodPost, route, basicAuth, nil, bytesData)
	if err != nil {
		return model.CreateTagResponse{}, fmt.Errorf("create tag error: %w", err)
	}

	resData := model.CreateTagResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.CreateTagResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}
