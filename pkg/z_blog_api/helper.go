package zBlogApi

import "fmt"

func (t *ZBlogAPIClient) retry(f func() error) error {
	err := f()
	if err != nil {
		err = t.Login()
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
