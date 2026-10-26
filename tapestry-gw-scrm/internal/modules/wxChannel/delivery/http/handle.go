package wxchannelhttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

//wxChannelList 可用的微信场景值列表
func wxChannelList(c *gin.Context) {
	var request interface{}
	response, err := microscrm.WxChannelList(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}
