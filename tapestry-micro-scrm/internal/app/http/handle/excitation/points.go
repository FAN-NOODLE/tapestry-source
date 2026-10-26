package excitation

import (
	"github.com/gin-gonic/gin"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/app/http/handle"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationSvc "micro-scrm/internal/app/service/excitation"
)

// PointsHandle 导购积分模块
type PointsHandle struct {
	handle.BaseHandle
	pointsProduceSvc *excitationSvc.PointsProduceService
	pointsRecordSvc  *excitationSvc.PointsRecordService
}

func NewPointsHandle() PointsHandle {
	return PointsHandle{
		pointsProduceSvc: excitationSvc.NewPointsProduceService(),
		pointsRecordSvc:  excitationSvc.NewPointsRecordService(),
	}
}

// AddManualPoints  手动添加积分
func (h PointsHandle) AddManualPoints(c *gin.Context) {
	var request excitationReq.ManualPointsReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	var userIds []string
	for _, data := range request.MedalPointsData {
		userIds = append(userIds, data.UserId)
	}
	members, err := membersDao.NewMembers().GetMany(h.WarpContext(c), map[string]interface{}{
		"enterprise_id": request.EnterpriseId,
		"userid":        userIds,
	})
	if err != nil {
		return
	}

	var dataList []excitationReq.ManualPointsSave //要写入的积分数据
	var userToMember = make(map[string]int64)
	for _, member := range members {
		userToMember[member.Userid] = member.ID
	}
	for _, data := range request.MedalPointsData {
		dataList = append(dataList, excitationReq.ManualPointsSave{
			MemberId:       userToMember[data.UserId],
			RuleDetailName: data.RuleDetailName,
			RuleId:         data.RuleId,
			PointsRemark:   data.PointsRemark,
			Points:         data.Points,
			HasCoins:       data.HasCoins,
		})
	}

	err = h.pointsProduceSvc.ManualPoints(h.WarpContext(c), request.EnterpriseId, dataList)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

// RecordList   获取导购积分获取记录列表
func (h PointsHandle) RecordList(c *gin.Context) {

	var request excitationReq.RecordListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	list, total, err := h.pointsRecordSvc.GetRecordList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, request.Page, request.Size, list)

	return
}

//GetPointsList 小程序端导购积分列表
func (h PointsHandle) GetPointsList(c *gin.Context) {
	var request excitationReq.GetPointsListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	list, total, err := h.pointsRecordSvc.GetPointsList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, request.Page, request.Size, list)

	return
}
