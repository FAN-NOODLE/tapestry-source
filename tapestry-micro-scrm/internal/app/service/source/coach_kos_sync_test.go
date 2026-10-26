package source

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
	"testing"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/models"
	sourcemodel "micro-scrm/internal/app/models/source"
	coachkos "micro-scrm/internal/pkg/sdk/coach_kos"
)

func TestCoachKOSSync批量写入使用缓冲UID(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "coach_kos_sync.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}

	bufferedCalls := 0
	directCalls := 0
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Name.Name != "upsertCoachKOSPost" {
			continue
		}
		ast.Inspect(function.Body, func(node ast.Node) bool {
			selector, ok := node.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			packageName, ok := selector.X.(*ast.Ident)
			if !ok || packageName.Name != "app" {
				return true
			}
			switch selector.Sel.Name {
			case "UID":
				bufferedCalls++
			case "UIDContext":
				directCalls++
			}
			return true
		})
	}

	if directCalls != 0 || bufferedCalls != 3 {
		t.Fatalf("upsertCoachKOSPost ID 调用：app.UID=%d，app.UIDContext=%d，期望 3 和 0", bufferedCalls, directCalls)
	}
}

type fakeCoachKOSClient struct {
	requests []coachkos.ListRequest
	pages    map[int]coachkos.PageData
	errors   map[int]error
}

func TestBuildCoachKOSGoodsTask保留商品号副本(t *testing.T) {
	productCodes := []string{"CU068_B4MPL", "CU068_B4/T2", "CU068_B4PUB", "CU068_B4CBD"}
	task := buildCoachKOSGoodsTask(1001, 2002, productCodes)

	if task.EnterpriseID != 1001 || task.SourceID != 2002 || !reflect.DeepEqual(task.ProductCodes, productCodes) {
		t.Fatalf("补全任务 = %+v", task)
	}
	productCodes[0] = "changed"
	if task.ProductCodes[0] != "CU068_B4MPL" {
		t.Fatalf("补全任务商品号被原切片修改：%#v", task.ProductCodes)
	}
}

func TestCoachKOS商品关系增量对账计划(t *testing.T) {
	existing := []sourcemodel.RelationGoodsModel{
		{
			MicroBaseModel: models.MicroBaseModel{BaseModel: models.BaseModel{ID: 1}},
			SkuCode:        "CV758 B4NQ4",
			SKU:            "CV758",
		},
		{
			MicroBaseModel: models.MicroBaseModel{BaseModel: models.BaseModel{ID: 2}},
			SkuCode:        "CU068 B4PUB",
			SKU:            "",
		},
		{
			MicroBaseModel: models.MicroBaseModel{BaseModel: models.BaseModel{ID: 3}},
			SkuCode:        "OLD CODE",
			SKU:            "OLD",
		},
	}
	productCodes := []string{"CV758_B4NQ4", "CU068_B4PUB", "NEW_X"}

	plan := planCoachKOSRelationGoods(existing, productCodes)
	if !reflect.DeepEqual(plan.CreateProductCodes, []string{"NEW_X"}) {
		t.Fatalf("待创建商品号=%v", plan.CreateProductCodes)
	}
	if !reflect.DeepEqual(plan.EnrichProductCodes, []string{"CU068_B4PUB", "NEW_X"}) {
		t.Fatalf("待补全商品号=%v，完整商品不应再次查询 PIM", plan.EnrichProductCodes)
	}
	if !reflect.DeepEqual(plan.DeleteIDs, []int64{3}) {
		t.Fatalf("待删除关系ID=%v", plan.DeleteIDs)
	}

	emptyPlan := planCoachKOSRelationGoods(existing, nil)
	if len(emptyPlan.CreateProductCodes) != 0 || len(emptyPlan.EnrichProductCodes) != 0 ||
		!reflect.DeepEqual(emptyPlan.DeleteIDs, []int64{1, 2, 3}) {
		t.Fatalf("product_id 为空时应清理全部旧关系：%+v", emptyPlan)
	}
}

func TestValidateCoachKOSPageRejectsResponsesThatCannotAdvance(t *testing.T) {
	tests := []struct {
		name     string
		received int64
		data     coachkos.PageData
		want     string
	}{
		{name: "page 非法", data: coachkos.PageData{Page: 0, Size: 10}, want: "page=0"},
		{name: "size 非法", data: coachkos.PageData{Page: 1, Size: 0}, want: "size=0"},
		{name: "total 非法", data: coachkos.PageData{Page: 1, Size: 10, Total: -1}, want: "total=-1"},
		{name: "响应页不匹配", data: coachkos.PageData{Page: 2, Size: 10}, want: "与请求 page=1 不一致"},
		{name: "空页但还有数据", data: coachkos.PageData{Page: 1, Size: 10, Total: 11}, want: "短页"},
		{name: "大分页短页会漏数", data: coachkos.PageData{Page: 1, Size: 1000, Total: 10, List: make([]coachkos.PostResult, 5)}, want: "短页"},
		{name: "累计超过 total", received: 3, data: coachkos.PageData{Page: 1, Size: 1, Total: 3, List: []coachkos.PostResult{{}}}, want: "超过 total=3"},
		{name: "列表超过声明 size", data: coachkos.PageData{Page: 1, Size: 1, Total: 3, List: []coachkos.PostResult{{}, {}}}, want: "超过 size=1"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := validateCoachKOSPage(1, test.received, test.data)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("validateCoachKOSPage() error = %v，期望包含 %q", err, test.want)
			}
		})
	}

	next, err := validateCoachKOSPage(1, 0, coachkos.PageData{Page: 1, Size: 2, Total: 3, List: make([]coachkos.PostResult, 2)})
	if err != nil || next != 2 {
		t.Fatalf("正常非末页 next=%d, error=%v，期望累计 2", next, err)
	}
	next, err = validateCoachKOSPage(2, next, coachkos.PageData{Page: 2, Size: 2, Total: 3, List: []coachkos.PostResult{{}}})
	if err != nil || next != 3 {
		t.Fatalf("正常末页短页 next=%d, error=%v，期望累计完成 3", next, err)
	}
}

func TestCoachKOSSync已存在分组赋权失败后重试并缓存成功结果(t *testing.T) {
	var rootGrants, monthGrants int
	service := &CoachKOSSyncService{
		getOrCreateGroup: func(_ app.Ctx, _ int64, parentID int64, _ string) (int64, error) {
			if parentID == 0 {
				return 100, nil
			}
			return 200, nil
		},
		grantGroupAccess: func(_ app.Ctx, _ int64, groupID int64) error {
			if groupID == 100 {
				rootGrants++
				return nil
			}
			monthGrants++
			if monthGrants == 1 {
				return errors.New("grant failed")
			}
			return nil
		},
	}
	cache := make(map[string]int64)

	if _, err := service.ensureCoachKOSGroup(ctx.New(), 1001, "2026-07-01", cache); err == nil {
		t.Fatal("首次月份分组赋权失败时应返回错误")
	}
	if _, err := service.ensureCoachKOSGroup(ctx.New(), 1001, "2026-07-02", cache); err != nil {
		t.Fatalf("第二次应重试成功：%v", err)
	}
	if _, err := service.ensureCoachKOSGroup(ctx.New(), 1001, "2026-07-03", cache); err != nil {
		t.Fatalf("成功缓存后再次确保失败：%v", err)
	}

	if rootGrants != 1 {
		t.Fatalf("一级分组赋权次数 = %d，期望 1", rootGrants)
	}
	if monthGrants != 2 {
		t.Fatalf("月份分组赋权次数 = %d，期望失败后重试且成功后缓存，共 2 次", monthGrants)
	}
}

func TestCoachKOSSync已有素材只更新KOS字段并保留本地字段(t *testing.T) {
	existingSource := sourcemodel.SourceModel{
		MicroBaseModel: models.MicroBaseModel{BaseModel: models.BaseModel{ID: 11, CreatedAt: 101, UpdatedAt: 102, DeletedAt: 99}, EnterpriseID: 1001},
		GroupID:        5, Type: sourcemodel.TypeLink, Status: sourcemodel.StatusEnable,
		AdminID: 77, StartTime: 88, EndTime: 99, MemberID: 66,
		DownloadTimes: 12, ForwardingTimes: 13, CopyTimes: 14, CollectNum: 15, IsTop: 1,
	}
	existingLink := sourcemodel.LinkModel{
		BaseModel:    models.BaseModel{ID: 11, CreatedAt: 201, UpdatedAt: 202, DeletedAt: 99},
		EnterpriseID: 1001, Title: "旧标题", Content: "本地正文",
	}
	existingKOS := sourcemodel.LinkKOSModel{
		MicroBaseModel: models.MicroBaseModel{BaseModel: models.BaseModel{ID: 33, CreatedAt: 301}, EnterpriseID: 1001},
		SourceID:       11, SourceLinkID: 11, PostID: "post-1",
	}
	generatedSource, generatedLink, generatedKOS := buildCoachKOSModels(1001, 11, 9, 999, coachkos.PostResult{
		PostID: "post-1", DataDate: "2026-06-03", Name: "新标题", Content: "新描述", URL: "https://new", CoverURL: "https://cover", InfluencerName: "新作者",
	}, "https://cdn.example.com/cover.jpg?imageMogr2/quality/80/format/jpg")
	wantCreatedAt := time.Date(2026, 6, 3, 0, 0, 0, 0, time.FixedZone("Asia/Shanghai", 8*60*60)).Unix()

	sourceModel, linkModel, kosModel, err := mergeExistingCoachKOSModels("post-1", existingKOS, existingSource, existingLink, generatedSource, generatedLink, generatedKOS)
	if err != nil {
		t.Fatal(err)
	}
	if sourceModel.GroupID != 9 || sourceModel.CreatedAt != wantCreatedAt || sourceModel.UpdatedAt != 999 || sourceModel.DeletedAt != 0 {
		t.Fatalf("KOS 可更新字段异常：%+v", sourceModel)
	}
	if sourceModel.AdminID != 0 || sourceModel.StartTime != 88 || sourceModel.MemberID != 66 || sourceModel.DownloadTimes != 12 || sourceModel.ForwardingTimes != 13 || sourceModel.CopyTimes != 14 || sourceModel.CollectNum != 15 || sourceModel.IsTop != 1 {
		t.Fatalf("本地 source 字段被覆盖：%+v", sourceModel)
	}
	if linkModel.Title != "新标题" || linkModel.LinkURL != "https://new" || linkModel.Content != "" || linkModel.CreatedAt != wantCreatedAt {
		t.Fatalf("link 合并结果异常：%+v", linkModel)
	}
	if kosModel.ID != 33 || kosModel.CreatedAt != 301 || kosModel.SourceID != 11 || kosModel.SourceLinkID != 11 {
		t.Fatalf("KOS 扩展 ID/created_at 未复用：%+v", kosModel)
	}

	generatedSource.MemberID = 88
	generatedSource.AdminID = 99
	sourceModel, _, _, err = mergeExistingCoachKOSModels("post-1", existingKOS, existingSource, existingLink, generatedSource, generatedLink, generatedKOS)
	if err != nil || sourceModel.MemberID != 88 || sourceModel.AdminID != 99 {
		t.Fatalf("已有素材上传人未更新：source=%+v err=%v", sourceModel, err)
	}
}

func TestCoachKOSSync按页传递上传人(t *testing.T) {
	client := &fakeCoachKOSClient{pages: map[int]coachkos.PageData{
		1: {Total: 2, Page: 1, Size: 2, List: []coachkos.PostResult{
			{PostID: "post-1", Name: "素材1", DataDate: "2026-07-01", EmployeeID: " member-1 "},
			{PostID: "post-2", Name: "素材2", DataDate: "2026-07-01", EmployeeID: "missing"},
		}},
	}}
	saved := make(map[string]coachKOSUploader)
	service := &CoachKOSSyncService{
		client:   client,
		pageSize: 2,
		resolveUploaders: func(_ app.Ctx, enterpriseID int64, posts []coachkos.PostResult) (map[string]coachKOSUploader, error) {
			if enterpriseID != 1001 || len(posts) != 2 {
				t.Fatalf("按页查询上传人参数异常：enterpriseID=%d posts=%+v", enterpriseID, posts)
			}
			return map[string]coachKOSUploader{"member-1": {memberID: 12, adminID: 34}}, nil
		},
		ensureGroup: func(_ app.Ctx, _ int64, _ string, _ map[string]int64) (int64, error) { return 9, nil },
		upsertPost: func(_ app.Ctx, _, _ int64, post coachkos.PostResult, _ string, uploader coachKOSUploader) (CoachKOSGoodsTask, error) {
			saved[post.PostID] = uploader
			return CoachKOSGoodsTask{}, nil
		},
	}
	if err := service.SyncEnterprise(ctx.New(), 1001, time.Now()); err != nil {
		t.Fatal(err)
	}
	if saved["post-1"] != (coachKOSUploader{memberID: 12, adminID: 34}) || saved["post-2"] != (coachKOSUploader{}) {
		t.Fatalf("传给素材保存逻辑的上传人异常：%+v", saved)
	}
}

func TestCoachKOSSync异常KOS关联在写入前返回错误(t *testing.T) {
	tests := []struct {
		name           string
		kos            sourcemodel.LinkKOSModel
		source         sourcemodel.SourceModel
		link           sourcemodel.LinkModel
		expectedErrors []string
	}{
		{
			name:           "关联ID为零",
			kos:            sourcemodel.LinkKOSModel{SourceID: 0, SourceLinkID: 0},
			expectedErrors: []string{"post_id=post-1", "source_id=0", "source_link_id=0"},
		},
		{
			name:           "source缺失",
			kos:            sourcemodel.LinkKOSModel{SourceID: 11, SourceLinkID: 12},
			link:           sourcemodel.LinkModel{BaseModel: models.BaseModel{ID: 12}},
			expectedErrors: []string{"post_id=post-1", "source_id=11", "source_link_id=12"},
		},
		{
			name:           "link缺失",
			kos:            sourcemodel.LinkKOSModel{SourceID: 11, SourceLinkID: 12},
			source:         sourcemodel.SourceModel{MicroBaseModel: models.MicroBaseModel{BaseModel: models.BaseModel{ID: 11}}},
			expectedErrors: []string{"post_id=post-1", "source_id=11", "source_link_id=12"},
		},
		{
			name:           "source和link主键不一致",
			kos:            sourcemodel.LinkKOSModel{SourceID: 11, SourceLinkID: 12},
			source:         sourcemodel.SourceModel{MicroBaseModel: models.MicroBaseModel{BaseModel: models.BaseModel{ID: 11}}},
			link:           sourcemodel.LinkModel{BaseModel: models.BaseModel{ID: 12}},
			expectedErrors: []string{"post_id=post-1", "source_id=11", "source_link_id=12"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, _, _, err := mergeExistingCoachKOSModels("post-1", test.kos, test.source, test.link, sourcemodel.SourceModel{}, sourcemodel.LinkModel{}, sourcemodel.LinkKOSModel{})
			if err == nil {
				t.Fatal("异常关联应在保存前返回错误")
			}
			for _, expected := range test.expectedErrors {
				if !strings.Contains(err.Error(), expected) {
					t.Fatalf("错误 %q 未包含 %q", err, expected)
				}
			}
		})
	}
}

func (f *fakeCoachKOSClient) List(_ ctx.Context, request coachkos.ListRequest) (coachkos.PageData, error) {
	f.requests = append(f.requests, request)
	if err := f.errors[request.Page]; err != nil {
		return coachkos.PageData{}, err
	}
	return f.pages[request.Page], nil
}

func TestCoachKOSSync每页落库后立即补全再拉下一页(t *testing.T) {
	client := &fakeCoachKOSClient{pages: map[int]coachkos.PageData{
		1: {
			Total: 3,
			Page:  1,
			Size:  2,
			List: []coachkos.PostResult{
				{PostID: "post-1", Name: "素材1", DataDate: "2026-07-01", ProductID: []string{"sku-1"}},
				{PostID: "post-2", Name: "素材2", DataDate: "2026-07-02"},
			},
		},
		2: {
			Total: 3,
			Page:  2,
			Size:  2,
			List: []coachkos.PostResult{
				{PostID: "post-3", Name: "素材3", DataDate: "2026-07-03", ProductID: []string{"sku-3"}},
			},
		},
	}}

	var events []string
	service := &CoachKOSSyncService{
		client:   client,
		pageSize: 2,
		ensureGroup: func(_ app.Ctx, _ int64, _ string, _ map[string]int64) (int64, error) {
			return 9001, nil
		},
		upsertPost: func(_ app.Ctx, enterpriseID, _ int64, post coachkos.PostResult, _ string, _ coachKOSUploader) (CoachKOSGoodsTask, error) {
			events = append(events, "save:"+post.PostID)
			return CoachKOSGoodsTask{EnterpriseID: enterpriseID, SourceID: int64(len(events)), ProductCodes: post.ProductID}, nil
		},
		enrichGoods: func(_ app.Ctx, tasks []CoachKOSGoodsTask) {
			if len(tasks) != 1 || len(tasks[0].ProductCodes) != 1 {
				t.Fatalf("本页补全任务=%+v，期望仅包含有待补商品的素材", tasks)
			}
			events = append(events, fmt.Sprintf("enrich:requests=%d:product=%s", len(client.requests), tasks[0].ProductCodes[0]))
		},
	}

	err := service.SyncEnterprise(ctx.New(), 1001, time.Date(2026, 7, 7, 10, 30, 0, 0, time.Local))
	if err != nil {
		t.Fatalf("SyncEnterprise() error = %v", err)
	}
	if got, want := strings.Join(events, ","), "save:post-1,save:post-2,enrich:requests=1:product=sku-1,save:post-3,enrich:requests=2:product=sku-3"; got != want {
		t.Fatalf("处理顺序=%q，期望=%q", got, want)
	}
}

func TestCoachKOSSync分页并聚合单条错误(t *testing.T) {
	client := &fakeCoachKOSClient{pages: map[int]coachkos.PageData{
		1: {
			Total: 3,
			Page:  1,
			Size:  2,
			List: []coachkos.PostResult{
				{PostID: "post-1", Name: "素材1", DataDate: "2026-07-01", ProductID: []string{"sku-1"}},
				{PostID: "post-2", Name: "素材2", DataDate: "2026-07-02", ProductID: []string{"sku-2"}},
			},
		},
		2: {
			Total: 3,
			Page:  2,
			Size:  2,
			List: []coachkos.PostResult{
				{PostID: "post-3", Name: "素材3", DataDate: "2026-07-03", ProductID: []string{"sku-3"}},
			},
		},
	}}

	var saved []string
	var enrichedTasks []CoachKOSGoodsTask
	service := &CoachKOSSyncService{
		client:   client,
		pageSize: 2,
		ensureGroup: func(_ app.Ctx, _ int64, _ string, _ map[string]int64) (int64, error) {
			return 9001, nil
		},
		upsertPost: func(_ app.Ctx, enterpriseID, groupID int64, post coachkos.PostResult, _ string, _ coachKOSUploader) (CoachKOSGoodsTask, error) {
			saved = append(saved, post.PostID)
			if post.PostID == "post-2" {
				return CoachKOSGoodsTask{}, errors.New("save failed")
			}
			return CoachKOSGoodsTask{EnterpriseID: enterpriseID, SourceID: int64(len(saved)), ProductCodes: post.ProductID}, nil
		},
		enrichGoods: func(_ app.Ctx, tasks []CoachKOSGoodsTask) {
			enrichedTasks = append(enrichedTasks, tasks...)
		},
	}

	cursor := time.Date(2026, 7, 7, 10, 30, 0, 0, time.FixedZone("Asia/Shanghai", 8*60*60))
	err := service.SyncEnterprise(ctx.New(), 1001, cursor)

	if len(client.requests) != 2 {
		t.Fatalf("请求页数 = %d，期望 2", len(client.requests))
	}
	for index, request := range client.requests {
		if request.Page != index+1 || request.Size != 2 {
			t.Fatalf("第 %d 次请求 = %+v", index+1, request)
		}
		if request.UpdatedAt != "2026-07-07 10:30:00" {
			t.Fatalf("第 %d 次请求游标 = %q", index+1, request.UpdatedAt)
		}
	}
	if got := strings.Join(saved, ","); got != "post-1,post-2,post-3" {
		t.Fatalf("保存顺序 = %q", got)
	}
	if len(enrichedTasks) != 2 || enrichedTasks[0].ProductCodes[0] != "sku-1" || enrichedTasks[1].ProductCodes[0] != "sku-3" {
		t.Fatalf("成功任务 = %+v", enrichedTasks)
	}
	if err == nil || !strings.Contains(err.Error(), "enterprise_id=1001") || !strings.Contains(err.Error(), "post_id=post-2") || !strings.Contains(err.Error(), "save failed") {
		t.Fatalf("聚合错误 = %v", err)
	}
}

func TestCoachKOSSync首次失效素材在创建分组前跳过(t *testing.T) {
	client := &fakeCoachKOSClient{pages: map[int]coachkos.PageData{
		1: {
			Total: 1,
			Page:  1,
			Size:  1,
			List: []coachkos.PostResult{{
				PostID:   "invalid-post",
				Name:     "失效素材",
				DataDate: "2026-07-01",
				IsValid:  "N",
			}},
		},
	}}
	var queried bool
	service := &CoachKOSSyncService{
		client:   client,
		pageSize: 1,
		getKOSByPostID: func(_ app.Ctx, enterpriseID int64, postID string) (sourcemodel.LinkKOSModel, error) {
			queried = true
			if enterpriseID != 1001 || postID != "invalid-post" {
				t.Fatalf("查询参数 enterpriseID=%d postID=%s", enterpriseID, postID)
			}
			return sourcemodel.LinkKOSModel{}, nil
		},
		ensureGroup: func(_ app.Ctx, _ int64, _ string, _ map[string]int64) (int64, error) {
			t.Fatal("首次失效素材不应创建分组")
			return 0, nil
		},
		upsertPost: func(_ app.Ctx, _, _ int64, _ coachkos.PostResult, _ string, _ coachKOSUploader) (CoachKOSGoodsTask, error) {
			t.Fatal("首次失效素材不应写入素材表")
			return CoachKOSGoodsTask{}, nil
		},
	}

	err := service.SyncEnterprise(ctx.New(), 1001, time.Now())
	if err != nil {
		t.Fatalf("跳过首次失效素材不应报错：%v", err)
	}
	if !queried {
		t.Fatal("失效素材应先查询本地 KOS 记录")
	}
}

func TestCoachKOSSync空标题素材在任何落库准备前跳过(t *testing.T) {
	client := &fakeCoachKOSClient{pages: map[int]coachkos.PageData{
		1: {
			Total: 2,
			Page:  1,
			Size:  2,
			List: []coachkos.PostResult{
				{PostID: "empty-title", Name: "", DataDate: "2026-07-01", ProductID: []string{"sku-1"}},
				{PostID: "blank-title", Name: "  \t", DataDate: "2026-07-02", ProductID: []string{"sku-2"}},
			},
		},
	}}
	var (
		groupCalls  int
		upsertCalls int
		enrichCalls int
		logs        []string
	)
	service := &CoachKOSSyncService{
		client:   client,
		pageSize: 2,
		ensureGroup: func(_ app.Ctx, _ int64, _ string, _ map[string]int64) (int64, error) {
			groupCalls++
			return 9001, nil
		},
		upsertPost: func(_ app.Ctx, _, _ int64, _ coachkos.PostResult, _ string, _ coachKOSUploader) (CoachKOSGoodsTask, error) {
			upsertCalls++
			return CoachKOSGoodsTask{}, nil
		},
		enrichGoods: func(_ app.Ctx, _ []CoachKOSGoodsTask) {
			enrichCalls++
		},
		logInfo: func(_ app.Ctx, format string, args ...any) {
			logs = append(logs, fmt.Sprintf(format, args...))
		},
	}

	if err := service.SyncEnterprise(ctx.New(), 1001, time.Now()); err != nil {
		t.Fatalf("跳过空标题素材不应报错：%v", err)
	}
	if groupCalls != 0 || upsertCalls != 0 || enrichCalls != 0 {
		t.Fatalf("空标题素材触发了后续处理：group=%d upsert=%d enrich=%d", groupCalls, upsertCalls, enrichCalls)
	}
	if len(logs) == 0 || !strings.Contains(logs[len(logs)-1], "success=0") ||
		!strings.Contains(logs[len(logs)-1], "skipped_empty_title=2") {
		t.Fatalf("结束日志未记录空标题跳过数量：%v", logs)
	}
}

func TestCoachKOSSync空列表输出进度日志(t *testing.T) {
	client := &fakeCoachKOSClient{pages: map[int]coachkos.PageData{
		1: {Total: 0, Page: 1, Size: 10, List: []coachkos.PostResult{}},
	}}
	var logs []string
	service := &CoachKOSSyncService{
		client:   client,
		pageSize: 10,
		logInfo: func(_ app.Ctx, format string, args ...any) {
			logs = append(logs, fmt.Sprintf(format, args...))
		},
	}
	cursor := time.Date(2026, 7, 13, 10, 56, 57, 0, time.Local)

	err := service.SyncEnterprise(ctx.New(), 1001, cursor)
	if err != nil {
		t.Fatalf("空列表同步不应报错：%v", err)
	}
	if len(logs) != 3 {
		t.Fatalf("日志数量=%d，期望开始、分页、结束共3条：%v", len(logs), logs)
	}
	for _, want := range []string{"KOS素材同步开始", "enterprise_id=1001", "updated_at=2026-07-13 10:56:57", "page_size=10"} {
		if !strings.Contains(logs[0], want) {
			t.Fatalf("开始日志 %q 缺少 %q", logs[0], want)
		}
	}
	for _, want := range []string{"KOS素材分页返回", "page=1", "total=0", "list_count=0"} {
		if !strings.Contains(logs[1], want) {
			t.Fatalf("分页日志 %q 缺少 %q", logs[1], want)
		}
	}
	for _, want := range []string{"KOS素材同步结束", "received=0", "success=0", "skipped_invalid=0", "failed=0", "goods_tasks=0"} {
		if !strings.Contains(logs[2], want) {
			t.Fatalf("结束日志 %q 缺少 %q", logs[2], want)
		}
	}
}

func TestCoachKOSSync页错误立即结束并保留记录错误(t *testing.T) {
	client := &fakeCoachKOSClient{
		pages: map[int]coachkos.PageData{
			1: {Total: 3, Page: 1, Size: 2, List: []coachkos.PostResult{
				{PostID: "post-1", Name: "素材1", DataDate: "2026-07-01"},
				{PostID: "post-2", Name: "素材2", DataDate: "2026-07-01"},
			}},
		},
		errors: map[int]error{2: errors.New("api failed")},
	}
	var enrichedTasks []CoachKOSGoodsTask
	service := &CoachKOSSyncService{
		client:   client,
		pageSize: 2,
		ensureGroup: func(_ app.Ctx, _ int64, _ string, _ map[string]int64) (int64, error) {
			return 9001, nil
		},
		upsertPost: func(_ app.Ctx, _, _ int64, post coachkos.PostResult, _ string, _ coachKOSUploader) (CoachKOSGoodsTask, error) {
			if post.PostID == "post-1" {
				return CoachKOSGoodsTask{}, errors.New("record failed")
			}
			return CoachKOSGoodsTask{SourceID: 2, ProductCodes: []string{"sku-2"}}, nil
		},
		enrichGoods: func(_ app.Ctx, tasks []CoachKOSGoodsTask) {
			enrichedTasks = append(enrichedTasks, tasks...)
		},
	}

	err := service.SyncEnterprise(ctx.New(), 1001, time.Date(2026, 7, 14, 10, 30, 0, 0, time.Local))

	if len(client.requests) != 2 {
		t.Fatalf("请求页数 = %d，期望 API 失败后停止在第 2 页", len(client.requests))
	}
	if len(enrichedTasks) != 1 || enrichedTasks[0].SourceID != 2 {
		t.Fatalf("成功任务 = %+v", enrichedTasks)
	}
	if err == nil || !strings.Contains(err.Error(), "post_id=post-1") || !strings.Contains(err.Error(), "page=2") || !strings.Contains(err.Error(), "api failed") {
		t.Fatalf("聚合错误 = %v", err)
	}
}
