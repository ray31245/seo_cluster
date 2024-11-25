package origin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	wordpressError "github.com/ray31245/seo_cluster/pkg/wordpress_api/error"
	"github.com/ray31245/seo_cluster/pkg/wordpress_api/model"
)

const (
	APIPath = "wp-json/wp/v2"
)

func doRequest(ctx context.Context, baseURL string, method string, route string, basicAuth model.BasicAuthentication, parameter map[string]interface{}, body []byte) ([]byte, error) {
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse url error: %w", err)
	}

	reqURL = reqURL.JoinPath(APIPath)
	reqURL = reqURL.JoinPath(route)

	values := reqURL.Query()
	for k, v := range parameter {
		values.Add(k, fmt.Sprintf("%v", v))
	}

	reqURL.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("new request error: %w", err)
	}

	if !basicAuth.IsAnonymous {
		req.SetBasicAuth(basicAuth.Username, basicAuth.Password)
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

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

	// log.Println(string(resBody))

	statusCodeErr := wordpressError.NewHTTPStatusCodeError(res.StatusCode)
	if statusCodeErr != nil {
		errRes := model.ErrorResponse{}
		if err := json.Unmarshal(resBody, &errRes); err != nil {
			return nil, fmt.Errorf("request error: %w with message: %s", statusCodeErr, resBody)
		}

		// TODO: decode error message form utf-8 to string
		return nil, fmt.Errorf("request error: %w with message: %s", statusCodeErr, resBody)
	}

	return resBody, nil
}
