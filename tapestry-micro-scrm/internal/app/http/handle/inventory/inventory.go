package inventory

import (
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	"micro-scrm/internal/app/http/handle"
	inventoryReq "micro-scrm/internal/app/http/request/inventory"
	inventoryResp "micro-scrm/internal/app/http/response/inventory"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	"micro-scrm/internal/app/service/department"
	inventorySvc "micro-scrm/internal/app/service/inventory"
	"strings"

	"github.com/gin-gonic/gin"
)

// 积分规则

type InvHandle struct {
	handle.BaseHandle
	InventorySvc *inventorySvc.InventoryService
}

func firstStoreCode(storeCodes []string) string {
	if len(storeCodes) == 0 {
		return ""
	}
	return storeCodes[0]
}

func buildUpcStoreCodeParams(storeCodeParams []string) []string {
	if len(storeCodeParams) == 1 && storeCodeParams[0] == "OC226" {
		return storeCodeParams
	}
	return nil
}

func buildUpcResponseStoreCode(storeCodeParams []string) string {
	if len(storeCodeParams) == 1 && storeCodeParams[0] != "OC226" {
		return storeCodeParams[0]
	}
	return ""
}

func filterUpcSkuListByStore(skuList []inventoryResp.SkuInfo, upcCode, storeCode string) []inventoryResp.SkuInfo {
	if storeCode == "" || len(skuList) == 0 {
		return skuList
	}
	for _, sku := range skuList {
		if sku.LocCode == storeCode {
			return []inventoryResp.SkuInfo{sku}
		}
	}
	return []inventoryResp.SkuInfo{
		{
			UpcInfo: inventoryResp.UpcInfo{
				UpcCode:   upcCode,
				LocCode:   storeCode,
				DclOnHand: skuList[0].DclOnHand,
			},
		},
	}
}

func NewInvHandle() InvHandle {
	return InvHandle{
		InventorySvc: inventorySvc.NewInventoryService(),
	}
}

// List   库存列表
func (h InvHandle) List(c *gin.Context) {
	var request inventoryReq.ListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	//allStoreCode, _ := department.NewDepartmentsService().GetAllStoreCode(h.WarpContext(c), request.EnterpriseId)

	var storeCodeParams []string
	if len(request.StoreCodeArr) > 0 {
		//如果传了专门的门店参数则查导购门店
		storeCodeParams = request.StoreCodeArr
	}
	//} else {
	//	//没传则用全部门店来查
	//	storeCodeParams = allStoreCode
	//}
	var storeType int32
	if len(storeCodeParams) == 1 {
		storeDepartmentInfo, _ := microWorkDao.NewStoreDepartment().GetInfo(h.WarpContext(c), map[string]any{
			"enterprise_id":    request.EnterpriseId,
			"third_store_code": storeCodeParams[0],
			"deleted_at":       0,
		})
		storeType = storeDepartmentInfo.StoreType
	}
	var (
		result = inventoryResp.ListResult{
			List:            make([]inventoryResp.SkuInfo, 0),
			InventoryStatus: inventorySvc.InventoryStatusSuccess,
			FailedUPCs:      make([]string, 0),
			DCLStatus:       inventorySvc.DCLStatusNotRequired,
		}
		err error
	)
	keywordType := ""
	searchKeyword := ""
	if request.SkuItem != "" {
		request.SkuItem = strings.ToUpper(request.SkuItem)
		searchKeyword = request.SkuItem
		keywordType = "sku_item"
		upcList, err := h.InventorySvc.GetUpcListBySkuItem(h.WarpContext(c), request.EnterpriseId, request.SkuItem)
		if err != nil {
			h.ResError(c, err)
			return
		}
		if len(upcList) == 0 {
			h.ResSuccess(c, result)
			return
		}
		//传所有门店编码相当于不传
		//通过sku_item查的时候需要把同渠道有库存的商品都查出来
		result, err = h.InventorySvc.FetchConcurrentlyBySkuItem(h.WarpContext(c), request.SkuItem, upcList, firstStoreCode(storeCodeParams), storeType)
		if err != nil {
			h.ResError(c, err)
			return
		}
	}

	if request.UpcCode != "" {
		searchKeyword = request.UpcCode
		keywordType = "upc_code"
		result, err = h.InventorySvc.GetListByUpcCode(h.WarpContext(c), request.UpcCode, buildUpcStoreCodeParams(storeCodeParams), storeType)
		if err != nil {
			h.ResError(c, err)
			return
		}
		result.List = filterUpcSkuListByStore(result.List, request.UpcCode, buildUpcResponseStoreCode(storeCodeParams))
	}

	if len(result.List) > 0 {
		var storeCodeSlice []string
		for _, sku := range result.List {
			if sku.LocCode != "" {
				storeCodeSlice = append(storeCodeSlice, sku.LocCode)
			}
		}
		var storeMap = make(map[string]microWorkModel.StoreDepartmentModel)
		if len(storeCodeSlice) > 0 {
			storeMap, err = department.NewDepartmentsService().GetDpListByStoreCode(h.WarpContext(c), request.EnterpriseId, storeCodeSlice)
		}
		if err != nil {
			return
		}
		for i, sku := range result.List {
			if val, ok := storeMap[sku.LocCode]; ok {
				if val.ClosedAt == 0 {
					result.List[i].StoreName = val.InventoryName
					result.List[i].City = val.City
					result.List[i].Phone = val.Phone
				}
			} else if sku.LocCode == "OC226" {
				result.List[i].StoreName = sku.LocCode
				result.List[i].InventoryName = "网仓"
			}
		}
		for i, sku := range result.List {
			if sku.LocCode != "" {
				if result.List[i].ThirdStoreCode == "" {
					result.List[i].ThirdStoreCode = sku.LocCode
				}
				if result.List[i].StoreName == "" {
					result.List[i].StoreName = sku.LocCode
				}
				if result.List[i].InventoryName == "" {
					result.List[i].InventoryName = sku.LocCode
				}
			}
		}
	}
	if searchKeyword != "" && keywordType != "" && !request.Scan && len(request.StoreCodeArr) > 0 {
		//如果传了skuItem或者upcCode则查询baozun_search表
		err = h.InventorySvc.InsertKeywords(h.WarpContext(c), searchKeyword, keywordType, request.EnterpriseId, request.MemberId)
		if err != nil {
			h.ResError(c, err)
			return
		}
		err = h.InventorySvc.CheckKeywordsNum(h.WarpContext(c), request.EnterpriseId, request.MemberId, keywordType)
		if err != nil {
			h.ResError(c, err)
			return
		}
	}
	h.ResSuccess(c, result)
	return
}

func (h InvHandle) DcList(c *gin.Context) {
	var request inventoryReq.DcListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	// 通过upCode查询大仓的库存
	skuList, err := h.InventorySvc.GetDcListByUpCode(h.WarpContext(c), request.UpcCode, request.EnterpriseId, request.MemberId)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, skuList)

}

func (h InvHandle) SearchKeywords(c *gin.Context) {
	var request inventoryReq.KeywordsReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	keywords, err := h.InventorySvc.GetSearchKeywords(h.WarpContext(c), request.EnterpriseId, request.MemberId)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, keywords)
}

func (h InvHandle) DelKeywords(c *gin.Context) {
	var request inventoryReq.DelKeywordsReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.InventorySvc.DelKeywords(h.WarpContext(c), request.EnterpriseId, request.MemberId, request.Ids, request.IsAll)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}
