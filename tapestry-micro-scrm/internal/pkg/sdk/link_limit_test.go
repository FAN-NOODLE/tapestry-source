package sdk_test

import (
	"strings"
	"testing"
	"unicode/utf8"

	"micro-scrm/internal/pkg/sdk/work"
)

func TestLink序列化时按企微字节限制裁剪(t *testing.T) {
	original := work.Link{
		Title: strings.Repeat("中", 43) + "A",
		Desc:  strings.Repeat("文", 171) + "B",
		Url:   "https://example.com",
	}

	got := work.LimitLinkForQy(original)

	if got.Title != strings.Repeat("中", 42) {
		t.Fatalf("标题裁剪结果=%q，字节数=%d", got.Title, len(got.Title))
	}
	if got.Desc != strings.Repeat("文", 170) {
		t.Fatalf("描述裁剪结果字节数=%d", len(got.Desc))
	}
	if len(got.Title) >= 128 || len(got.Desc) >= 512 {
		t.Fatalf("仍超过企微限制：title=%d desc=%d", len(got.Title), len(got.Desc))
	}
	if !utf8.ValidString(got.Title) || !utf8.ValidString(got.Desc) {
		t.Fatal("裁剪后产生了无效 UTF-8 字符串")
	}
	if got.Url != original.Url {
		t.Fatalf("非裁剪字段被修改：url=%q", got.Url)
	}
}

func TestLink序列化时保留限制内内容(t *testing.T) {
	original := work.Link{
		Title: strings.Repeat("a", 127),
		Desc:  strings.Repeat("b", 511),
	}

	got := work.LimitLinkForQy(original)
	if got.Title != original.Title || got.Desc != original.Desc {
		t.Fatalf("限制内内容不应变化：title=%d desc=%d", len(got.Title), len(got.Desc))
	}
}

func TestMomentLink序列化时按企微字节限制裁剪标题(t *testing.T) {
	original := work.MomentLink{
		Title: strings.Repeat("😀", 32),
		Url:   "https://example.com",
	}

	got := work.LimitMomentLinkForQy(original)
	if got.Title != strings.Repeat("😀", 31) || len(got.Title) >= 128 || !utf8.ValidString(got.Title) {
		t.Fatalf("朋友圈标题裁剪异常：title=%q bytes=%d", got.Title, len(got.Title))
	}
	if got.Url != original.Url {
		t.Fatalf("非裁剪字段被修改：url=%q", got.Url)
	}
}
