package zBlogApi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goTool/pkg/util"
	"goTool/pkg/z_blog_api/model"
	"io"
	"net/http"
)

func (t *ZBlogAPIClient) postArticle(art model.PostArticleRequest) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	requestUrl := t.baseURL
	values := requestUrl.Query()
	values.Add("mod", "post")
	values.Add("act", "post")
	requestUrl.RawQuery = values.Encode()
	bytesData, err := util.EscapeHTMLMarshal(art)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, requestUrl.String(), bytes.NewReader(bytesData))
	if err != nil {
		return fmt.Errorf("postArticle: new request error: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("post article error: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read body error: %w", err)
	}
	resData := model.PostArticleResponse{}
	if err := json.Unmarshal(body, &resData); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}
	if resData.Code != 200 {
		return fmt.Errorf("post article error: %s", resData.Message)
	}
	return nil
}

func (t *ZBlogAPIClient) listArticle(model.ListArticleRequest) (model.ListArticleResponse, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	requestUrl := t.baseURL
	values := requestUrl.Query()
	values.Add("mod", "post")
	values.Add("act", "list")
	requestUrl.RawQuery = values.Encode()
	req, err := http.NewRequest(http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return model.ListArticleResponse{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return model.ListArticleResponse{}, fmt.Errorf("list article error: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return model.ListArticleResponse{}, fmt.Errorf("status code error: %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return model.ListArticleResponse{}, fmt.Errorf("read body error: %w", err)
	}
	resData := model.ListArticleResponse{}
	if err := json.Unmarshal(body, &resData); err != nil {
		return model.ListArticleResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}
	if resData.Code != 200 {
		return model.ListArticleResponse{}, fmt.Errorf("list article error: %s", resData.Message)
	}
	return resData, nil
}

func (t *ZBlogAPIClient) deleteArticle(id string) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	requestUrl := t.baseURL
	values := requestUrl.Query()
	values.Add("mod", "post")
	values.Add("act", "delete")
	values.Add("id", id)
	requestUrl.RawQuery = values.Encode()
	req, err := http.NewRequest(http.MethodPost, requestUrl.String(), nil)
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
	resData := model.DeleteArticleResponse{}
	if err := json.Unmarshal(body, &resData); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}
	if resData.Code != 200 {
		return fmt.Errorf("delete article error: %s", resData.Message)
	}
	return nil
}
