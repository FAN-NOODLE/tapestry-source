package display

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"gw-work-scrm/internal/app/http/handle"
	request2 "gw-work-scrm/internal/app/http/request"
	"gw-work-scrm/internal/app/http/request/display"
	"gw-work-scrm/internal/config"
	display2 "gw-work-scrm/internal/pkg/sdk/microscrm/display"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type ApiHandle struct {
	handle.BaseHandle
	apiSdk *display2.Api
}

func NewApiHandle() *ApiHandle {
	return &ApiHandle{
		apiSdk: display2.NewApi(),
	}
}

func (h ApiHandle) TaskInfo(c *gin.Context) {
	var request display.MemberStoreTaskInfoReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	info, err := h.apiSdk.TaskInfo(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)

}

func (h ApiHandle) Save(c *gin.Context) {
	var request display.StoreTaskSaveInfoReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	request.MemberName = memberInfo.Name
	request.MemberUserId = memberInfo.Userid
	info, err := h.apiSdk.SaveTask(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) TaskList(c *gin.Context) {
	var request display.MemberTaskListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	info, err := h.apiSdk.TaskList(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) Look(c *gin.Context) {
	var request display.LookReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	info, err := h.apiSdk.Look(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) Auth(c *gin.Context) {
	var request display.AuthReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	info, err := h.apiSdk.Auth(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) VmList(c *gin.Context) {
	var request request2.EnterpriseBaseReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = h.EnterpriseID(c)
	info, err := h.apiSdk.VmList(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) StoreList(c *gin.Context) {
	var request display.StoreSearchReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = h.EnterpriseID(c)
	list, err := h.apiSdk.StoreList(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h ApiHandle) ListTips(c *gin.Context) {
	var request display.ListTipsReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	request.MemberUserId = memberInfo.Userid
	info, err := h.apiSdk.ListTips(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) CheckList(c *gin.Context) {
	var request display.MemberTaskListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	request.MemberUserId = memberInfo.Userid
	info, err := h.apiSdk.CheckList(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) CheckStoreList(c *gin.Context) {
	var request display.StoreListInfoReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	request.MemberUserId = memberInfo.Userid
	info, err := h.apiSdk.CheckStoreList(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) StoreTaskInfo(c *gin.Context) {
	var request display.MemberStoreTaskInfoReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	info, err := h.apiSdk.StoreTaskInfo(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) VmCheck(c *gin.Context) {
	var request display.StoreTaskReviewReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	request.MemberUserId = memberInfo.Userid
	request.MemberName = memberInfo.Name
	info, err := h.apiSdk.VmCheck(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) SendNotice(c *gin.Context) {
	var request display.SendNoticeReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	info, err := h.apiSdk.VmSendNotice(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) WeLearnToken(c *gin.Context) {
	var request display.WeLearnToken
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if request.LoginStr == "" {
		h.ResError(c, errors.New("login_str 不能为空"))
		return
	}
	loginStr := request.LoginStr
	configInfo := config.CfgData
	salt := configInfo.WeLearningTokenSalt
	hasher := sha256.New()
	hasher.Write([]byte(salt + loginStr))
	hashInBytes := hasher.Sum(nil)

	// 将哈希结果转为十六进制字符串
	hashString := hex.EncodeToString(hashInBytes)
	h.ResSuccess(c, gin.H{
		"login_token": hashString,
	})

}

func (h ApiHandle) SaveImage(c *gin.Context) {
	var req display.SaveStoreImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.EnterpriseId = h.EnterpriseID(c)
	info, err := h.apiSdk.SaveImage(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) DelImage(c *gin.Context) {
	var req display.DelImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.EnterpriseId = h.EnterpriseID(c)
	info, err := h.apiSdk.DelImage(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) SaveComment(c *gin.Context) {
	var req display.SaveStoreCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.EnterpriseId = h.EnterpriseID(c)
	info, err := h.apiSdk.SaveComment(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) DelComment(c *gin.Context) {
	var req display.DelImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.EnterpriseId = h.EnterpriseID(c)
	info, err := h.apiSdk.DelComment(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) AllList(c *gin.Context) {
	var request request2.EnterpriseBaseReq
	request.EnterpriseId = h.EnterpriseID(c)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.apiSdk.AllList(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) DashBoard(c *gin.Context) {
	var req display.DashBoard
	req.EnterpriseId = h.EnterpriseID(c)
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.apiSdk.Dashboard(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) FailList(c *gin.Context) {
	var req display.MemberTaskReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.apiSdk.FailList(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) DashBoardList(c *gin.Context) {
	var req display.DashBoardListReq
	req.EnterpriseId = h.EnterpriseID(c)
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.apiSdk.DashboardList(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) StoreTaskList(c *gin.Context) {
	var req display.StoreListReq
	req.EnterpriseId = h.EnterpriseID(c)
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.apiSdk.StoreTaskList(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h ApiHandle) SpecialSave(c *gin.Context) {
	var request display.StoreTaskReviewReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	request.MemberUserId = memberInfo.Userid
	request.MemberName = memberInfo.Name
	info, err := h.apiSdk.SpecialSave(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}
