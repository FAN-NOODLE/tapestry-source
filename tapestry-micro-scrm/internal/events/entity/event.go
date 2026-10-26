package entity

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"

// 事件监听的数据结构。
// 不一定每个字段用到
// 必填字段：EnterpriseID、Ctx
type Data struct {
	Ctx          ctx.Context
	EnterpriseID int64

	// 外部联系人标签变更事件
	ExternalUserID int64
	RelationID     int64

	// 客户分组删除事件
	CustomerGroupID int64
}
