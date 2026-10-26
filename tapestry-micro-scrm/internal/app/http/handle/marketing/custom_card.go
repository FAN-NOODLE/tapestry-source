package marketing

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	templateReq "micro-scrm/internal/app/http/request/template"
	marketingService "micro-scrm/internal/app/service/marketing"
)

type CustomCard struct {
	handle.BaseHandle
}

func NewCustomCard() CustomCard {
	return CustomCard{}
}

func (h CustomCard) Lists(context *gin.Context) {
	var Req templateReq.ReqList
	if err := context.ShouldBind(&Req); err != nil {
		h.ResError(context, err)
		return
	}
	var result templateReq.ResponseList
	Total, countErr := marketingService.NewCustomCard().GetCount(h.WarpContext(context), Req)
	if countErr != nil {
		h.ResError(context, countErr)
		return
	}
	if Total < 1 {
		result.List = []templateReq.TemplateInfo{}
		result.Total = Total
		h.ResSuccess(context, result)
		return
	}
	list, err := marketingService.NewCustomCard().GetTemplateList(h.WarpContext(context), Req)
	if err != nil {
		h.ResError(context, err)
		return
	}
	result.List = list
	result.Total = Total
	h.ResSuccess(context, result)
}

func (h CustomCard) Create(context *gin.Context) {
	var templateInfo templateReq.ReqCreate
	if err := context.ShouldBind(&templateInfo); err != nil {
		h.ResError(context, err)
		return
	}
	if templateInfo.ID > 0 {
		if err := marketingService.NewCustomCard().UpdateInfo(h.WarpContext(context), templateInfo); err != nil {
			h.ResError(context, err)
			return
		}
	} else {
		if err := marketingService.NewCustomCard().Create(h.WarpContext(context), templateInfo); err != nil {
			h.ResError(context, err)
			return
		}
	}
	h.ResOk(context)
}

func (h CustomCard) Delete(context *gin.Context) {
	var RequestDel templateReq.RequestInfo
	if err := context.ShouldBind(&RequestDel); err != nil {
		h.ResError(context, err)
		return
	}
	if err := marketingService.NewCustomCard().Delete(h.WarpContext(context), RequestDel); err != nil {
		h.ResError(context, err)
		return
	}
	h.ResOk(context)
}

func (h CustomCard) Check(context *gin.Context) {
	var RequestDel templateReq.RequestInfo
	if err := context.ShouldBind(&RequestDel); err != nil {
		h.ResError(context, err)
		return
	}
	total, err := marketingService.NewCustomCard().Check(h.WarpContext(context), RequestDel)
	if err != nil {
		h.ResError(context, err)
		return
	}
	var result templateReq.ResponseCheck
	result.Check = true
	if total > int64(0) {
		result.Check = false
	}
	h.ResSuccess(context, result)
}

func (h CustomCard) Detail(context *gin.Context) {
	var RequestDel templateReq.RequestInfo
	if err := context.ShouldBind(&RequestDel); err != nil {
		h.ResError(context, err)
		return
	}
	Info, err := marketingService.NewCustomCard().Detail(h.WarpContext(context), RequestDel)
	if err != nil {
		h.ResError(context, err)
		return
	}
	h.ResSuccess(context, Info)
}
