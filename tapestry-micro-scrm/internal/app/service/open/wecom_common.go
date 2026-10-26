package openSvc

import (
	"micro-scrm/internal/app"
	openReq "micro-scrm/internal/app/http/request/open"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
)

type WecomCommonService struct {
}

func NewWecomCommonService() *WecomCommonService {
	return &WecomCommonService{}
}

// ConvertOpenUserIdToUserId 将企微代开发应用的open_userid通过自建应用转成userid
func (s *WecomCommonService) ConvertOpenUserIdToUserId(c app.Ctx, req openReq.ConvertOpenUserIdRequest) (qyWork.OpenUseridToUseridResp, error) {
	return qyWork.NewQyWork(req.EnterpriseID).OpenUseridToUserid(c, qyWork.OpenUseridToUseridReq{
		OpenUserIDList: req.OpenUseridList,
		SourceAgentID:  req.SourceAgentId,
	})
}

// ConvertOpenExternalUseridToExternalUserid 将企微代开发应用的open_external_userid通过自建应用转成external_userid
func (s *WecomCommonService) ConvertOpenExternalUseridToExternalUserid(c app.Ctx, req openReq.ConvertOpenExternalUserIdRequest) (qyWork.OpenExternalUseridToExternalUseridResp, error) {
	return qyWork.NewQyWork(req.EnterpriseID).OpenExternalUseridToExternalUserid(c, qyWork.OpenExternalUseridToExternalUseridReq{
		ExternalUserid: req.ExternalUserid,
		SourceAgentID:  req.SourceAgentId,
	})
}
