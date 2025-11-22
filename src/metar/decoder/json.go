// Package decoder
package decoder

import (
	"encoding/json"
	"strings"

	"github.com/mdaverde/jsonpath"
)

type JsonDecoder struct {
}

func (j *JsonDecoder) Decode(raw []byte, selector string, reverse bool, multiline string) (bool, string, error) {
	var data interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return false, "", err
	}
	value, err := jsonpath.Get(&data, selector)
	if err != nil {
		return false, "", err
	}
	// 如果value是数组
	if valueArr, ok := value.([]interface{}); ok {
		// 如果valueArr的长度是0
		if len(valueArr) == 0 {
			return false, "", nil
		}
		// 如果valueArr的元素不是string
		if _, ok := valueArr[0].(string); !ok {
			return false, "", nil
		}
		// 如果reverse为true
		if reverse {
			return true, valueArr[len(valueArr)-1].(string), nil
		}
		return true, valueArr[0].(string), nil
	}
	// 如果value是字符串
	if valueStr, ok := value.(string); ok {
		// 如果multiline为空
		if multiline == "" {
			return true, valueStr, nil
		}
		lines := strings.Split(valueStr, multiline)
		if reverse {
			return true, lines[len(lines)-1], nil
		}
		return true, lines[0], nil
	}
	// 如果value不是数组也不是字符串
	return false, "", nil
}
