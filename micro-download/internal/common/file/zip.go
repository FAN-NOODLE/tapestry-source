package file

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-download/internal/common"
	"micro-download/internal/modules/download/delivery"
)

type ZipGenerate struct {
}

func (z ZipGenerate) Generate(c ctx.Context, filename string, startLine int, data [][]any, overview *delivery.SourceOverview, complete bool) error {
	if complete {
		folder := fmt.Sprintf("/tmp/%d", overview.ID)
		return common.ZipDir(folder, filename)
	}
	return nil
}
