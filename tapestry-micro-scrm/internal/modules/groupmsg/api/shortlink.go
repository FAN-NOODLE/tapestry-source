package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/pkg/common"
)

type ShortLinkApi struct {
}

// CreateShortParamID 生成短参
func (s *ShortLinkApi) CreateShortParamID(c ctx.Context, param common.Params) (int64, error) {
	paramID, err := common.CreateShortParamId(c, param)
	if err != nil {
		log.Info(c).Msgf("生成短参失败, err:%s, param:%v", err.Error(), param)
		return 0, err
	}
	return paramID, nil
}
