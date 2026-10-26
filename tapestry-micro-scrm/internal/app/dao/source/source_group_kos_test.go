package source

import (
	"strings"
	"testing"
)

func TestLinkKOSSyncFlagSelect(t *testing.T) {
	wants := []string{
		"EXISTS (",
		"source_link_kos",
		"slk.enterprise_id = source_link.enterprise_id",
		"slk.source_id = s.id",
		"slk.source_link_id = source_link.id",
		"slk.deleted_at = 0",
		"AS is_kos_sync",
	}

	for _, want := range wants {
		if !strings.Contains(linkKOSSyncFlagSelect, want) {
			t.Errorf("linkKOSSyncFlagSelect missing %q", want)
		}
	}
}

func TestLinkKOSSourceOriginFilter(t *testing.T) {
	tests := []struct {
		name         string
		sourceOrigin int
		want         string
		notWant      string
	}{
		{name: "全部", sourceOrigin: 0},
		{name: "手动创建", sourceOrigin: 1, want: "NOT EXISTS ("},
		{name: "CH同步", sourceOrigin: 2, want: "EXISTS (", notWant: "NOT EXISTS ("},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filter := linkKOSSourceOriginFilter(test.sourceOrigin)
			if test.want == "" && filter != "" {
				t.Fatalf("全部来源不应增加过滤条件：%q", filter)
			}
			if test.want != "" && !strings.Contains(filter, test.want) {
				t.Fatalf("过滤条件 %q 未包含 %q", filter, test.want)
			}
			if test.notWant != "" && strings.Contains(filter, test.notWant) {
				t.Fatalf("过滤条件 %q 不应包含 %q", filter, test.notWant)
			}
			if filter != "" && !strings.Contains(filter, "source_link_kos") {
				t.Fatalf("过滤条件未关联 KOS 扩展表：%q", filter)
			}
		})
	}
}
