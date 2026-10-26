package inventorySvc

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	inventoryResp "micro-scrm/internal/app/http/response/inventory"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/sdk/inventory"
	"sort"
	"strings"
	"sync"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

type InventoryService struct {
	service.Base
	dpDao             *microWorkDao.StoreDepartment
	postInventory     func(ctx.Context, bool, inventory.ParamsReq, *inventory.Response, ...clients.Option) error
	postInventoryOnce func(ctx.Context, bool, inventory.ParamsReq, *inventory.Response, ...clients.Option) error
}

const dclLocCode = "CN10"

const (
	InventoryStatusSuccess = "success"
	InventoryStatusPartial = "partial"
	DCLStatusSuccess       = "success"
	DCLStatusFailed        = "failed"
	DCLStatusNotRequired   = "not_required"
)

func NewInventoryService() *InventoryService {
	return &InventoryService{
		dpDao:             microWorkDao.NewStoreDepartment(),
		postInventory:     inventory.Post,
		postInventoryOnce: inventory.PostOnce,
	}
}

func newListResult() inventoryResp.ListResult {
	return inventoryResp.ListResult{
		List:            make([]inventoryResp.SkuInfo, 0),
		InventoryStatus: InventoryStatusSuccess,
		FailedUPCs:      make([]string, 0),
		DCLStatus:       DCLStatusNotRequired,
	}
}

func dclSegmentForStoreType(storeType int32) string {
	switch storeType {
	case microWorkDao.STORE_TYPE_RETAIL:
		return "RTL"
	case microWorkDao.STORE_TYPE_OUTLETS:
		return "RTO"
	default:
		return ""
	}
}

func isDclInventory(storeType int32, item inventory.RespDetail) bool {
	if item.LocCode != dclLocCode {
		return false
	}
	return (storeType == microWorkDao.STORE_TYPE_RETAIL && item.Segment == "RTL") ||
		(storeType == microWorkDao.STORE_TYPE_OUTLETS && item.Segment == "RTO")
}

// FetchConcurrentlyBySkuItem 通过skuItem获取到所有upc编码，然后并发获取商品库存
func (s *InventoryService) FetchConcurrentlyBySkuItem(c ctx.Context, skuItem string, upcList []string, storeCode string, storeType int32) (inventoryResp.ListResult, error) {
	result := newListResult()
	// Step 1: 将skuItem统一转换为大写，保证数据一致性
	skuItem = strings.ToUpper(skuItem)

	// Step 2: 初始化结果map，key为UPC编码，用于去重和快速查找
	var resUpcMap = make(map[string]inventoryResp.SkuInfo, 0)

	// Step 3: 【协程1】异步查询大仓(DCL)库存
	// 大仓库存查询与门店库存查询并行执行，提高效率
	dclCh := make(chan inventory.Response)
	if storeType > 0 {
		result.DCLStatus = DCLStatusSuccess
		go func(skuItem string, ch chan inventory.Response) {
			var tmpResponse inventory.Response
			// 调用库存系统接口，查询CN10大仓的所有UPC库存
			// InvLookType="L"表示按层级查询，SummaryType="S"表示汇总查询
			err := s.postInventory(c, true, inventory.ParamsReq{
				LocCode:         []string{dclLocCode}, // CN10是大仓位置编码
				UpcCode:         upcList,              // 批量查询所有UPC
				InvLookType:     "L",
				SummaryType:     "S",
				ExpectedSegment: dclSegmentForStoreType(storeType),
			}, &tmpResponse)
			if err != nil {
				// 发生错误时，将错误信息封装到Response中返回
				ch <- inventory.Response{Error: err}
				return
			}
			ch <- tmpResponse
		}(skuItem, dclCh)
	}

	// Step 4: 【并发查询】使用信号量控制并发数量(最多50个)，防止对库存系统造成过大压力
	// 定义结果结构体，用于channel传递查询结果
	type onHandResult struct {
		upcCode string             // UPC编码
		resp    inventory.Response // 查询响应数据
		err     error              // 查询错误信息
	}
	// 创建带缓冲的channel，缓冲区大小为UPC数量，避免阻塞
	onHandCh := make(chan onHandResult, len(upcList))
	// WaitGroup用于等待所有goroutine执行完毕
	var wg sync.WaitGroup
	// 信号量channel，用于控制最大并发数(30个)
	const maxConcurrency = 30
	sem := make(chan struct{}, maxConcurrency)

	// 为每个UPC启动一个goroutine进行并发查询(受信号量限制)
	for _, upc := range upcList {
		wg.Add(1) // 增加等待计数
		go func(upcCode string) {
			// 获取信号量，如果已满则阻塞等待
			sem <- struct{}{}
			defer func() {
				// 释放信号量
				<-sem
				wg.Done() // goroutine结束时减少等待计数
			}()
			var resp inventory.Response
			// 调用库存接口查询该UPC在所有门店的库存
			// 注意：库存接口限制一次只能查询一个UPC
			log.Info(c).Msgf("库存sku查询-门店库存查询开始, sku_item:%s, upc:%s", skuItem, upcCode)
			err := s.postInventoryOnce(c, false, inventory.ParamsReq{
				//LocCode: allStoreCode,      // 查询所有门店
				UpcCode: []string{upcCode}, // 单次查询一个UPC
			}, &resp)
			// 将结果发送到channel
			onHandCh <- onHandResult{upcCode: upcCode, resp: resp, err: err}
		}(upc)
	}

	// Step 5: 启动goroutine等待所有查询完成后关闭channel
	go func() {
		wg.Wait()       // 阻塞等待所有goroutine完成
		close(onHandCh) // 关闭channel，通知接收方所有数据已发送完毕
	}()

	// Step 6: 合并所有并发请求的结果
	var onHandResp inventory.Response
	var onHandSuccessCount int
	var onHandFailCount int
	var onHandContractErr error
	var onHandSuccessUPCs = make(map[string]struct{}, len(upcList))
	for onHandItem := range onHandCh {
		// 跳过查询失败的UPC，继续处理其他结果
		if onHandItem.err != nil {
			onHandFailCount++
			if inventory.IsContractError(onHandItem.err) {
				onHandContractErr = onHandItem.err
				log.Error(c).Msgf("库存sku查询-门店库存数据契约异常, sku_item:%s, upc:%s, err:%v", skuItem, onHandItem.upcCode, onHandItem.err)
				continue
			}
			result.FailedUPCs = append(result.FailedUPCs, onHandItem.upcCode)
			log.Info(c).Msgf("库存sku查询-门店库存查询失败, sku_item:%s, upc:%s, err:%v", skuItem, onHandItem.upcCode, onHandItem.err)
			continue
		}
		onHandSuccessCount++
		onHandSuccessUPCs[onHandItem.upcCode] = struct{}{}
		// 将该UPC的库存明细追加到总结果中
		onHandResp.DetailList = append(onHandResp.DetailList, onHandItem.resp.DetailList...)
	}
	log.Info(c).Msgf("库存sku查询-门店库存查询汇总, sku_item:%s, upc_count:%d, success_count:%d, fail_count:%d, detail_count:%d", skuItem, len(upcList), onHandSuccessCount, onHandFailCount, len(onHandResp.DetailList))
	if len(result.FailedUPCs) > 0 {
		sort.Strings(result.FailedUPCs)
		result.InventoryStatus = InventoryStatusPartial
	}

	// Step 7: 接收大仓库存查询结果
	var dclResp inventory.Response
	if storeType > 0 {
		// 阻塞读取协程的请求结果，等待大仓查询完成
		dclResp = <-dclCh
	}
	if onHandContractErr != nil {
		return newListResult(), onHandContractErr
	}
	// 大仓查询出错时返回空结果
	if dclResp.Error != nil {
		result.DCLStatus = DCLStatusFailed
	}

	// Step 8: 将大仓库存数据转换为map，便于按UPC快速查找
	// 同时根据门店类型筛选有效的大仓库存数据
	var dclMap = make(map[string]inventory.RespDetail)
	for _, val := range dclResp.DetailList {
		// 按门店类型筛选：零售店(RTL)或奥莱店(RTO)
		if isDclInventory(storeType, val) {
			dclMap[val.UpcCode] = val
		}
	}

	// Step 9: 遍历门店库存数据，构建最终返回结果
	for _, item := range onHandResp.DetailList {
		if item.LocCode == storeCode {
			// 当前门店数据：构建完整的SkuInfo信息
			tmp := inventoryResp.SkuInfo{
				SkuItem: skuItem,
				// 构建商品图片URL，格式：Coach/{skuItem}_{color}_a0
				ImageUrl: "http://s7d2.scene7.com/is/image/Coach/" + strings.ToLower(skuItem+"_"+strings.Replace(item.SkuAttr1, "/", "", 1)) + "_a0",
				Color:    item.SkuAttr1, // 颜色属性
				Style:    item.SkuAttr2, // 款式属性
				UpcInfo: inventoryResp.UpcInfo{
					UpcCode:   item.UpcCode,   // UPC编码
					InTransit: item.InTransit, // 在途库存
					LocCode:   item.LocCode,   // 门店编码
					OnHand:    item.OnHand,    // 可用库存
				},
			}
			// 关联大仓库存：如果该UPC在大仓有库存，则填充DclOnHand字段
			if val, ok := dclMap[item.UpcCode]; ok {
				tmp.UpcInfo.DclOnHand = val.OnHand
			}
			// 存入结果map
			resUpcMap[item.UpcCode] = tmp
		} else {
			// 其他门店数据：仅记录UPC编码(避免覆盖当前门店的完整数据)
			if _, ok := resUpcMap[item.UpcCode]; !ok {
				tmp := inventoryResp.SkuInfo{
					SkuItem:  skuItem,
					ImageUrl: "http://s7d2.scene7.com/is/image/Coach/" + strings.ToLower(skuItem+"_"+strings.Replace(item.SkuAttr1, "/", "", 1)) + "_a0",
					Color:    item.SkuAttr1,
					Style:    item.SkuAttr2,
					UpcInfo: inventoryResp.UpcInfo{
						UpcCode: item.UpcCode,
					},
				}
				if val, exists := dclMap[item.UpcCode]; exists {
					tmp.UpcInfo.DclOnHand = val.OnHand
				}
				resUpcMap[item.UpcCode] = tmp
			}
		}
	}

	// HHT 正常返回空时，只要大仓有库存，仍应展示该 UPC 的大仓库存。
	// HHT 请求失败的 UPC 不在 onHandSuccessUPCs 中，避免将未知的门店库存误判为 0。
	appendDclOnlySkuRows(resUpcMap, dclMap, onHandSuccessUPCs, skuItem)

	// Step 10: 将结果map转换为切片返回
	values := make([]inventoryResp.SkuInfo, 0, len(resUpcMap))
	for _, v := range resUpcMap {
		values = append(values, v)
	}
	sort.SliceStable(values, func(i, j int) bool {
		if values[i].UpcCode != values[j].UpcCode {
			return values[i].UpcCode < values[j].UpcCode
		}
		return false
	})
	result.List = values
	return result, nil
}

func appendDclOnlySkuRows(resUpcMap map[string]inventoryResp.SkuInfo, dclMap map[string]inventory.RespDetail, onHandSuccessUPCs map[string]struct{}, skuItem string) {
	for upcCode, item := range dclMap {
		if item.OnHand <= 0 {
			continue
		}
		if _, ok := onHandSuccessUPCs[upcCode]; !ok {
			continue
		}
		if _, ok := resUpcMap[upcCode]; ok {
			continue
		}
		resUpcMap[upcCode] = inventoryResp.SkuInfo{
			SkuItem:  skuItem,
			ImageUrl: "http://s7d2.scene7.com/is/image/Coach/" + strings.ToLower(skuItem+"_"+strings.Replace(item.SkuAttr1, "/", "", 1)) + "_a0",
			Color:    item.SkuAttr1,
			Style:    item.SkuAttr2,
			UpcInfo: inventoryResp.UpcInfo{
				UpcCode:   upcCode,
				DclOnHand: item.OnHand,
			},
		}
	}
}

// GetListByUpcCode 通过upcCode获取商品库存
func (s *InventoryService) GetListByUpcCode(c ctx.Context, upcCode string, storeCode []string, storeType int32) (inventoryResp.ListResult, error) {
	result := newListResult()
	var upcList = make([]inventoryResp.SkuInfo, 0)
	//由于需要调用两次库存接口，为提高效率使用协程处理首次请求
	dclCh := make(chan inventory.Response)
	if storeType > 0 {
		result.DCLStatus = DCLStatusSuccess
		go func(skuItem string, ch chan inventory.Response) {
			var tmpResponse inventory.Response
			err := s.postInventory(c, true, inventory.ParamsReq{
				LocCode:         []string{dclLocCode},
				UpcCode:         []string{upcCode},
				InvLookType:     "L",
				SummaryType:     "S",
				ExpectedSegment: dclSegmentForStoreType(storeType),
			}, &tmpResponse)
			if err != nil {
				ch <- inventory.Response{Error: err}
				return
			}
			ch <- tmpResponse
		}(upcCode, dclCh)
	}

	var onHandResp inventory.Response
	err := s.postInventory(c, false, inventory.ParamsReq{
		LocCode: storeCode,
		UpcCode: []string{upcCode},
	}, &onHandResp)
	if err != nil {
		return result, err
	}

	var dclResp inventory.Response
	if storeType > 0 {
		dclResp = <-dclCh
	}
	if dclResp.Error != nil {
		result.DCLStatus = DCLStatusFailed
	}

	if len(onHandResp.DetailList) > 0 {
		var dclMap = make(map[string]inventory.RespDetail)
		for _, val := range dclResp.DetailList {
			if isDclInventory(storeType, val) {
				dclMap[val.UpcCode] = val
			}
		}
		for _, item := range onHandResp.DetailList {
			tmp := inventoryResp.UpcInfo{
				UpcCode:   item.UpcCode,
				InTransit: item.InTransit,
				LocCode:   item.LocCode,
				OnHand:    item.OnHand,
			}
			if val, ok := dclMap[item.UpcCode]; ok {
				tmp.DclOnHand = val.OnHand
			}
			upcList = append(upcList, inventoryResp.SkuInfo{
				UpcInfo: tmp,
			})
		}
	} else if len(dclResp.DetailList) > 0 {
		var onHandMap = make(map[string]inventory.RespDetail)
		for _, val := range onHandResp.DetailList {
			onHandMap[val.UpcCode] = val
		}
		for _, item := range dclResp.DetailList {
			if isDclInventory(storeType, item) {
				tmp := inventoryResp.UpcInfo{
					DclOnHand: item.OnHand,
					UpcCode:   item.UpcCode,
				}
				if val, ok := onHandMap[item.UpcCode]; ok {
					tmp.InTransit = val.InTransit
					tmp.LocCode = val.LocCode
					tmp.OnHand = val.OnHand
				}
				upcList = append(upcList, inventoryResp.SkuInfo{
					UpcInfo: tmp,
				})
			}
		}
	}

	result.List = upcList
	return result, nil
}

// GetDcListByUpCode 通过skuItem获取大仓商品库存
func (s *InventoryService) GetDcListByUpCode(c ctx.Context, upcCode string, enterpriseId, memberId int64) (inventoryResp.SkuInfo, error) {
	var upcList inventoryResp.SkuInfo
	memberInfo, err := membersDao.NewMembers().GetMemberInfo(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"id":            memberId,
	})
	if err != nil {
		return upcList, err
	}
	if len(memberInfo.MemberDepartments) == 0 {
		return upcList, nil
	}
	var thirdId int
	for _, department := range memberInfo.MemberDepartments {
		thirdId = department.Department.ThirdID
	}
	storeInfo, err := microWorkDao.NewStoreDepartment().GetInfo(c, map[string]interface{}{
		"enterprise_id":        enterpriseId,
		"wechat_department_id": thirdId,
	})
	if err != nil {
		return upcList, err
	}
	storeType := storeInfo.StoreType
	var dclResp inventory.Response
	req := inventory.ParamsReq{
		LocCode:     []string{dclLocCode},
		UpcCode:     []string{upcCode},
		InvLookType: "L",
		SummaryType: "S",
	}

	err = inventory.Post(c, true, req, &dclResp)
	if err != nil {
		return upcList, err
	}
	for _, item := range dclResp.DetailList {
		if isDclInventory(storeType, item) {
			upcList = inventoryResp.SkuInfo{
				UpcInfo: inventoryResp.UpcInfo{
					UpcCode:   item.UpcCode,
					DclOnHand: item.OnHand,
				},
			}
		}
	}
	return upcList, nil
}

func (s *InventoryService) InsertKeywords(c ctx.Context, keyword, KeywordType string, enterpriseId, memberId int64) error {
	info, err := microScrmDao.NewInventorySearchDao().GetOne(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"member_id":     memberId,
		"keywords":      keyword,
		"keyword_type":  KeywordType,
	})
	if err != nil {
		return err
	}
	nowTime := time.Now().Unix()
	if info.ID > 0 {
		err = microScrmDao.NewInventorySearchDao().UpdateData(c, map[string]interface{}{
			"id": info.ID,
		}, map[string]interface{}{
			"updated_at": nowTime,
		})
		return err
	}
	err = microScrmDao.NewInventorySearchDao().CreateData(c, microScrmModel.InventorySearchModel{
		MicroBaseModel: models.MicroBaseModel{
			BaseModel: models.BaseModel{
				ID:        app.UID(),
				CreatedAt: nowTime,
				UpdatedAt: nowTime,
				DeletedAt: 0,
			},
			EnterpriseID: enterpriseId,
		},
		MemberID:    memberId,
		KeywordType: KeywordType,
		Keywords:    keyword,
	})
	return err
}

func (s *InventoryService) GetSearchKeywords(c ctx.Context, enterpriseId, memberId int64) ([]microScrmModel.InventorySearchModel, error) {
	return microScrmDao.NewInventorySearchDao().GetKeywords(c, enterpriseId, memberId)

}

func (s *InventoryService) DelKeywords(c ctx.Context, enterpriseId, memberId int64, ids []int64, isAll int) error {
	whereMap := map[string]interface{}{
		"enterprise_id": enterpriseId,
		"member_id":     memberId,
	}
	switch isAll {
	case microScrmModel.InventorySearchTypeSkuItem:
		whereMap["keyword_type"] = "sku_item"
	case microScrmModel.InventorySearchTypeUpcCode:
		whereMap["keyword_type"] = "upc_code"
	default:
		whereMap["id"] = ids
	}
	return microScrmDao.NewInventorySearchDao().DeleteData(c, whereMap)
}

func (s *InventoryService) CheckKeywordsNum(c ctx.Context, enterpriseId, memberId int64, keywordType string) error {
	limitInfo, err := microScrmDao.NewInventorySearchDao().GetLimitInfo(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"member_id":     memberId,
		"keyword_type":  keywordType,
	})
	if err != nil || limitInfo.ID == 0 {
		return err
	}
	if limitInfo.ID > 0 {
		return microScrmDao.NewInventorySearchDao().DeleteData(c, "enterprise_id=? and member_id=? and keyword_type=? and updated_at <=?", enterpriseId, memberId, keywordType, limitInfo.UpdatedAt)
	}
	return nil
}
