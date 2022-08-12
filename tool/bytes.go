package tool

import (
	"bytes"
	"encoding/hex"
	"strings"
)

// 字节工具包

// BytesToHex 字节转换Hex字符串
// padding[0]: 左填充
// padding[1]：右填充(最后一个字节不填充)
func BytesToHex(data []byte, padding ...string) string {

	hexText := bytesToHex(data)

	if len(padding) > 0 {
		buffer := new(bytes.Buffer)
		raw := []byte(hexText)

		prefix := padding[0]
		suffix := ""
		if len(padding) > 1 {
			suffix = padding[1]
		}
		size := len(data)
		if len(hexText)/2 != len(data) { // hex格式应与byte格式长度匹配
			return ""
		}
		for i := 0; i < size; i++ {
			buffer.WriteString(prefix)
			buffer.Write(raw[2*i+0 : 2*i+2])
			if i+1 < size {
				buffer.WriteString(suffix)
			}
		}
		return buffer.String()
	} else {
		return hexText
	}
}

func bytesToHex(data []byte) string {
	return strings.ToLower(hex.EncodeToString(data))
}
