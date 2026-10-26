package open

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/middleware"
	openReq "micro-scrm/internal/app/http/request/open"
	messageSvc "micro-scrm/internal/app/service/message"
	openSvc "micro-scrm/internal/app/service/open"
	"micro-scrm/internal/pkg/sdk/microwork"
	"net/url"
	"unicode/utf8"
)

type MessageSendHandle struct {
	handle.BaseHandle
	openSvc *openSvc.MessageSendService
}

func NewMessageSendHandle() *MessageSendHandle {
	return &MessageSendHandle{
		openSvc: openSvc.NewMessageSendService(),
	}
}

func (h MessageSendHandle) SendWaitPayment(c *gin.Context) {
	var request openReq.SendWaitPaymentRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	request.EnterpriseID = c.GetInt64(middleware.EnterpriseId)

	resContent, _ := url.QueryUnescape(request.Content)
	if utf8.RuneCountInString(resContent) > 30 {
		resContent = string([]rune(resContent)[:27]) + "..."
	}

	content := make([]map[string]string, 0)
	content = append(content, map[string]string{"key": "订单编号", "value": request.OrderNo})
	content = append(content, map[string]string{"key": "温馨提示", "value": resContent})

	// 改业务了这里直接会传递员工编号
	/*userid, err := h.openSvc.GetOrderMemberUserId(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}*/

	params := microwork.SendMiniMessageRequest{
		EnterpriseID: request.EnterpriseID,
		MemUserIds:   []string{request.UserId},
		Title:        "线上订单即将过期提醒",
		Content:      content,
		Page:         "/pages/online-order/online-order?scene=list",
	}

	resp, err := messageSvc.NewMessageService().SendMiniMessage(h.WarpContext(c), params)

	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, resp)
}
