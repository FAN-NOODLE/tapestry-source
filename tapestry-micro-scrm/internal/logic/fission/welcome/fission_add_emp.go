package welcome

// 裂变活动 - 加BA且成为品牌VIP，需要给客户发送 欢迎语 + 裂变活动link

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/orm/model/micro_scrm/fission"
)

var (
	_ welcomeGetter = &welcomeFissionAndEmp{}
)

type welcomeFissionAndEmp struct {
	ctx  ctx.Context
	data fission.FissionActivityModel
}

func NewWelcomeFissionAndEmp(ctx ctx.Context, data fission.FissionActivityModel) welcomeGetter {
	return &welcomeFissionAndEmp{
		ctx:  ctx,
		data: data,
	}
}

func (w welcomeFissionAndEmp) Get() ([]map[string]interface{}, error) {
	return DecodeWelcomeContent(w.data.DrainageWelcomes.Content)
}
