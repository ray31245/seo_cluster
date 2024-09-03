package aiassist

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	aiassistinterface "github.com/ray31245/seo_cluster/pkg/ai_assist/ai_assist_interface"
	"github.com/ray31245/seo_cluster/pkg/ai_assist/model"
	"google.golang.org/api/option"
)

var _ aiassistinterface.AIAssistInterface = &AIAssist{}

type AIAssist struct {
	token    string
	client   *genai.Client
	rewriter *genai.GenerativeModel
}

func NewAIAssist(ctx context.Context, token string) (*AIAssist, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(token))
	if err != nil {
		return nil, fmt.Errorf("failed to create new client: %w", err)
	}

	// The Gemini 1.5 models are versatile and work with both text-only and multimodal prompts
	rewriter := client.GenerativeModel("gemini-1.5-flash")
	rewriter.GenerationConfig = genai.GenerationConfig{
		ResponseMIMEType: "application/json",
	}

	return &AIAssist{
		token:    token,
		client:   client,
		rewriter: rewriter,
	}, nil
}

func (a *AIAssist) Close() error {
	return a.client.Close()
}

func (a *AIAssist) Rewrite(ctx context.Context, text string) (model.RewriteResponse, error) {
	//nolint:gosmopolitan // prompt is a string
	prompt := "你是一位收悉區塊鏈的專欄作家，請你將以下內容用你的話重新闡述文章中的內容，並訂一個標題。請使用json格式輸出：{Title: string,Content: string}"

	resp, err := a.rewriter.GenerateContent(ctx, genai.Text(fmt.Sprintf("%s\n%s", prompt, text)))
	if err != nil {
		return model.RewriteResponse{}, fmt.Errorf("failed to rewrite content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return model.RewriteResponse{}, errors.New("no content generated")
	}

	res := model.RewriteResponse{}

	err = json.Unmarshal([]byte(fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])), &res)
	if err != nil {
		return model.RewriteResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return res, nil
}
