package util

import (
	"bytes"
	"encoding/json"
)

func EscapeHTMLMarshual(art interface{}) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	if err := jsonEncoder.Encode(art); err != nil {
		return nil, err
	}

	return bytes.TrimSuffix(bf.Bytes(), []byte{'\n'}), nil
}
