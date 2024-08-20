package zblogapi

import (
	"context"
	"fmt"
)

func (t *Client) retry(ctx context.Context, f func() error) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	err := f()
	if err != nil {
		err = t.Login(ctx)
		if err != nil {
			return fmt.Errorf("login error: %w", err)
		}

		err = f()
		if err != nil {
			return fmt.Errorf("retry error: %w", err)
		}
	}

	return nil
}
