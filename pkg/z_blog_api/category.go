package zblogapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (t *ZblogAPIClient) listCategory() (ListCategoryResponse, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	requestUrl := t.baseURL
	values := requestUrl.Query()
	values.Add("mod", "category")
	values.Add("act", "list")
	requestUrl.RawQuery = values.Encode()
	req, err := http.NewRequest(http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return ListCategoryResponse{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return ListCategoryResponse{}, fmt.Errorf("list category error: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return ListCategoryResponse{}, fmt.Errorf("status code error: %d", res.StatusCode)
	}
	var resData ListCategoryResponse
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		return ListCategoryResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	if resData.Code != 200 {
		return ListCategoryResponse{}, fmt.Errorf("list category error: %s", resData.Message)
	}
	return resData, nil
}
