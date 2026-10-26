package common

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
)

// Md5 - Calculate the md5 hash of a string
func Md5(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

// Base64Encode - Encodes data with MIME base64
func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// Base64Decode - Decodes data encoded with MIME base64
func Base64Decode(s string) (string, error) {
	bs, err := base64.StdEncoding.DecodeString(s)
	return string(bs), err
}

// JsonEncode - Returns the JSON representation of a value
func JsonEncode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// JsonDecode - Decodes a JSON string
func JsonDecode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
