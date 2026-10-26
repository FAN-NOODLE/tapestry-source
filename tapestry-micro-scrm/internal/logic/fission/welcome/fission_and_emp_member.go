package welcome

// 裂变活动 - 加BA且成为品牌VIP，需要给客户发送 欢迎语 + 小程序

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/orm/model/micro_scrm/fission"
)

var (
	_ welcomeGetter = &welcomeFissionAndEmpMember{}
)

type welcomeFissionAndEmpMember struct {
	ctx  ctx.Context
	data fission.FissionActivityModel
}

func NewWelcomeFissionAndEmpMember(ctx ctx.Context, data fission.FissionActivityModel) welcomeGetter {
	return &welcomeFissionAndEmpMember{
		ctx:  ctx,
		data: data,
	}
}

func (w welcomeFissionAndEmpMember) Get() ([]map[string]interface{}, error) {
	return DecodeWelcomeContent(w.data.DrainageWelcomes.Content)
}
