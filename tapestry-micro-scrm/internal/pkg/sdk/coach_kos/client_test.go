package coachkos

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

func TestClientList成功(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("请求方法 = %q", r.Method)
		}
		if r.URL.Path != "/api/api/biz/coach_mark_post_results" {
			t.Fatalf("请求路径 = %q", r.URL.Path)
		}
		if got := r.Header.Get("X-Coach-Mark-Token"); got != "token" {
			t.Fatalf("Token = %q", got)
		}
		if got := r.URL.Query().Get("updated_at"); got != "2026-07-07 10:00:00" {
			t.Fatalf("updated_at = %q", got)
		}
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Fatalf("page = %q", got)
		}
		if got := r.URL.Query().Get("size"); got != "1000" {
			t.Fatalf("size = %q", got)
		}
		_, _ = w.Write([]byte(`{"code":0,"message":"OK","data":{"date":"2026-07-07","updated_at":"2026-07-07 11:00:00","channel":"RED","total":2,"page":1,"size":1000,"list":[{"post_id":"post-1","source_url":"https://www.xiaohongshu.com/explore/post-1","product_id":["CV758_B4NQ4"],"comment_sentiment":{"positive":0.83}},{"post_id":"post-2","product_id":["CU068_B4MPL&CU068_B4/T2&CU068_B4PUB&CU068_B4CBD"],"comment_sentiment":[]}]}}`))
	}))
	defer server.Close()

	data, err := NewClient(server.URL+"/api", "token", 1000).List(ctx.New(), ListRequest{UpdatedAt: "2026-07-07 10:00:00", Page: 1})
	if err != nil {
		t.Fatal(err)
	}
	if data.Total != 2 || data.Page != 1 || data.Size != 1000 || len(data.List) != 2 || data.List[0].PostID != "post-1" || data.List[1].PostID != "post-2" {
		t.Fatalf("返回数据异常：%+v", data)
	}
	if got := data.List[0].ProductID; len(got) != 1 || got[0] != "CV758_B4NQ4" {
		t.Fatalf("product_id = %#v", got)
	}
	encodedPost, err := json.Marshal(data.List[0])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(encodedPost), `"source_url":"https://www.xiaohongshu.com/explore/post-1"`) {
		t.Fatalf("source_url 未解析：%s", encodedPost)
	}
	if got := data.List[1].ProductID; len(got) != 1 || got[0] != "CU068_B4MPL&CU068_B4/T2&CU068_B4PUB&CU068_B4CBD" {
		t.Fatalf("product_id = %#v", got)
	}
	var sentiment map[string]float64
	if err := json.Unmarshal(data.List[0].CommentSentiment, &sentiment); err != nil || sentiment["positive"] != 0.83 {
		t.Fatalf("comment_sentiment = %s, err = %v", data.List[0].CommentSentiment, err)
	}
	if !json.Valid(data.List[0].CommentSentiment) || !json.Valid(data.List[1].CommentSentiment) {
		t.Fatalf("comment_sentiment raw JSON = %s, %s", data.List[0].CommentSentiment, data.List[1].CommentSentiment)
	}
	var emptySentiment []json.RawMessage
	if err := json.Unmarshal(data.List[1].CommentSentiment, &emptySentiment); err != nil || len(emptySentiment) != 0 {
		t.Fatalf("comment_sentiment = %s, err = %v", data.List[1].CommentSentiment, err)
	}
}

func TestClientList分页大小规则(t *testing.T) {
	tests := []struct {
		name          string
		configured    int
		requested     int
		expectedQuery string
	}{
		{name: "请求未指定时使用构造配置", configured: 200, requested: 0, expectedQuery: "200"},
		{name: "构造配置无效时使用默认值", configured: 0, requested: 0, expectedQuery: "1000"},
		{name: "请求值优先于构造配置", configured: 200, requested: 300, expectedQuery: "300"},
		{name: "请求值超过上限时截断", configured: 200, requested: 1001, expectedQuery: "1000"},
		{name: "构造配置超过上限时截断", configured: 1001, requested: 0, expectedQuery: "1000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := r.URL.Query().Get("size"); got != tt.expectedQuery {
					t.Fatalf("size = %q", got)
				}
				_, _ = w.Write([]byte(`{"code":0,"message":"OK","data":{"total":0,"page":1,"size":1,"list":[]}}`))
			}))
			defer server.Close()

			if _, err := NewClient(server.URL, "token", tt.configured).List(ctx.New(), ListRequest{Page: 1, Size: tt.requested}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestClientListHTTP非200返回错误(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "不可用", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	_, err := NewClient(server.URL, "token", 1000).List(ctx.New(), ListRequest{})
	if err == nil || !strings.Contains(err.Error(), "503") {
		t.Fatalf("error = %v", err)
	}
}

func TestClientList业务码异常返回错误(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		message string
	}{
		{name: "一般业务异常", code: -1, message: "参数错误"},
		{name: "未授权", code: 401, message: "Token无效"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(`{"code":` + strconv.Itoa(tt.code) + `,"message":"` + tt.message + `","data":{}}`))
			}))
			defer server.Close()

			_, err := NewClient(server.URL, "token", 1000).List(ctx.New(), ListRequest{})
			expected := "COACH KOS接口返回失败，code=" + strconv.Itoa(tt.code) + "，message=" + tt.message
			if err == nil || err.Error() != expected {
				t.Fatalf("error = %v", err)
			}
		})
	}
}

func TestClientList非法JSON返回错误(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"code":`))
	}))
	defer server.Close()

	_, err := NewClient(server.URL, "token", 1000).List(ctx.New(), ListRequest{})
	if err == nil || !strings.Contains(err.Error(), "JSON") {
		t.Fatalf("error = %v", err)
	}
}

func TestClientList空配置不发请求(t *testing.T) {
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&requests, 1)
		_, _ = w.Write([]byte(`{"code":0,"message":"OK","data":{}}`))
	}))
	defer server.Close()

	tests := []struct {
		name      string
		baseURL   string
		token     string
		errorText string
	}{
		{name: "baseURL为空", baseURL: "", token: "token", errorText: "COACH KOS配置错误：baseURL不能为空"},
		{name: "token为空", baseURL: server.URL, token: "", errorText: "COACH KOS配置错误：token不能为空"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.baseURL, tt.token, 1000).List(ctx.New(), ListRequest{})
			if err == nil || err.Error() != tt.errorText {
				t.Fatalf("error = %v", err)
			}
		})
	}
	if got := atomic.LoadInt32(&requests); got != 0 {
		t.Fatalf("请求次数 = %d", got)
	}
}
