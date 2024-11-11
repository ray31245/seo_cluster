package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func EscapeHTMLMarshal(art interface{}) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)

	if err := jsonEncoder.Encode(art); err != nil {
		return nil, fmt.Errorf("EscapeHTMLMarshal: json encode error: %w", err)
	}

	return bytes.TrimSuffix(bf.Bytes(), []byte{'\n'}), nil
}

func MdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

type UnixTime struct {
	time.Time
}

// MarshalJSON implements the json.Marshaler interface.
func (t *UnixTime) UnmarshalJSON(data []byte) error {
	var timestamp StringNumber
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return fmt.Errorf("UnixTime.UnmarshalJSON: json unmarshal error: %w", err)
	}

	t.Time = time.Unix(int64(timestamp), 0)

	return nil
}

type StringNumber int

func (i *StringNumber) UnmarshalJSON(data []byte) error {
	var str interface{}
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("StringNumber.UnmarshalJSON: json unmarshal error: %w", err)
	}

	switch v := str.(type) {
	case string:
		number, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("StringNumber.UnmarshalJSON: strconv.Atoi error: %w", err)
		}

		*i = StringNumber(number)
	case float64:
		*i = StringNumber(int(v))
	case int:
		*i = StringNumber(v)
	default:
		return fmt.Errorf("StringNumber.UnmarshalJSON: unknown type: %T", v)
	}

	return nil
}

type NumberString string

// MarshalJSON implements the json.Marshaler interface.
func (i *NumberString) UnmarshalJSON(data []byte) error {
	var str interface{}
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("NumberString.UnmarshalJSON: json unmarshal error: %w", err)
	}

	switch v := str.(type) {
	case string:
		*i = NumberString(v)
	case float64:
		*i = NumberString(fmt.Sprintf("%d", int(v)))
	case int:
		*i = NumberString(fmt.Sprintf("%d", v))
	default:
		return fmt.Errorf("NumberString.UnmarshalJSON: unknown type: %T", v)
	}

	return nil
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}

	return string(b)
}
