package groupmember

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"

type Service struct {
	repo Repository
}

func (ser Service) SelectGroupNum(c ctx.Context, param GroupParam) (int64, error) {
	return ser.repo.SelectGroupNum(c, param)
}

func (ser Service) SelectReGroupUv(c ctx.Context, param GroupParam) (int64, error) {
	return ser.repo.SelectReGroupUv(c, param)
}

func (ser Service) SelectAddGroupUv(c ctx.Context, param GroupParam) (int64, error) {
	return ser.repo.SelectAddGroupUv(c, param)
}

func (ser Service) SelectLossGroupUv(c ctx.Context, param GroupParam) (int64, error) {
	return ser.repo.SelectLossGroupUv(c, param)
}
