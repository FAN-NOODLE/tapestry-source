// Package callback 对接Open 平台
package guide

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

//go:generate mockgen -source=${GOFILE} -destination=mocks/mock_${GOFILE}

// Repository 存储定义
type Repository interface {
	repos.IRepository
	AddUserBingGuide(ctx.Context, *BasicUserBingGuideInfoEntity) (int64, error)
	UpdateUserBingGuide(ctx.Context, *BasicUserBingGuideInfoEntity) (int64, error)
	CountUserBingGuide(ctx.Context, *BasicUserBingGuideInfoEntity) int64
	DeleteUserBingGuide(ctx.Context, *BasicUserBingGuideInfoEntity) (int64, error)
	DeleteBatchUserBingGuide(ctx.Context, *BatchUserBingGuideInfoEntity) error
	SelectUserBingGuide(ctx.Context, *BatchUserBingGuideInfoEntity) (*[]BasicUserBingGuideInfoEntity, error)
	FindGuideRelationByUserID(cc ctx.Context, enterpriseID int64, userID int64) (*CustomerGuideRelation, error)
}
