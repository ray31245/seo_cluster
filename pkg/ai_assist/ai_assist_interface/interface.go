package aiassistinterface

import (
	"context"

	"github.com/ray31245/seo_cluster/pkg/ai_assist/model"
)

type AIAssistInterface interface {
	CustomRewrite(ctx context.Context, systemPrompt string, prompt string, content []byte) (string, error)
	Rewrite(ctx context.Context, text []byte) (model.RewriteResponse, error)
	ExtendRewrite(ctx context.Context, text []byte) (model.ExtendRewriteResponse, error)
	MultiSectionsRewrite(ctx context.Context, systemPrompt string, content string) (string, error)
	Comment(ctx context.Context, text []byte) (model.CommentResponse, error)
	FindKeyWords(ctx context.Context, text []byte) (model.FindKeyWordsResponse, error)
	SelectCategory(ctx context.Context, req model.SelectCategoryRequest) (model.SelectCategoryResponse, error)
	MakeTitle(ctx context.Context, systemPrompt string, prompt string, content []byte) (string, error)
	Lock()
	Unlock()
	TryLock() bool
}
