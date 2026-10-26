package chat

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/pkg/uid"
)

// 媒体文件服务
type MsgMediaService struct {
	repo MsgMediaRepository
}

// 媒体文件服务新增
func (s *MsgMediaService) Add(c ctx.Context, row MsgMedia) error {
	if row.ID == 0 {
		row.ID = uid.ID(c)
	}
	_, err := s.repo.Add(c, row)
	if err != nil {
		log.Error(ctx.Wrap(c)).Msgf("新增媒体文件记录错误:" + err.Error())

		return err
	}

	return nil
}

// ListByFileId 根据媒体id列表查列表
func (s MsgMediaService) ListByFileId(c ctx.Context, enterpriseID int64, sdkFileidList []string) []MsgMedia {
	list, err := s.repo.ListByFileId(c, enterpriseID, sdkFileidList)
	if err != nil {
		log.Error(ctx.Wrap(c)).Msgf("根据媒体id列表查列表错误:" + err.Error())

		return []MsgMedia{}
	}

	return list
}
