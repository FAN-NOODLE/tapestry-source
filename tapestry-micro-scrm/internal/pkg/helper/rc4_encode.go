package helper

import (
	"crypto/rc4"
	"encoding/base64"
	"strconv"
)

// SecretKey 密钥
const SecretKey = "REDACTED_FOR_SOURCE_SHARING"

// Decode 消息解密
func Decode(key string, b []byte) []byte {
	dest := make([]byte, len(b))
	cipher2, _ := rc4.NewCipher([]byte(key))
	cipher2.XORKeyStream(dest, b)
	return dest
}

// Encode 消息加密
func Encode(key string, src []byte) []byte {
	dest := make([]byte, len(src))
	cipher1, _ := rc4.NewCipher([]byte(key))
	cipher1.XORKeyStream(dest, src)
	return dest
}

// EncodeToString returns the base64 encoding of src.
func EncodeToString(str string) string {
	// 要加密的源数据
	src := []byte(str)
	// 加密字符串，返回值是加密后的字符串
	return base64.StdEncoding.EncodeToString(Encode(SecretKey, src))
}

// DecodeString returns the bytes represented by the base64 string s.
func DecodeString(s string) (string, error) {
	// 解密字符串，对加密后的字符串进行解密，返回值是解密后的数据（byte数组形式）
	decodeDataByteArr, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return s, err
	}
	return string(Decode(SecretKey, decodeDataByteArr)), nil
}

// GetBirthdayArr 生日区间加密后的数组
func GetBirthdayArr(birthdayStart, birthdayEnd int) []string {

	var birthdayArr []string
	for i := birthdayStart; i <= birthdayEnd; i++ {
		var birthdayMonth string
		birthdayMonth = strconv.Itoa(i)
		if i <= 9 {
			birthdayMonth = "0" + birthdayMonth
		}
		birthdayArr = append(birthdayArr, EncodeToString(birthdayMonth))
	}
	return birthdayArr
}
