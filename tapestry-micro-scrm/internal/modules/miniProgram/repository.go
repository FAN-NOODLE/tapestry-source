// Package callback 对接Open 平台
package miniProgram

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

//go:generate mockgen -source=${GOFILE} -destination=mocks/mock_${GOFILE}

// Repository 存储定义
type Repository interface {
	repos.IRepository

	// AuthorizeMiniProgramDetail 授权小程序详情
	AuthorizeMiniProgramDetail(cc ctx.Context, enterpriseId int64) (int64, error)
}
