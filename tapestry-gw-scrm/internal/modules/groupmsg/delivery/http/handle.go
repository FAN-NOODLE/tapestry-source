package groupmsg_http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/groupmsg/api"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/excelizeLib"
	"gw-scrm/internal/pkg/helper"
	"gw-scrm/internal/pkg/tencent/tms"
	"net/http"
)

// list 群发消息列表
func list(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.GroupMsgApi{}
	res, rErr := act.List(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// groupSendAdd 新增群发消息
func groupSendAdd(c *gin.Context) {
	var (
		tokenInfo *common.TokenInfo
		act       = &api.GroupMsgApi{}
		params    map[string]interface{}
		res       interface{}
		err       error
	)
	params = common.GetRawMap(c)
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	// 处理非法关键字
	content := params["content"].(map[string]interface{})
	text := content["text"].(map[string]interface{})
	tempText := text["content"].(string)
	if len(tempText) > 0 {
		apiResponse, ok, err := tms.Verification(text["content"].(string))
		if err != nil {
			ginplus.ResError(c, err, http.StatusOK)
			return
		}
		if !ok {
			resp := &ginplus.APIResponse{
				Code:      ginplus.ParamsErrCode,
				Message:   "存在非法关键字",
				RequestID: ctx.Wrap(c).ID(),
				Data:      apiResponse,
			}
			ginplus.ResJSON(c, http.StatusOK, resp)
			return
		}
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	res, err = act.GroupSendAdd(ctx.Wrap(c), params)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

// groupSendEdit编辑群发消息
func groupSendEdit(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	// 处理非法关键字
	content := params["content"].(map[string]interface{})
	text := content["text"].(map[string]interface{})
	tempText := text["content"].(string)
	if len(tempText) > 0 {
		apiResponse, ok, err := tms.Verification(text["content"].(string))
		if err != nil {
			ginplus.ResError(c, err, http.StatusOK)
			return
		}
		if !ok {
			resp := &ginplus.APIResponse{
				Code:      ginplus.ParamsErrCode,
				Message:   "存在非法关键字",
				RequestID: ctx.Wrap(c).ID(),
				Data:      apiResponse,
			}
			ginplus.ResJSON(c, http.StatusOK, resp)
			return
		}
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.GroupMsgApi{}
	res, rErr := act.GroupSendEdit(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// memberList任务明细-执行群主下拉列表
func memberList(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.GroupMsgApi{}
	res, rErr := act.MemberList(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// TaskInfo 群发任务详情
func TaskInfo(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.GroupMsgApi{}
	res, rErr := act.TaskInfo(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)

		return
	}
	ginplus.ResSuccess(c, res)
}

// TaskDetail 群发任务明细
func TaskDetail(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.GroupMsgApi{}
	res, rErr := act.TaskDetail(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)
		return
	}
	ginplus.ResSuccess(c, res)
}

func TaskMemberDetail(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	params["enterprise_id"] = tokenInfo.EnterpriseId
	act := &api.GroupMsgApi{}
	res, rErr := act.TaskMemberDetail(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)
		return
	}
	ginplus.ResSuccess(c, res)
}

func TaskMemberSendList(c *gin.Context) {
	params := common.GetRawMap(c)
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	params["enterprise_id"] = tokenInfo.EnterpriseId
	res, rErr := new(api.GroupMsgApi).TaskMemberSendList(ctx.Wrap(c), params)
	if rErr != nil {
		ginplus.ResError(c, rErr)
		return
	}

	if res.Size != 0 && res.Page != 0 {
		ginplus.ResSuccess(c, res)
		return
	}

	// 数据导出
	data := make([]map[string]interface{}, 0)
	var isExternalUser bool
	for _, send := range res.List {
		if send.ExternalUserID != "" {
			isExternalUser = true
		}
		tmp := map[string]interface{}{
			"status":    convertStatus(send.Status),
			"send_time": "-",
		}
		if send.ExternalUserID != "" {
			tmp["id"] = send.ExternalUserID
		}
		if send.ChatId != "" {
			tmp["id"] = send.ChatId
		}
		if send.SendTime > 0 {
			tmp["send_time"] = helper.GetFormatDate(send.SendTime, 4)
		}
		if send.ChatId != "" {
			tmp["name"] = send.ChatName
		}
		if send.ExternalUserID != "" {
			tmp["name"] = send.UserName
		}
		data = append(data, tmp)
	}

	dataKey := make([]map[string]string, 0)
	if isExternalUser {
		dataKey = append(dataKey, map[string]string{"key": "id", "title": "客户ID", "width": "50", "is_num": "0"})
		dataKey = append(dataKey, map[string]string{"key": "name", "title": "顾客昵称", "width": "20", "is_num": "0"})
		dataKey = append(dataKey, map[string]string{"key": "status", "title": "发送状态", "width": "30", "is_num": "0"})
		dataKey = append(dataKey, map[string]string{"key": "send_time", "title": "发送时间", "width": "20", "is_num": "0"})
	} else {
		dataKey = append(dataKey, map[string]string{"key": "id", "title": "客户群ID", "width": "50", "is_num": "0"})
		dataKey = append(dataKey, map[string]string{"key": "name", "title": "客户群名称", "width": "20", "is_num": "0"})
		dataKey = append(dataKey, map[string]string{"key": "status", "title": "发送状态", "width": "30", "is_num": "0"})
		dataKey = append(dataKey, map[string]string{"key": "send_time", "title": "发送时间", "width": "20", "is_num": "0"})
	}
	excel := excelizeLib.NewMyExcel()
	excel.ExportToWeb(c, dataKey, data, "群发明细-")
	return
}

func convertStatus(status int) string {
	switch status {
	case 0:
		return "未发送"
	case 1:
		return "已发送"
	case 2:
		return "因客户不是好友导致发送失败"
	case 3:
		return "因客户已经收到其他群发消息导致发送失败"
	}
	return ""
}
