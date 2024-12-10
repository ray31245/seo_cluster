package origin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	zBlogErr "github.com/ray31245/seo_cluster/pkg/z_blog_api/error"
	"github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
)

const (
	APIPath = "zb_system/api.php"

	ModMember   = "member"
	ModPost     = "post"
	ModCategory = "category"
	ModComment  = "comment"
	ModTag      = "tag"

	ActGet    = "get"
	ActList   = "list"
	ActLogin  = "login"
	ActPost   = "post"
	ActDelete = "delete"

	ParamMod = "mod"
	ParamAct = "act"
)

func doRequest(ctx context.Context, baseURL string, method string, token string, parameter map[string]interface{}, body []byte) ([]byte, error) {
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse url error: %w", err)
	}

	reqURL = reqURL.JoinPath(APIPath)

	values := reqURL.Query()
	for k, v := range parameter {
		values.Add(k, fmt.Sprintf("%v", v))
	}

	reqURL.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("new request error: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error: %w", err)
	}

	err = zBlogErr.NewHTTPStatusCodeError(res.StatusCode)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	resData := model.BasicResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	err = zBlogErr.NewHTTPStatusCodeErrWithMsg(resData.Code, resData.Message)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	return resBody, nil
}
