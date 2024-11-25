package origin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ray31245/seo_cluster/pkg/wordpress_api/model"
)

// ListComment is a function to list comment.
func ListComment(ctx context.Context, baseURL string, basicAuth model.BasicAuthentication, args model.ListCommentArgs) (model.ListCommentResponse, error) {
	param, err := json.Marshal(args)
	if err != nil {
		return model.ListCommentResponse{}, fmt.Errorf("marshal error: %w", err)
	}

	paramsMap := map[string]interface{}{}

	err = json.Unmarshal(param, &paramsMap)
	if err != nil {
		return model.ListCommentResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	// paramsMap["_fields"] = "id"

	route := "comments"

	resBody, err := doRequest(ctx, baseURL, http.MethodGet, route, basicAuth, paramsMap, nil)
	if err != nil {
		return model.ListCommentResponse{}, fmt.Errorf("list comment error: %w", err)
	}

	// log.Println(string(resBody))

	resData := model.ListCommentResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.ListCommentResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}

// CreateComment is a function to create comment.
func CreateComment(ctx context.Context, baseURL string, basicAuth model.BasicAuthentication, args model.CreateCommentArgs) (model.CreateCommentResponse, error) {
	bytesData, err := json.Marshal(args)
	if err != nil {
		return model.CreateCommentResponse{}, fmt.Errorf("marshal error: %w", err)
	}

	route := "comments"

	resBody, err := doRequest(ctx, baseURL, http.MethodPost, route, basicAuth, nil, bytesData)
	if err != nil {
		return model.CreateCommentResponse{}, fmt.Errorf("create comment error: %w", err)
	}

	resData := model.CreateCommentResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.CreateCommentResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}
