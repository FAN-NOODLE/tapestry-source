package moment_http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/moment"
	"gw-scrm/internal/modules/moment/api"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/errs"
	"gw-scrm/internal/pkg/tencent/tms"
	"net/http"
)

// topicCreate 朋友圈创建
func topicCreate(c *gin.Context) {
	params := &moment.TopicCreateRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	if len(params.Content["text"].(string)) > 0 {
		apiResponse, ok, err := tms.Verification(fmt.Sprintf("%s%s", params.Name, params.Content["text"].(string)))
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

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseId = tokenInfo.EnterpriseId
	params.CreatedMemberId = tokenInfo.AdminId
	act := &api.MomentApi{}
	res, err := act.TopicCreate(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, res)
}

// topicEdit 朋友圈编辑
func topicEdit(c *gin.Context) {
	params := &moment.TopicEditRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	if len(params.Content["text"].(string)) > 0 {
		apiResponse, ok, err := tms.Verification(params.Content["text"].(string))
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
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseId = tokenInfo.EnterpriseId
	act := &api.MomentApi{}
	res, err := act.TopicEdit(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, res)
}

// topicDelete 朋友圈删除
func topicDelete(c *gin.Context) {
	params := &moment.TopicDeleteRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	act := &api.MomentApi{}
	res, err := act.TopicDelete(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, res)
}

// topicDisplay 朋友圈展示
func topicDisplay(c *gin.Context) {
	params := &moment.TopicDisplayRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	act := &api.MomentApi{}
	res, err := act.TopicDisplay(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, res)
}

// memberTopic 一个员工的朋友圈列表
func memberTopic(c *gin.Context) {
	params := &moment.MemberTopicRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	act := &api.MomentApi{}
	res, err := act.MemberTopic(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, res)
}

// contentList 朋友圈按内容列表展示
func contentList(c *gin.Context) {
	params := &moment.ContentListRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseId = tokenInfo.EnterpriseId
	act := &api.MomentApi{}
	res, err := act.ContentList(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, res)
}

// memberList 朋友圈列表按员工搜索
func memberList(c *gin.Context) {
	params := &moment.MemberListRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseId = tokenInfo.EnterpriseId
	act := &api.MomentApi{}
	res, err := act.MemberList(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, res)
}

// appearanceGet 获取朋友圈形象
func appearanceGet(c *gin.Context) {
	params := &moment.MomentMemberGetRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
		params.EnterpriseId = tokenInfo.EnterpriseId
	}
	act := &api.MomentApi{}
	res, err := act.AppearanceGet(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, res)
}

// appearanceEdit 编辑朋友圈形象
func appearanceEdit(c *gin.Context) {
	params := &moment.MomentMemberEditRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseId = tokenInfo.EnterpriseId
	act := &api.MomentApi{}
	res, err := act.AppearanceEdit(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, res)
}

func jsConfig(c *gin.Context) {
	var request moment.JsConfigRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	act := &api.MomentApi{}
	res, err := act.JsConfig(ctx.Wrap(c), &request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

// topicMemberDelete 删除员工朋友圈
func topicMemberDelete(c *gin.Context) {
	params := &moment.TopicMemberDeleteRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseId = tokenInfo.EnterpriseId
	act := &api.MomentApi{}
	res, err := act.TopicMemberDelete(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, res)
}

// topicMemberEdit 编辑员工朋友圈
func topicMemberEdit(c *gin.Context) {
	params := &moment.TopicMemberEditRequest{}
	err := c.BindJSON(params)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	params.EnterpriseId = tokenInfo.EnterpriseId
	act := &api.MomentApi{}
	res, err := act.TopicMemberEdit(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, res)
}
