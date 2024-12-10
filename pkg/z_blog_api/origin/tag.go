package origin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ray31245/seo_cluster/pkg/util"
	"github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
)

func ListTag(ctx context.Context, baseURL string, token string, req model.ListTagRequest) (model.ListTagResponse, error) {
	params, err := json.Marshal(req)
	if err != nil {
		return model.ListTagResponse{}, fmt.Errorf("marshal error: %w", err)
	}

	paramsMap := map[string]interface{}{}

	err = json.Unmarshal(params, &paramsMap)
	if err != nil {
		return model.ListTagResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	paramsMap[ParamMod] = ModTag
	paramsMap[ParamAct] = ActList

	resBody, err := doRequest(ctx, baseURL, http.MethodGet, token, paramsMap, nil)
	if err != nil {
		return model.ListTagResponse{}, fmt.Errorf("list tag error: %w", err)
	}

	var resData model.ListTagResponse
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.ListTagResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}

func PostTag(ctx context.Context, baseURL string, token string, req model.PostTagRequest) (model.PostTagResponse, error) {
	bytesData, err := util.EscapeHTMLMarshal(req)
	if err != nil {
		return model.PostTagResponse{}, fmt.Errorf("marshal error: %w", err)
	}

	resBody, err := doRequest(ctx, baseURL, http.MethodPost, token, map[string]interface{}{ParamMod: ModTag, ParamAct: ActPost}, bytesData)
	if err != nil {
		return model.PostTagResponse{}, fmt.Errorf("post tag error: %w", err)
	}

	resData := model.PostTagResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return model.PostTagResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}
