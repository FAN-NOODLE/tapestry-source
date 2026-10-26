package customer

import (
	"github.com/gin-gonic/gin"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	"micro-scrm/internal/app/http/handle"
	customerReq "micro-scrm/internal/app/http/request/customer"
)

type RelationHandle struct {
	handle.BaseHandle
}

func NewRelationHandle() RelationHandle {
	return RelationHandle{}
}

func (h RelationHandle) GetUserMemberRelation(c *gin.Context) {
	var request customerReq.UserMemberRelationReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	many, err := microScrmDao.NewUserRelation().GetMany(h.WarpContext(c), map[string]interface{}{
		"enterprise_id":         request.EnterpriseID,
		"member_id":             request.MemberIds,
		"work_external_user_id": request.WorkExternalUserIds,
	})
	if err != nil {
		h.ResError(c, err)
		return
	}
	var resp []customerReq.UserMemberRelationResp
	for _, model := range many {
		respTmp := customerReq.UserMemberRelationResp{
			MemberId:           model.MemberID,
			WorkExternalUserId: model.WorkExternalUserID,
			RelationStatus:     model.RelationStatus,
		}
		resp = append(resp, respTmp)
	}

	h.ResSuccess(c, resp)
	return
}
