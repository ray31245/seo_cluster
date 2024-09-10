package origin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ray31245/seo_cluster/pkg/z_blog_api/model"
)

func PostComment(ctx context.Context, baseURL string, token string, comment model.PostCommentRequest) error {
	bytesData, err := json.Marshal(comment)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	resBody, err := doRequest(ctx, baseURL, http.MethodPost, token, map[string]interface{}{ParamMod: ModComment, ParamAct: ActPost}, bytesData)
	if err != nil {
		return fmt.Errorf("post comment error: %w", err)
	}

	resData := model.PostCommentResponse{}
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	return nil
}
