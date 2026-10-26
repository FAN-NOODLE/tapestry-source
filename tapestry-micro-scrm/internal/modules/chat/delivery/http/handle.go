package chat_http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/chat"
	"micro-scrm/internal/pkg/common"
)

var validate *validator.Validate
var chatService *chat.Service
var memberService *chat.MemberService
var objectService *chat.ObjectService
var msgMediaService *chat.MsgMediaService

// msgTypeList 获取会话存档同意信息
func msgTypeList(c *gin.Context) {
	ginplus.ResSuccess(c, map[string]interface{}{
		"list": chatService.MsgTypeList(),
	})

	return
}

// agreeInfo 获取会话存档同意信息
func agreeInfo(c *gin.Context) {
	req := chat.AgreeInfoRequest{
		EnterpriseID: common.EnterpriseID(ctx.Wrap(c)),
	}
	_ = c.ShouldBindJSON(&req)
	err := validate.Struct(req)
	if err != nil {
		ginplus.ResError(c, err)

		return
	}

	res, resErr := chatService.GetCheckSingleAgree(ctx.Wrap(c), req)
	if resErr != nil {
		ginplus.ResError(c, resErr)

		return
	}

	ginplus.ResSuccess(c, map[string]interface{}{
		"list": res,
	})

	return
}

// msgList 消息分页列表
func msgList(c *gin.Context) {
	req := chat.MsgPageListRequest{
		EnterpriseID: common.EnterpriseID(ctx.Wrap(c)),
	}

	_ = c.ShouldBindJSON(&req)
	err := validate.Struct(req)
	if err != nil {
		ginplus.ResError(c, err)

		return
	}

	page, pErr := chatService.MsgPageList(ctx.Wrap(c), *msgMediaService, req)
	if pErr != nil {
		ginplus.ResError(c, pErr)

		return
	}

	ginplus.ResSuccess(c, page)

	return
}

// groupMsgList 群消息分页列表
func groupMsgList(c *gin.Context) {
	req := chat.GroupMsgListRequest{
		EnterpriseID: common.EnterpriseID(ctx.Wrap(c)),
	}

	_ = c.ShouldBindJSON(&req)
	err := validate.Struct(req)
	if err != nil {
		ginplus.ResError(c, err)

		return
	}

	page, pErr := chatService.GroupMsgPageList(ctx.Wrap(c), *msgMediaService, req)
	if pErr != nil {
		ginplus.ResError(c, pErr)

		return
	}

	ginplus.ResSuccess(c, page)

	return
}

// memberList 员工分页列表
func memberList(c *gin.Context) {
	req := chat.MemberPageListRequest{
		EnterpriseID: common.EnterpriseID(ctx.Wrap(c)),
	}
	_ = c.ShouldBindJSON(&req)
	err := validate.Struct(req)
	if err != nil {
		ginplus.ResError(c, err)

		return
	}
	page, pErr := memberService.MemberPageList(ctx.Wrap(c), req)
	if pErr != nil {
		ginplus.ResError(c, pErr)

		return
	}

	ginplus.ResSuccess(c, page)

	return
}

// msgListByID 根据消息ID找附近的消息
func msgListByID(c *gin.Context) {
	req := chat.MsgListByIDRequest{
		EnterpriseID: common.EnterpriseID(ctx.Wrap(c)),
	}
	_ = c.ShouldBindJSON(&req)
	err := validate.Struct(req)
	if err != nil {
		ginplus.ResError(c, err)

		return
	}
	page, pErr := chatService.MsgListByID(ctx.Wrap(c), req, *msgMediaService)
	if pErr != nil {
		ginplus.ResError(c, pErr)

		return
	}

	ginplus.ResSuccess(c, page)

	return
}

// memberToMemberList 员工与员工聊天分页列表
func memberToMemberList(c *gin.Context) {
	req := chat.MemberToMemberPageListRequest{
		EnterpriseID: common.EnterpriseID(ctx.Wrap(c)),
	}
	_ = c.ShouldBindJSON(&req)
	err := validate.Struct(req)
	if err != nil {
		ginplus.ResError(c, err)

		return
	}
	page, pErr := objectService.MemberToMemberList(ctx.Wrap(c), req)
	if pErr != nil {
		ginplus.ResError(c, pErr)

		return
	}

	ginplus.ResSuccess(c, page)

	return
}

// memberToExternalList 员工与外部联系人聊天分页列表
func memberToExternalList(c *gin.Context) {
	req := chat.MemberToExternalPageListRequest{
		EnterpriseID: common.EnterpriseID(ctx.Wrap(c)),
	}
	_ = c.ShouldBindJSON(&req)
	err := validate.Struct(req)
	if err != nil {
		ginplus.ResError(c, err)

		return
	}
	page, pErr := objectService.MemberToExternalList(ctx.Wrap(c), req)
	if pErr != nil {
		ginplus.ResError(c, pErr)

		return
	}

	ginplus.ResSuccess(c, page)

	return
}

// memberToGroupList 员工群聊天分页列表
func memberToGroupList(c *gin.Context) {
	req := chat.MemberToGroupPageListRequest{
		EnterpriseID: common.EnterpriseID(ctx.Wrap(c)),
	}
	_ = c.ShouldBindJSON(&req)
	err := validate.Struct(req)
	if err != nil {
		ginplus.ResError(c, err)

		return
	}
	page, pErr := objectService.MemberToGroupList(ctx.Wrap(c), req)
	if pErr != nil {
		ginplus.ResError(c, pErr)

		return
	}

	ginplus.ResSuccess(c, page)

	return
}

// getUserInfo 获取用户信息
func getUserInfo(c *gin.Context) {
	var req chat.UserInfo
	if err := c.ShouldBind(&req); err != nil {
		ginplus.ResError(c, err)
		return
	}

	userType := chatService.ConversionUserType(req.UserID)

	userIds := []string{req.UserID}

	list, err := chatService.GetUserInfo(ctx.Wrap(c), userType, userIds, req.EnterpriseID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if len(*list) == 0 {
		ginplus.ResSuccess(c, req)
		return
	}
	ginplus.ResSuccess(c, (*list)[0])
	return
}

func test(c *gin.Context) {
	var a int64
	a = 100
	var b uint
	b = uint(a)

	fmt.Println("bbbb", a, b)

	ginplus.ResSuccess(c, map[string]interface{}{
		"data": b,
	})
	return
	/* 代码质量检测不通过
	msg, err := chatService.Detail(ctx.Wrap(c), 120584607814336)
	if err != nil {
		ginplus.ResError(c, err)

		return
	}

	ginplus.ResSuccess(c, map[string]interface{}{
		"msg": msg,
		"fileid_arr": chatService.GetSdkFileId(msg.Type, msg.Content),
	})
	return

	*/
}
