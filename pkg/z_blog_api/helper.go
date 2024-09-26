package zblogapi

import (
	"context"
	"errors"
	"fmt"

	zBlogErr "github.com/ray31245/seo_cluster/pkg/z_blog_api/error"
)

func (t *Client) retry(ctx context.Context, f func() error) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	err := f()
	if errors.Is(err, zBlogErr.ErrHTTPUnauthorized) || errors.Is(err, zBlogErr.ErrHTTPForbidden) || errors.Is(err, zBlogErr.ErrIllegalAccess) {
		err = t.Login(ctx)
		if err != nil {
			return fmt.Errorf("retry error: %w", err)
		}

		err = f()
		if err != nil {
			return fmt.Errorf("retry error: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("retry error: %w", err)
	}

	return nil
}
