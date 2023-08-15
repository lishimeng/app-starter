package tool

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	mathRand "math/rand"
	"strings"
	"time"
)

// 字符串工具包
//

const letters = "qwertyuipasdfghjkzxcvbnmQWERTYUIPASDFGHJKLZXCVBNM123456789!@#$%&[]"

var src = mathRand.NewSource(time.Now().UnixNano())
var size = len(letters)

// RandStr 随机字符串,返回n个字符
func RandStr(n int) (s string) {
	var count = 0
	buf := bytes.Buffer{}
	for {
		bs := Int64ToBytes(src.Int63())
		for _, b := range bs {
			i := b & 0xff
			index := int(i) % size
			buf.WriteByte(letters[index])
			count++
			if count >= n {
				break
			}
		}
		if count >= n {
			break
		}
	}
	s = buf.String()
	return
}

// Int64ToBytes int64-->8byte
func Int64ToBytes(n int64) (bs []byte) {
	bs = make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(n))
	return
}

// RandHexStr 随机字符串,hex表示.
//
// n:字节数,返回2n个hex字符
func RandHexStr(n int) (s string) {
	randBytes := make([]byte, n)
	_, _ = rand.Read(randBytes)
	s = fmt.Sprintf("%x", randBytes)
	return
}

func UUIDString() (s string) {
	u, err := uuid.NewRandom()
	if err != nil {
		return
	}
	s = u.String()
	s = strings.ToLower(strings.ReplaceAll(s, "-", ""))
	return
}

func Join(delimiter string, s ...string) string {
	return strings.Join(s, delimiter)
}
