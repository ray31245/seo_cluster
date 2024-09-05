package util

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ListImageSrcFromHtml(body []byte) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var images []string

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		// validate src is a reachable URL for image
		// check if the response content type is image
		// ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types
		if resp, err := http.Get(src); err == nil && strings.HasPrefix(resp.Header.Get("Content-Type"), "image") {
			images = append(images, src)
		}
	})

	return images, nil
}

func GenImageListEncodeDiv(body []byte) (string, error) {
	images, err := ListImageSrcFromHtml(body)
	if err != nil {
		return "", fmt.Errorf("GenImageListEncodeDiv: ListImageSrcFromHtml error: %w", err)
	}

	var imageDiv string
	for _, image := range images {
		imageDiv = imageDiv + fmt.Sprintf("<img src=\"%s\" />", image)
	}

	// encode imageDiv to base64
	imageDiv = base64.StdEncoding.EncodeToString([]byte(imageDiv))

	// wrap imageDiv with div tag
	imageDiv = fmt.Sprintf("<div class=\"encodeImageList\">%s</div>", imageDiv)

	return imageDiv, nil
}

func DecodeImageListDivFromHTMl(body []byte) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	var imageDiv string

	doc.Find("div.encodeImageList").Each(func(i int, s *goquery.Selection) {
		imageDiv = s.Text()
	})

	// decode imageDiv from base64
	decodeImageDiv, err := base64.StdEncoding.DecodeString(imageDiv)
	if err != nil {
		return "", fmt.Errorf("DecodeImageListDivFromHTMl: base64 decode error: %w", err)
	}

	doc.Find("div.encodeImageList").Each(func(i int, s *goquery.Selection) {
		s.ReplaceWithHtml(string(decodeImageDiv))
	})

	res := ""

	var callBackErr error

	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		// log.Println(s.Text())
		res, err = s.Html()
		if err != nil {
			callBackErr = fmt.Errorf("DecodeImageListDivFromHTMl: set text error: %w", err)
		}
	})

	if callBackErr != nil {
		return "", callBackErr
	}

	// res, err := doc.Html()
	// if err != nil {
	// 	return "", fmt.Errorf("DecodeImageListDivFromHTMl: set text error: %w", err)
	// }

	return res, nil
}
