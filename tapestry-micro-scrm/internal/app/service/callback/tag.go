package callbackSvc

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/sdk/work"
)

var (
	// UpdateTag 更新标签事件
	UpdateTag = "update_tag"
)

type Tag struct {
	service.Base
	bossDetail boss.EnterpriseModel
}

func NewTag() *Tag {
	return &Tag{}
}

func (s *Tag) Handel(c ctx.Context, data work.ContactCallbackTag, bossDetail boss.EnterpriseModel) error {
	s.Log().Info(c).Msgf("通讯录回调,事件类型：%s,标签事件数据：%v,企业数据：%v", data.ChangeType, data, bossDetail)
	return nil
}
