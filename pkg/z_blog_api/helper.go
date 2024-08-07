package zblogapi

func (t *ZblogAPI) retry(f func() error) error {
	err := f()
	if err != nil {
		t.Login()
		err = f()
		if err != nil {
			return err
		}
	}
	return nil
}
