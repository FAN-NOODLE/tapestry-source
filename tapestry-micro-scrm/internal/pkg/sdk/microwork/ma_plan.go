package microwork

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"

type MaPlan struct {
}

func NewMaPlan() *MaPlan {
	return &MaPlan{}
}

// PlanSendResult 任务的发送结果回执
func (s MaPlan) PlanSendResult(c ctx.Context, request SendMaResultReq) (ResponseBase, error) {
	var response ResponseBase
	err := httpClient.Post(c, GetURL("%s/service/ma/send-result"), request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// PlanSendFinish 任务完成回执
func (s MaPlan) PlanSendFinish(c ctx.Context, request SendMaFinishReq) (ResponseBase, error) {
	var response ResponseBase
	err := httpClient.Post(c, GetURL("%s/service/ma/send-finish"), request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// PlanSendMomentResult  朋友圈任务发送结果及可见范围回执给到MA
func (s MaPlan) PlanSendMomentResult(c ctx.Context, request SendMaMomentResultReq) (ResponseBase, error) {
	var response ResponseBase
	err := httpClient.Post(c, GetURL("%s/service/ma/moment-send-result"), request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// PlanSendMomentComment  朋友圈任务评论及点赞回执给到MA
func (s MaPlan) PlanSendMomentComment(c ctx.Context, request SendMaMomentCommentsReq) (ResponseBase, error) {
	var response ResponseBase
	err := httpClient.Post(c, GetURL("%s/service/ma/moment-comments"), request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
