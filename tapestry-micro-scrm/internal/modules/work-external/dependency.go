package workexternal

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"

type EnterpriseService interface {
	GetCorpID(c ctx.Context, enterpriseID int64) (string, error)
}
