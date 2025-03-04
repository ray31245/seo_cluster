package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
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

func HTMLToMd(htmlCode string) (string, error) {
	md, err := htmltomarkdown.ConvertString(htmlCode)
	if err != nil {
		return "", fmt.Errorf("HTMLToMd: convert error: %w", err)
	}

	return md, nil
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

const TimeFormat = "2006-01-02 15:04:05"

// MarshalJSON implements the json.Marshaler interface.
// ref: https://github.com/gin-gonic/gin/issues/2479#issuecomment-1502705948
func (t UnixTime) MarshalJSON() ([]byte, error) {
	return EncodeFormatedTime(t.Time)
}

func EncodeFormatedTime(t time.Time) ([]byte, error) {
	if time.Time.IsZero(t) {
		return []byte("\"\""), nil
	}
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, TimeFormat)
	b = append(b, '"')

	return b, nil
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

func (i NumberString) String() string {
	return string(i)
}

func (i NumberString) Int() int {
	n, err := strconv.Atoi(i.String())
	if err != nil {
		return 0
	}

	return n
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}

	return string(b)
}

func WaitGroupChan(wg *sync.WaitGroup) chan bool {
	done := make(chan bool)

	go func() {
		wg.Wait()

		done <- true
	}()

	return done
}
