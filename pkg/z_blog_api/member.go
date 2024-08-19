package zblogapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"goTool/pkg/z_blog_api/model"
)

func (t *Client) listMember(ctx context.Context) (model.ListMemberResponse, error) {
	resBopdy, err := t.requestWithBlock(ctx, http.MethodGet, map[string]interface{}{"mod": "member", "act": "list"}, nil)
	if err != nil {
		return model.ListMemberResponse{}, fmt.Errorf("list member error: %w", err)
	}

	resData := model.ListMemberResponse{}
	if err := json.Unmarshal(resBopdy, &resData); err != nil {
		return model.ListMemberResponse{}, fmt.Errorf("unmarshal error: %w", err)
	}

	return resData, nil
}
