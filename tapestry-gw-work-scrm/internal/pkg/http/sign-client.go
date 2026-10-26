package http

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/json"
	"io/ioutil"
	"strings"
)

// ------------------ sign 签名算法 ----------------------
type Sign interface {
	// TODO
	Sign(in string) string
}

type SignFunc func(string) string

func (s SignFunc) Sign(in string) string {
	return s(in)
}

type MD5Sign struct {
}

func (M *MD5Sign) Sign(in string) string {
	m := md5.New()
	m.Write([]byte(in))
	return fmt.Sprintf("%s", hex.EncodeToString(m.Sum(nil)))
}

type ShA256Sign struct {
}

func (M *ShA256Sign) Sign(in string) string {
	m := sha256.New()
	m.Write([]byte(in))
	return fmt.Sprintf("%s", hex.EncodeToString(m.Sum(nil)))
}

// ------------------ sign 签名数据源获取 ----------------------
type SignOrigin interface {
	Origin(data interface{}) string
}

type JsonOrigin struct {
}

func (t *JsonOrigin) Origin(data interface{}) string {
	return json.MarshalToString(data)
}

type UrlOrigin struct {
}

func (t *UrlOrigin) Origin(data interface{}) string {
	return json.MarshalToString(data)
}

type SortSplitOrigin struct {
	s string
}

func NewSortSplitOrigin(s string) *SortSplitOrigin {
	return &SortSplitOrigin{s: s}
}

func (t *SortSplitOrigin) Origin(data interface{}) string {
	a := data.(map[string]string)
	var values []string
	for _, v := range a {
		values = append(values, v)
	}
	return strings.Join(values, t.s)
}

// ------------------ TransformClient 网络传输方式 ----------------------
type TransformClient interface {
	Send(string, []byte) ([]byte, error)
}

type PostTransformClient struct {
}

func (p PostTransformClient) Send(url string, bytes []byte) ([]byte, error) {
	httpClient := clients.NewHTTPClient()
	res, err := httpClient.Post(ctx.New(), url, clients.JSONBody(bytes))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(res.Response().Body)
}

// ------------------ Client 客户端使用代理 ----------------------
type Client struct {
	sign   Sign
	origin SignOrigin
	client TransformClient
}

func NewClient(sign Sign, origin SignOrigin, client TransformClient) *Client {
	return &Client{sign: sign, origin: origin, client: client}
}

func (c *Client) Request(url string, data interface{}) string {
	res, err := c.client.Send(url+"?sing="+c.sign.Sign(json.MarshalToString(data)), []byte(json.MarshalToString(data)))

	if err != nil {
		return ""
	}

	return string(res)
}

func (c *Client) Http(url string, method string, data interface{}) string {
	res, err := c.client.Send(url+"?sing="+c.sign.Sign(json.MarshalToString(data)), []byte(json.MarshalToString(data)))

	if err != nil {
		return ""
	}

	switch strings.ToUpper(method) {
	case "POSt":
		p := PostTransformClient{}
		p.Send(url, data.([]byte))
	case "GEt":
		//p := GetTransformClient{}
		//p.Send(url, data.([]byte))
	}

	return string(res)
}

func (c *Client) Post(url string, data interface{}) string {
	return c.Http(url, "POST", data)
}
