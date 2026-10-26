package user

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	userReq "gw-scrm/internal/app/http/request/user"
	userEntity "gw-scrm/internal/modules/user"
	"gw-scrm/internal/pkg/sdk/microauthority"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type SSOHandle struct {
	handle.BaseHandle
	UserApi microscrm.UserApi
}

func NewSSOHandle() SSOHandle {
	return SSOHandle{
		UserApi: microscrm.NewUserApi(),
	}
}

func (h SSOHandle) SSOLogin(c *gin.Context) {
	var req userReq.SSOLoginReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	appCode, ok := userEntity.NormalizeAppCode(req.AppCode)
	if !ok {
		h.ResError(c, errors.New("不支持的app_code"))
		return
	}
	req.AppCode = appCode
	userInfo, err := h.UserApi.SSOLogin(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	jsonStr, _ := json.Marshal(userInfo)
	resp := userReq.SSOLoginResp{}
	err = json.Unmarshal(jsonStr, &resp)
	if err != nil {
		return
	}
	if resp.Mobile == "" {
		h.ResError(c, errors.New("没有获取到用户"))
		return
	}

	res, err := microauthority.LoginByMobile(h.WarpContext(c), map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
		"client":        userEntity.ClientForAppCode(req.AppCode),
		"mobile":        resp.Mobile,
	})
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}
