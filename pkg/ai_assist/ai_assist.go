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

const (
	maxTemperature = 2
)

type AIAssist struct {
	token     string
	client    *genai.Client
	rewriter  *genai.GenerativeModel
	commenter *genai.GenerativeModel
	evaluator *genai.GenerativeModel
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

	commenter := client.GenerativeModel("gemini-1.5-flash")
	commenter.GenerationConfig = genai.GenerationConfig{
		ResponseMIMEType: "application/json",
		Temperature: func() *float32 {
			temp := float32(maxTemperature)

			return &temp
		}(),
	}

	evaluator := client.GenerativeModel("gemini-1.5-flash")
	evaluator.GenerationConfig = genai.GenerationConfig{
		ResponseMIMEType: "application/json",
	}

	return &AIAssist{
		token:     token,
		client:    client,
		rewriter:  rewriter,
		commenter: commenter,
		evaluator: evaluator,
	}, nil
}

func (a *AIAssist) Close() error {
	err := a.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close client: %w", err)
	}

	return nil
}

func (a *AIAssist) Rewrite(ctx context.Context, text []byte) (model.RewriteResponse, error) {
	//nolint:gosmopolitan // prompt is a string
	prompt := "你是一位收悉区块链的简体中文语系专栏作家，请你将以下内容用你的话重新阐述文章中的内容，并订一个标题。请使用json格式输出：{Title: string,Content: string}"

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

func (a *AIAssist) Comment(ctx context.Context, text []byte) (model.CommentResponse, error) {
	//nolint:gosmopolitan // prompt is a string
	prompt := "你在网路上看到以下文章，请随性且简洁地在这篇文章下留言。并且以一位看新闻的人的角度记录这个文章能够为你提供的价值。最低0分滿分100分。请使用json格式输出：{Comment: string, Score: int}\n。"

	resp, err := a.rewriter.GenerateContent(ctx, genai.Text(fmt.Sprintf("%s\n%s", prompt, text)))
	if err != nil {
		return model.CommentResponse{}, fmt.Errorf("failed to comment content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return model.CommentResponse{}, errors.New("no content generated")
	}

	res := model.CommentResponse{}

	err = json.Unmarshal([]byte(fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])), &res)
	if err != nil {
		return model.CommentResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return res, nil
}

func (a *AIAssist) Evaluate(ctx context.Context, text []byte) (model.EvaluateResponse, error) {
	//nolint:gosmopolitan // prompt is a string
	prompt := "你是一位区块链专栏作家，你的文章被编辑修改过，你需要评价这篇文章的质量。最低0分滿分100分。请使用json格式输出：{ Score: int}。"

	resp, err := a.rewriter.GenerateContent(ctx, genai.Text(fmt.Sprintf("%s\n%s", prompt, text)))
	if err != nil {
		return model.EvaluateResponse{}, fmt.Errorf("failed to evaluate content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return model.EvaluateResponse{}, errors.New("no content generated")
	}

	res := model.EvaluateResponse{}

	err = json.Unmarshal([]byte(fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])), &res)
	if err != nil {
		return model.EvaluateResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return res, nil
}
