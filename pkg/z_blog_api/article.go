package zblogapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"goTool/pkg/util"
	"goTool/pkg/z_blog_api/model"
)

func (t *Client) postArticle(ctx context.Context, art model.PostArticleRequest) error {
	bytesData, err := util.EscapeHTMLMarshal(art)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	resBody, err := t.requestWithBlock(ctx, http.MethodPost, map[string]interface{}{"mod": "post", "act": "post"}, bytesData)
	if err != nil {
		return fmt.Errorf("post article error: %w", err)
	}

	resData := model.PostArticleResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	return nil
}

func (t *Client) listArticle(ctx context.Context, req model.ListArticleRequest) (model.ListArticleResponse, error) {
	params, err := json.Marshal(req)
	if err != nil {
		return model.ListArticleResponse{}, fmt.Errorf("marshal error: %w", err)
	}

	paramsMap := map[string]interface{}{}

	err = json.Unmarshal(params, &paramsMap)
	if err != nil {
		return model.ListArticleResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	paramsMap["mod"] = "post"
	paramsMap["act"] = "list"

	resBody, err := t.requestWithBlock(ctx, http.MethodGet, paramsMap, nil)
	if err != nil {
		return model.ListArticleResponse{}, fmt.Errorf("list article error: %w", err)
	}

	resData := model.ListArticleResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.ListArticleResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}

func (t *Client) deleteArticle(ctx context.Context, id string) error {
	resBody, err := t.requestWithBlock(ctx, http.MethodPost, map[string]interface{}{"mod": "post", "act": "delete", "id": id}, nil)
	if err != nil {
		return fmt.Errorf("delete article error: %w", err)
	}

	resData := model.DeleteArticleResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	return nil
}
