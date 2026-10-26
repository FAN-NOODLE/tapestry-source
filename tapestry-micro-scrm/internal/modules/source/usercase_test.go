package source

import (
	"reflect"
	"testing"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type sourceDeleteRepositoryStub struct {
	Repository
	sourceDeleteIDs [][]int64
	delIDs          []int64
}

func (r *sourceDeleteRepositoryStub) SourceDelete(_ ctx.Context, data *RequestSourceDelete) (int64, error) {
	r.sourceDeleteIDs = append(r.sourceDeleteIDs, append([]int64(nil), data.IDs...))
	return int64(len(data.IDs)), nil
}

func (r *sourceDeleteRepositoryStub) Del(_ ctx.Context, id int64) (int64, error) {
	r.delIDs = append(r.delIDs, id)
	return 1, nil
}

func TestSourceDelete单条和批量统一走完整删除逻辑(t *testing.T) {
	repo := &sourceDeleteRepositoryStub{}
	service := Service{repo: repo}

	for _, ids := range [][]int64{{1001}, {2001, 2002}} {
		if _, err := service.SourceDelete(ctx.New(), &RequestSourceDelete{IDs: ids}); err != nil {
			t.Fatal(err)
		}
	}

	wantIDs := [][]int64{{1001}, {2001, 2002}}
	if !reflect.DeepEqual(repo.sourceDeleteIDs, wantIDs) {
		t.Fatalf("完整删除逻辑收到的素材ID=%v，期望=%v", repo.sourceDeleteIDs, wantIDs)
	}
	if len(repo.delIDs) != 0 {
		t.Fatalf("仍调用了仅删除 source 主表的旧逻辑，素材ID=%v", repo.delIDs)
	}
}
