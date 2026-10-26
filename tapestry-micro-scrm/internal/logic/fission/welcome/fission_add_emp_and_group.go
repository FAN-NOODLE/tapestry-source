package welcome

// 裂变活动 - 加BA且成为品牌VIP，需要给客户发送 欢迎语 + 群二维码

import (
	"math/rand"
	"sync"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/orm/model/micro_scrm/fission"
)

var (
	_ welcomeGetter = &welcomeFissionAndEmpGroup{}
)

type welcomeFissionAndEmpGroup struct {
	ctx  ctx.Context
	data fission.FissionActivityModel
	sync.Mutex
}

func NewWelcomeFissionAndEmpGroup(ctx ctx.Context, data fission.FissionActivityModel) welcomeGetter {
	return &welcomeFissionAndEmpGroup{
		ctx:  ctx,
		data: data,
	}
}

func (w welcomeFissionAndEmpGroup) Get() ([]map[string]interface{}, error) {
	msg, err := DecodeWelcomeContent(w.data.DrainageWelcomes.Content)
	if err != nil {
		return nil, err
	}

	// 随机发二维码
	num := len(w.data.FissionChatGroupList)
	index := 0
	if num > 1 {
		index = rand.Intn(num - 1)
	}
	imageMap := make(map[string]interface{})
	valMap := make(map[string]interface{})
	valMap["pic_url"] = w.data.FissionChatGroupList[index].QrCode
	imageMap["image"] = valMap
	msg = append(msg, imageMap)
	return msg, nil
}
