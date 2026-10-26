package source

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	sourcemodel "micro-scrm/internal/app/models/source"
	coachkos "micro-scrm/internal/pkg/sdk/coach_kos"
)

func TestCoachKOSMapping月份分组名称(t *testing.T) {
	name, err := coachKOSMonthGroupName("2026-06-03")
	if err != nil {
		t.Fatal(err)
	}
	if name != "6月小红书素材" {
		t.Fatalf("分组名称 = %q", name)
	}

	if _, err := coachKOSMonthGroupName("2026-06-31"); err == nil {
		t.Fatal("非法日期未返回错误")
	}
}

func TestCoachKOS素材禁用条件(t *testing.T) {
	tests := []struct {
		name     string
		post     coachkos.PostResult
		disabled bool
	}{
		{
			name:     "is_valid为N",
			post:     coachkos.PostResult{IsValid: "N", CommentSentiment: json.RawMessage(`{"负面":0}`), PriceMentioned: "N"},
			disabled: true,
		},
		{
			name:     "真实数据负面超过阈值",
			post:     coachkos.PostResult{IsValid: "Y", CommentSentiment: json.RawMessage(`{"\u4e2d\u6027":0.5,"\u8d1f\u9762":0.5}`), PriceMentioned: "N"},
			disabled: true,
		},
		{
			name:     "负面等于阈值",
			post:     coachkos.PostResult{IsValid: "Y", CommentSentiment: json.RawMessage(`{"负面":0.25}`), PriceMentioned: "N"},
			disabled: false,
		},
		{
			name:     "价格被提及不再禁用",
			post:     coachkos.PostResult{IsValid: "Y", CommentSentiment: json.RawMessage(`[]`), PriceMentioned: "Y"},
			disabled: false,
		},
		{
			name:     "空数组",
			post:     coachkos.PostResult{IsValid: "Y", CommentSentiment: json.RawMessage(`[]`), PriceMentioned: "N"},
			disabled: false,
		},
		{
			name:     "空对象",
			post:     coachkos.PostResult{IsValid: "Y", CommentSentiment: json.RawMessage(`{}`), PriceMentioned: "N"},
			disabled: false,
		},
		{
			name:     "异常JSON",
			post:     coachkos.PostResult{IsValid: "Y", CommentSentiment: json.RawMessage(`invalid`), PriceMentioned: "N"},
			disabled: false,
		},
		{
			name:     "异常JSON且价格被提及不再禁用",
			post:     coachkos.PostResult{IsValid: "Y", CommentSentiment: json.RawMessage(`invalid`), PriceMentioned: "Y"},
			disabled: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sourceModel, _, _ := buildCoachKOSModels(1001, 2002, 3003, time.Now().Unix(), test.post, test.post.CoverURL)
			wantStatus := int32(sourcemodel.StatusEnable)
			if test.disabled {
				wantStatus = sourcemodel.StatusDisable
			}
			if sourceModel.Status != wantStatus {
				t.Fatalf("素材状态=%d，期望=%d，post=%+v", sourceModel.Status, wantStatus, test.post)
			}
		})
	}
}

func TestCoachKOSMapping描述按字符截断(t *testing.T) {
	content255 := strings.Repeat("文", 255)
	if got := truncateCoachKOSDescription(content255); got != content255 {
		t.Fatalf("255 字符不应截断，字符数 = %d", len([]rune(got)))
	}

	content256 := strings.Repeat("文", 256)
	got := truncateCoachKOSDescription(content256)
	if len([]rune(got)) != 255 {
		t.Fatalf("截断后字符数 = %d", len([]rune(got)))
	}
	if got != strings.Repeat("文", 254)+"…" {
		t.Fatalf("截断结果不符合预期：%q", got)
	}
}

func TestCoachKOSMapping正文为空时使用话题作为描述(t *testing.T) {
	hashtag := strings.Repeat("话", 256)
	post := coachkos.PostResult{
		Content: " \n\t",
		Hashtag: hashtag,
	}

	_, linkModel, kosModel := buildCoachKOSModels(1001, 2002, 3003, 1783987200, post, "")
	wantDescription := strings.Repeat("话", 254) + "…"
	if linkModel.Description != wantDescription {
		t.Fatalf("链接描述=%q，期望=%q", linkModel.Description, wantDescription)
	}
	if kosModel.Content != post.Content || kosModel.Hashtag != hashtag {
		t.Fatalf("KOS 原始正文或话题被修改：content=%q hashtag=%q", kosModel.Content, kosModel.Hashtag)
	}
}

func TestCoachKOSMapping模型字段转换(t *testing.T) {
	now := int64(1783987200)
	post := coachkos.PostResult{
		PostID:           "post-1",
		SourceURL:        "https://www.xiaohongshu.com/explore/post-1",
		DataDate:         "2026-06-03",
		PublishedAt:      "2026-06-03 10:11:12",
		Name:             "帖子标题",
		Content:          strings.Repeat("内容", 128),
		Hashtag:          "#新品",
		URL:              "https://example.com/post-1",
		CoverURL:         "https://example.com/cover.jpg",
		InfluencerID:     "influencer-1",
		InfluencerName:   "博主",
		EmployeeID:       "employee-1",
		EmployeeTag:      "员工标签",
		StoreID:          "store-1",
		StoreName:        "门店",
		ProductID:        []string{" CU068_B4MPL&CU068_B4/T2 ", "CU068_B4MPL&&CU068_B4PUB", " CU068_B4CBD "},
		ContentTag:       "内容标签",
		LikeCount:        11,
		CommentCount:     12,
		FavoriteCount:    13,
		ShareCount:       14,
		InteractionCount: 15,
		IsValid:          "N",
		CommentSentiment: json.RawMessage(`{"positive":0.83}`),
		PriceMentioned:   "价格信息",
		UpdatedAt:        "2026-06-04 10:11:12",
		SourceUpdatedAt:  "2026-06-05 10:11:12",
	}

	materialCoverURL := "https://cdn.example.com/uploads/cover.jpg?imageMogr2/quality/80/format/jpg"
	sourceModel, linkModel, kosModel := buildCoachKOSModels(1001, 2002, 3003, now, post, materialCoverURL)
	shanghai := time.FixedZone("Asia/Shanghai", 8*60*60)
	wantCreatedAt := time.Date(2026, 6, 3, 0, 0, 0, 0, shanghai).Unix()

	if sourceModel.ID != 2002 || sourceModel.EnterpriseID != 1001 || sourceModel.GroupID != 3003 {
		t.Fatalf("素材标识字段异常：%+v", sourceModel)
	}
	if sourceModel.Type != sourcemodel.TypeLink || sourceModel.Status != sourcemodel.StatusDisable {
		t.Fatalf("素材类型或状态异常：%+v", sourceModel)
	}
	if sourceModel.CreatedAt != wantCreatedAt || sourceModel.UpdatedAt != now || int64(sourceModel.StartTime) != now {
		t.Fatalf("素材创建时间未使用 data_date：%+v", sourceModel)
	}
	wantEnd := time.Date(2035, 12, 31, 23, 59, 59, 0, shanghai).Unix()
	if int64(sourceModel.EndTime) != wantEnd {
		t.Fatalf("结束时间 = %d，期望 %d", sourceModel.EndTime, wantEnd)
	}

	if linkModel.ID != 2002 || linkModel.CreatedAt != wantCreatedAt || linkModel.UpdatedAt != now {
		t.Fatalf("链接基础字段异常：%+v", linkModel)
	}
	if linkModel.Title != post.Name || linkModel.Description != strings.Repeat("内容", 127)+"…" ||
		linkModel.CoverURL != materialCoverURL || linkModel.LinkType != 1 || linkModel.LinkURL != post.URL ||
		linkModel.Author != post.InfluencerName || linkModel.Content != "" {
		t.Fatalf("链接字段异常：%+v", linkModel)
	}

	if kosModel.ID != 0 || kosModel.EnterpriseID != 1001 || kosModel.SourceID != 2002 || kosModel.SourceLinkID != 2002 {
		t.Fatalf("扩展表标识字段异常：%+v", kosModel)
	}
	if kosModel.CreatedAt != now || kosModel.UpdatedAt != now {
		t.Fatalf("扩展表基础时间未使用本地 now：%+v", kosModel)
	}
	wantProductIDs := `[" CU068_B4MPL&CU068_B4/T2 ","CU068_B4MPL&&CU068_B4PUB"," CU068_B4CBD "]`
	if !strings.Contains(kosModel.ProductID, "&") {
		t.Fatalf("原始商品数组未保留字面 &：%q", kosModel.ProductID)
	}
	if kosModel.PostID != post.PostID || kosModel.SourceURL != post.SourceURL || kosModel.DataDate != post.DataDate || kosModel.PublishedAt != post.PublishedAt ||
		kosModel.Name != post.Name || kosModel.Content != post.Content || kosModel.Hashtag != post.Hashtag ||
		kosModel.URL != post.URL || kosModel.CoverURL != post.CoverURL || kosModel.InfluencerID != post.InfluencerID ||
		kosModel.InfluencerName != post.InfluencerName || kosModel.EmployeeID != post.EmployeeID ||
		kosModel.EmployeeTag != post.EmployeeTag || kosModel.StoreID != post.StoreID || kosModel.StoreName != post.StoreName ||
		kosModel.ProductID != wantProductIDs || kosModel.ContentTag != post.ContentTag ||
		kosModel.LikeCount != post.LikeCount || kosModel.CommentCount != post.CommentCount ||
		kosModel.FavoriteCount != post.FavoriteCount || kosModel.ShareCount != post.ShareCount ||
		kosModel.InteractionCount != post.InteractionCount || kosModel.IsValid != post.IsValid ||
		kosModel.CommentSentiment != string(post.CommentSentiment) || kosModel.PriceMentioned != post.PriceMentioned ||
		kosModel.APIUpdatedAt != post.UpdatedAt || kosModel.SourceUpdatedAt != post.SourceUpdatedAt {
		t.Fatalf("扩展字段映射异常：%+v", kosModel)
	}

	post.IsValid = "Y"
	post.ProductID = nil
	sourceModel, _, kosModel = buildCoachKOSModels(1001, 2002, 3003, now, post, post.CoverURL)
	if sourceModel.Status != sourcemodel.StatusEnable {
		t.Fatalf("有效素材状态 = %d", sourceModel.Status)
	}
	if kosModel.ProductID != "[]" {
		t.Fatalf("空商品数组序列化结果 = %q", kosModel.ProductID)
	}
}

func TestCoachKOSMapping基础商品关系构造(t *testing.T) {
	productCodes := expandCoachKOSProductCodes([]string{
		" CU068_B4MPL&CU068_B4/T2 ",
		"CU068_B4MPL&&CU068_B4PUB",
		" CU068_B4CBD ",
	})
	goods := buildCoachKOSRelationGoods(1001, 2002, 1783987200, productCodes)
	if len(goods) != 4 {
		t.Fatalf("商品关系数量 = %d", len(goods))
	}
	for i, productNo := range []string{"CU068_B4MPL", "CU068_B4/T2", "CU068_B4PUB", "CU068_B4CBD"} {
		got := goods[i]
		if got.EnterpriseID != 1001 || got.SourceID != 2002 || got.Type != 3 || got.SKU != "" ||
			got.SkuCode != strings.ReplaceAll(productNo, "_", " ") || got.AddType != 2 {
			t.Fatalf("第 %d 条商品关系异常：%+v", i, got)
		}
		if got.CreatedAt != 1783987200 || got.UpdatedAt != 1783987200 {
			t.Fatalf("第 %d 条商品关系时间异常：%+v", i, got)
		}
	}
}

func TestExpandCoachKOSProductCodes拆分去空并稳定去重(t *testing.T) {
	got := expandCoachKOSProductCodes([]string{
		" CU068_B4MPL&CU068_B4/T2 ",
		"CU068_B4MPL&&CU068_B4PUB",
		" CU068_B4CBD ",
	})
	want := []string{"CU068_B4MPL", "CU068_B4/T2", "CU068_B4PUB", "CU068_B4CBD"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("商品号 = %#v，期望 %#v", got, want)
	}
}
