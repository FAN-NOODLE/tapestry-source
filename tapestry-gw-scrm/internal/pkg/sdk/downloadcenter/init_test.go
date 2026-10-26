package downloadcenter

import (
	"context"
	"fmt"
	"testing"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

func TestDownloadCenter(t *testing.T) {
	Init("https://ws-dev.vchangyi.com/micro-download/v1")
	ctx := ctx.Wrap(context.Background())

	res, err := DC.CreateExport(ctx, CreateExportRequest{
		FileName:      "t21",
		Type:          1,
		StoreId:       11,
		EnterpriseId:  3376928023100864,
		OperationId:   1,
		OperationName: "rabbit",
		SourceType:    1,
		SourceTitleAlias: map[string]string{
			"t1": "标题1",
			"t2": "标题2",
		},
		SourceTitle: []string{"t1", "t2"},

		CreateData: CreateData{
			SourceData: [][]interface{}{
				{"11111", "111111"},
				{"22222", "222222"},
				{"11111", "111111"},
				{"22222", "222222"},
				{"11111", "111111"},
				{"22222", "222222"},
				{"11111", "111111"},
				{"22222", "222222"},
			},
		},
	})

	fmt.Println(res, err)

	// deleteRes, deleteErr := DC.Delete(ctx, DeleteRequest{
	// 	FileId: 3379236470500800,
	// })

	// fmt.Println(deleteRes, deleteErr)

	// listRes, listErr := DC.ExportList(ctx, ExportListRequest{
	// 	EnterpriseId: 3376928023100864,
	// })

	// fmt.Println(listRes, listErr)

	// downloadRes, downloadErr := DC.Download(ctx, DownloadRequest{
	// 	FileId:       3379776417721792,
	// 	DownloadType: 1,
	// })

	// fmt.Println(downloadRes, downloadErr)
}
