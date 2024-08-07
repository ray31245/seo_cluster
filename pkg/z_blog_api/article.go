package zblogapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goTool/pkg/util"
	"io"
	"net/http"
)

func (t *ZblogAPI) postArticle(art PostArticleRequest) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	values := t.baseURL.Query()
	values.Add("mod", "post")
	values.Add("act", "post")
	t.baseURL.RawQuery = values.Encode()
	bytesData, err := util.EscapeHTMLMarshual(art)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, t.baseURL.String(), bytes.NewReader(bytesData))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read body error: %w", err)
	}
	resData := PostArticleResponse{}
	if err := json.Unmarshal(body, &resData); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}
	if resData.Code != 200 {
		return fmt.Errorf("post article error: %s", resData.Message)
	}
	return nil
}

func (t *ZblogAPI) listArticle(ListArticleRequest) (ListArticleResponse, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	values := t.baseURL.Query()
	values.Add("mod", "post")
	values.Add("act", "list")
	t.baseURL.RawQuery = values.Encode()
	req, err := http.NewRequest(http.MethodGet, t.baseURL.String(), nil)
	if err != nil {
		return ListArticleResponse{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return ListArticleResponse{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return ListArticleResponse{}, fmt.Errorf("status code error: %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ListArticleResponse{}, fmt.Errorf("read body error: %w", err)
	}
	resData := ListArticleResponse{}
	if err := json.Unmarshal(body, &resData); err != nil {
		return ListArticleResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}
	if resData.Code != 200 {
		return ListArticleResponse{}, fmt.Errorf("list article error: %s", resData.Message)
	}
	return resData, nil
}
