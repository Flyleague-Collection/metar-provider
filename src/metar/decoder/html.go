// Package decoder
package decoder

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type HtmlDecoder struct{}

func (h *HtmlDecoder) Decode(raw []byte, selector string, reverse bool, multiline string) (bool, string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(raw))
	if err != nil {
		return false, "", err
	}
	data := doc.Find(selector)
	// 如果没有选取到元素
	if data == nil || data.Length() == 0 {
		return false, "", nil
	}
	// 如果只选取到了一个元素
	if data.Length() == 1 {
		text := data.Get(0).FirstChild.Data
		// 如果multiline为空
		if multiline == "" {
			return true, text, nil
		}
		lines := strings.Split(text, multiline)
		if reverse {
			return true, lines[len(lines)-1], nil
		}
		return true, lines[0], nil
	}
	return false, "", nil
}
