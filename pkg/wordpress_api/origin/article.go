package origin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ray31245/seo_cluster/pkg/util"
	"github.com/ray31245/seo_cluster/pkg/wordpress_api/model"
)

func ListArticle(ctx context.Context, baseURL string, basicAuth model.BasicAuthentication, args model.ListArticleArgs) (model.ListArticleResponse, error) {
	param, err := json.Marshal(args)
	if err != nil {
		return model.ListArticleResponse{}, fmt.Errorf("marshal error: %w", err)
	}

	paramsMap := map[string]interface{}{}

	err = json.Unmarshal(param, &paramsMap)
	if err != nil {
		return model.ListArticleResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	// paramsMap["_fields"] = "id"

	route := "posts"

	resBody, _, err := doRequest(ctx, baseURL, http.MethodGet, route, basicAuth, paramsMap, nil)
	if err != nil {
		return model.ListArticleResponse{}, fmt.Errorf("list article error: %w", err)
	}

	// log.Println(string(resBody))

	resData := model.ListArticleResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.ListArticleResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}

func CreateArticle(ctx context.Context, baseURL string, basicAuth model.BasicAuthentication, args model.CreateArticleArgs) (model.CreateArticleResponse, error) {
	bytesData, err := util.EscapeHTMLMarshal(args)
	if err != nil {
		return model.CreateArticleResponse{}, fmt.Errorf("marshal error: %w", err)
	}

	// log.Println(string(bytesData))

	route := "posts"

	resBody, _, err := doRequest(ctx, baseURL, http.MethodPost, route, basicAuth, nil, bytesData)
	if err != nil {
		return model.CreateArticleResponse{}, fmt.Errorf("create article error: %w", err)
	}

	resData := model.CreateArticleResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.CreateArticleResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}

func UpdateArticle(ctx context.Context, baseURL string, basicAuth model.BasicAuthentication, args model.UpdateArticleArgs) (model.UpdateArticleResponse, error) {
	bytesData, err := util.EscapeHTMLMarshal(args)
	if err != nil {
		return model.UpdateArticleResponse{}, fmt.Errorf("marshal error: %w", err)
	}

	route := fmt.Sprintf("posts/%d", args.ID)

	resBody, _, err := doRequest(ctx, baseURL, http.MethodPost, route, basicAuth, nil, bytesData)
	if err != nil {
		return model.UpdateArticleResponse{}, fmt.Errorf("update article error: %w", err)
	}

	resData := model.UpdateArticleResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.UpdateArticleResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}

func RetrieveArticle(ctx context.Context, baseURL string, basicAuth model.BasicAuthentication, args model.RetrieveArticleArgs) (model.RetrieveArticleResponse, error) {
	param, err := json.Marshal(args)
	if err != nil {
		return model.RetrieveArticleResponse{}, fmt.Errorf("marshal error: %w", err)
	}

	paramsMap := map[string]interface{}{}

	err = json.Unmarshal(param, &paramsMap)
	if err != nil {
		return model.RetrieveArticleResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	route := fmt.Sprintf("posts/%d", args.ID)

	// paramsMap["_fields"] = "id,title,content"

	resBody, _, err := doRequest(ctx, baseURL, http.MethodGet, route, basicAuth, paramsMap, nil)
	if err != nil {
		return model.RetrieveArticleResponse{}, fmt.Errorf("retrieve article error: %w", err)
	}

	resData := model.RetrieveArticleResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.RetrieveArticleResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}
