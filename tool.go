package tool

import "math/rand"

// RandomString 返回一个随机字符串，包含数字，大写字母，小写字母
func RandomString(len int) []byte {
	b := make([]byte, len)
	for i := 0; i < len; i++ {
		x := rand.Intn(3)
		switch x {
		case 0:
			b[i] = byte(rand.Intn(26) + 65)
		case 1:
			b[i] = byte(rand.Intn(26) + 97)
		case 2:
			b[i] = byte(rand.Intn(10) + 48)
		}
	}

	return b
}

// RandomStringHEX 返回一个随机16进制字符串（已转换为ASCII表示）
func RandomStringHEX(len int) []byte {
	b := make([]byte, len)
	for i := 0; i < len; i++ {
		x := rand.Intn(2)
		switch x {
		case 0:
			b[i] = byte(rand.Intn(6) + 97)
		default:
			b[i] = byte(rand.Intn(10) + 48)
		}
	}
	return b
}
