package sourceGroup

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

//go:generate mockgen -source=${GOFILE} -destination=mocks/mock_${GOFILE}

// Repository 存储定义
type Repository interface {
	repos.IRepository
	IsRepeatName(cc ctx.Context, entity *SourceGroup) (bool, error)
	BatchDel(cc ctx.Context, ids []int64) (int64, error)
}
