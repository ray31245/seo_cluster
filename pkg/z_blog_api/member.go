package zblogapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (t *ZblogAPI) listMember() (string, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	values := t.baseURL.Query()
	values.Add("mod", "member")
	values.Add("act", "list")
	t.baseURL.RawQuery = values.Encode()
	req, err := http.NewRequest(http.MethodGet, t.baseURL.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("read body error: %w", err)
	}
	resData := ListMemberResponse{}
	if err := json.Unmarshal(body, &resData); err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}
	if resData.Code != 200 {
		return "", fmt.Errorf("list member error: %s", resData.Message)
	}
	return string(body), nil
}
