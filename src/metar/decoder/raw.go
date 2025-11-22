// Package decoder
package decoder

import (
	"strings"
)

type RawDecoder struct {
}

func (r *RawDecoder) Decode(raw []byte, _ string, reverse bool, multiline string) (bool, string, error) {
	data := string(raw)
	// 如果mutiline不为空，则切分数据
	if multiline != "" {
		dataLines := strings.Split(data, multiline)
		// 如果reverse为true，则取最后一行数据返回
		if reverse {
			return true, dataLines[len(dataLines)-1], nil
		}
		return true, dataLines[0], nil
	}
	// 否则直接返回数据
	return true, data, nil
}
