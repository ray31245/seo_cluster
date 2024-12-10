package aiassistinterface

import (
	"context"

	"github.com/ray31245/seo_cluster/pkg/ai_assist/model"
)

type AIAssistInterface interface {
	Rewrite(ctx context.Context, text []byte) (model.RewriteResponse, error)
	Comment(ctx context.Context, text []byte) (model.CommentResponse, error)
	FindKeyWords(ctx context.Context, text []byte) (model.FindKeyWordsResponse, error)
	SelectCategory(ctx context.Context, req model.SelectCategoryRequest) (model.SelectCategoryResponse, error)
	Lock()
	Unlock()
	TryLock() bool
}
