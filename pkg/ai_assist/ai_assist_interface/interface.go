package aiassistinterface

import (
	"context"

	"github.com/ray31245/seo_cluster/pkg/ai_assist/model"
)

type AIAssistInterface interface {
	Rewrite(ctx context.Context, text []byte) (model.RewriteResponse, error)
}
