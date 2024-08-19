package zblogapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	zBlogErr "github.com/ray31245/seo_cluster/pkg/z_blog_api/error"
	"github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
)

func (t *Client) retry(ctx context.Context, f func() error) error {
	err := f()
	if err != nil {
		err = t.Login(ctx)
		if err != nil {
			return fmt.Errorf("login error: %w", err)
		}

		err = f()
		if err != nil {
			return fmt.Errorf("retry error: %w", err)
		}
	}

	return nil
}

func (t *Client) requestWithBlock(ctx context.Context, method string, parameter map[string]interface{}, body []byte) ([]byte, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	requestURL := t.baseURL
	values := requestURL.Query()

	for k, v := range parameter {
		values.Add(k, fmt.Sprintf("%v", v))
	}

	requestURL.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, method, requestURL.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("new request error: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+t.token)

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read body error: %w", err)
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
