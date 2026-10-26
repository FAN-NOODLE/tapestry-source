package http

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

type TestTransformClient struct {
}

func (t *TestTransformClient) Send(s string, bytes []byte) ([]byte, error) {
	return []byte(s + string(bytes)), nil
}

var aa SignFunc = func(string2 string) string {
	m := md5.New()
	m.Write([]byte(string2))
	return fmt.Sprintf("%s", hex.EncodeToString(m.Sum(nil)))
}

func TestClient_Request(t *testing.T) {

	c := NewClient(aa, NewSortSplitOrigin("|"), new(TestTransformClient))
	res := c.Request("http://xxx.com", "xxx")
	assert.Equal(t, "http://xxx.com?sing=9a345c07c276088126cc92a22af24b97\"xxx\"", res)
}
