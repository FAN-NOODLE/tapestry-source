package inventory

import (
	"errors"
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	inventoryReq "micro-scrm/internal/app/http/request/inventory"
	inventorySvc "micro-scrm/internal/app/service/inventory_sw"
)

// 积分规则

type InvSwHandle struct {
	handle.BaseHandle
	InventorySwSvc *inventorySvc.InventorySwService
}

func NewInvSwHandle() InvSwHandle {
	return InvSwHandle{
		InventorySwSvc: inventorySvc.NewInventorySwService(),
	}
}

// List   库存列表
func (h InvSwHandle) List(c *gin.Context) {
	var request inventoryReq.ListSwReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if (request.SearchOtherStore && request.Plu == "") || request.StoreCode == "" || (request.Plu == "" && request.StkCode == "" && request.StkName == "") {
		h.ResError(c, errors.New("参数错误"))
		return
	}
	res, err := h.InventorySwSvc.GetList(h.WarpContext(c), request.ParamsReq)
	if err != nil {
		h.ResSuccess(c, res)
		return
	}
	keywordType := ""
	searchKeyword := ""
	if request.StkName != "" {
		searchKeyword = request.StkName
		keywordType = "stk_name"
	} else if request.StkCode != "" {
		searchKeyword = request.StkCode
		keywordType = "stk_code"
	} else if request.Plu != "" {
		searchKeyword = request.Plu
		keywordType = "plu"
	}

	if keywordType != "" && !request.Scan && !request.SearchOtherStore {
		// 历史记录
		_ = h.InventorySwSvc.InsertKeywords(h.WarpContext(c), searchKeyword, keywordType, request.EnterpriseId, request.MemberId)
		_ = h.InventorySwSvc.CheckKeywordsNum(h.WarpContext(c), request.EnterpriseId, request.MemberId, keywordType)
	}
	h.ResSuccess(c, res)
	return
}

func (h InvSwHandle) SearchKeywords(c *gin.Context) {
	var request inventoryReq.KeywordsReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	keywords, err := h.InventorySwSvc.GetSearchKeywords(h.WarpContext(c), request.EnterpriseId, request.MemberId)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, keywords)
}

func (h InvSwHandle) DelKeywords(c *gin.Context) {
	var request inventoryReq.DelKeywordsSwReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.InventorySwSvc.DelKeywords(h.WarpContext(c), request.EnterpriseId, request.MemberId, request.Ids, request.DelType)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}
