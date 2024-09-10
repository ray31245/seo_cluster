package origin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ray31245/seo_cluster/pkg/util"
	"github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
)

func PostArticle(ctx context.Context, baseURL string, token string, art model.PostArticleRequest) error {
	bytesData, err := util.EscapeHTMLMarshal(art)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	resBody, err := doRequest(ctx, baseURL, http.MethodPost, token, map[string]interface{}{ParamMod: ModPost, ParamAct: ActPost}, bytesData)
	if err != nil {
		return fmt.Errorf("post article error: %w", err)
	}

	resData := model.PostArticleResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	return nil
}

func GetArticle(ctx context.Context, baseURL string, token string, id string) (model.GetArticleResponse, error) {
	paramsMap := map[string]interface{}{}
	paramsMap["mod"] = "post"
	paramsMap["act"] = "get"
	paramsMap["id"] = id

	resBody, err := doRequest(ctx, baseURL, http.MethodGet, token, paramsMap, nil)
	if err != nil {
		return model.GetArticleResponse{}, fmt.Errorf("get article error: %w", err)
	}

	resData := model.GetArticleResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.GetArticleResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}

func ListArticle(ctx context.Context, baseURL string, token string, req model.ListArticleRequest) (model.ListArticleResponse, error) {
	params, err := json.Marshal(req)
	if err != nil {
		return model.ListArticleResponse{}, fmt.Errorf("marshal error: %w", err)
	}

	paramsMap := map[string]interface{}{}

	err = json.Unmarshal(params, &paramsMap)
	if err != nil {
		return model.ListArticleResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	paramsMap[ParamMod] = ModPost
	paramsMap[ParamAct] = ActList

	resBody, err := doRequest(ctx, baseURL, http.MethodGet, token, paramsMap, nil)
	if err != nil {
		return model.ListArticleResponse{}, fmt.Errorf("list article error: %w", err)
	}

	resData := model.ListArticleResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.ListArticleResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}

func DeleteArticle(ctx context.Context, baseURL string, token string, id string) error {
	paramsMap := map[string]interface{}{}
	paramsMap["mod"] = "post"
	paramsMap["act"] = "delete"
	paramsMap["id"] = id

	resBody, err := doRequest(ctx, baseURL, http.MethodGet, token, paramsMap, nil)
	if err != nil {
		return fmt.Errorf("delete article error: %w", err)
	}

	resData := model.DeleteArticleResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	return nil
}
