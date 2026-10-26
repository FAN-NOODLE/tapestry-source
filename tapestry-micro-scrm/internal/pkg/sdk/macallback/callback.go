package macallback

import (
	"strconv"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/sdk/cgema"
	"micro-scrm/internal/pkg/sdk/microwork"
)

func SendResult(c ctx.Context, originType int, request microwork.SendMaResultReq) error {
	if originType != microScrmModel.SHENYAN_TASK_ORIGIN_TYPE {
		_, err := microwork.NewMaPlan().PlanSendResult(c, request)
		return err
	}

	crowdResult := make([]cgema.CrowdResultItem, 0, len(request.CrowdResult))
	for _, item := range request.CrowdResult {
		crowdResult = append(crowdResult, cgema.CrowdResultItem{
			Date:                item.Date,
			SaName:              item.SaName,
			SaCode:              item.SaCode,
			FlowID:              item.FlowId,
			Campaign:            item.Campaign,
			CTaskRemark:         item.CTaskRemark,
			CTaskClassification: item.CTaskClassifation,
			UnionID:             item.Unionid,
			Result:              item.Result,
			ResultMessage:       item.ResultMessage,
			TaskRemark:          item.TaskRemark,
		})
	}

	_, err := cgema.NewClient().SendResult(c, cgema.SendResultRequest{
		TaskID:       strconv.FormatInt(request.TaskId, 10),
		CrowdID:      request.CrowdId,
		TaskCategory: request.TaskCategory,
		CrowdResult:  crowdResult,
	})
	return err
}

func SendFinish(c ctx.Context, request microwork.SendMaFinishReq) error {
	if request.OriginType != microScrmModel.SHENYAN_TASK_ORIGIN_TYPE {
		_, err := microwork.NewMaPlan().PlanSendFinish(c, request)
		return err
	}

	events := make([]cgema.EventItem, 0, len(request.Events))
	for _, item := range request.Events {
		events = append(events, cgema.EventItem{
			IdentityType:        item.IdentityType,
			IdentityValue:       item.IdentityValue,
			Date:                item.Date,
			CResult:             item.CResult,
			CSaName:             item.CSaName,
			CSaCode:             item.CSaCode,
			CTaskID:             item.CTaskId,
			CTaskName:           item.CTaskName,
			CCrowdID:            item.CCrowdId,
			FlowID:              item.FlowId,
			Campaign:            item.Campaign,
			CTaskRemark:         item.CTaskRemark,
			CTaskClassification: item.CTaskClassifation,
		})
	}

	_, err := cgema.NewClient().SendFinish(c, cgema.SendFinishRequest{
		Events: events,
		CommonProps: cgema.CommonProps{
			Event: request.CommonProps.Event,
		},
	})
	return err
}

func SendMomentResult(c ctx.Context, originType int, request microwork.SendMaMomentResultReq) error {
	if originType != microScrmModel.SHENYAN_TASK_ORIGIN_TYPE {
		_, err := microwork.NewMaPlan().PlanSendMomentResult(c, request)
		return err
	}

	_, err := cgema.NewClient().SendMomentResult(c, cgema.SendMomentResultRequest{
		SaCode:              request.SaCode,
		TaskID:              strconv.FormatInt(request.TaskId, 10),
		CrowdID:             request.CrowdId,
		Campaign:            request.Campaign,
		CustomerUnionIDList: request.CustomerUnionidList,
	})
	return err
}

func SendMomentComments(c ctx.Context, originType int, request microwork.SendMaMomentCommentsReq) error {
	if originType != microScrmModel.SHENYAN_TASK_ORIGIN_TYPE {
		_, err := microwork.NewMaPlan().PlanSendMomentComment(c, request)
		return err
	}

	commentList := make([]cgema.MomentCommentItem, 0, len(request.CommentList))
	for _, item := range request.CommentList {
		commentList = append(commentList, cgema.MomentCommentItem{
			UnionID:    item.UnionId,
			CreateTime: item.CreateTme,
		})
	}
	likeList := make([]cgema.MomentCommentItem, 0, len(request.LikeList))
	for _, item := range request.LikeList {
		likeList = append(likeList, cgema.MomentCommentItem{
			UnionID:    item.UnionId,
			CreateTime: item.CreateTme,
		})
	}

	_, err := cgema.NewClient().SendMomentComments(c, cgema.SendMomentCommentsRequest{
		SaCode:      request.SaCode,
		TaskID:      strconv.FormatInt(request.TaskId, 10),
		CrowdID:     request.CrowdId,
		Campaign:    request.Campaign,
		CommentList: commentList,
		LikeList:    likeList,
	})
	return err
}
