package tool

import "bytes"

// PKCS5Padding 末尾填充字节
func PKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize // 要填充的值和个数
	slice1 := []byte{byte(padding)}            // 要填充的单个二进制值
	slice2 := bytes.Repeat(slice1, padding)    // 要填充的二进制数组
	return append(data, slice2...)             // 填充到数据末端
}

func PKCS5UnPadding(data []byte) []byte {
	unpadding := data[len(data)-1]                // 获取二进制数组最后一个数值
	result := data[:(len(data) - int(unpadding))] // 截取开始至总长度减去填充值之间的有效数据
	return result
}

// ZerosPadding 末尾填充0
func ZerosPadding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize // 要填充的个数
	slice1 := []byte{0}                        // 要填充的单个0数据
	slice2 := bytes.Repeat(slice1, padding)    // 要填充的0二进制数组
	return append(data, slice2...)             // 填充到数据末端
}

func ZerosUnPadding(data []byte) []byte {
	return bytes.TrimRightFunc(data, func(r rune) bool { // 去除满足条件的子切片
		return r == 0
	})
}
