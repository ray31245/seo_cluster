package util

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return fmt.Errorf("UnixTime.UnmarshalJSON: json unmarshal error: %w", err)
	}

	t.Time = time.Unix(timestamp, 0)

	return nil
}
