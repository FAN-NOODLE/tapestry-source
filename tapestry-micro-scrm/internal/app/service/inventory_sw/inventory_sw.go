package inventorySwSvc

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service"
	baiqiuSdk "micro-scrm/internal/pkg/sdk/baiqiu"
	"micro-scrm/internal/pkg/sdk/inventory_sw"
	"time"
)

type InventorySwService struct {
	service.Base
	dpDao *microWorkDao.StoreDepartment
}

func NewInventorySwService() *InventorySwService {
	return &InventorySwService{
		dpDao: microWorkDao.NewStoreDepartment(),
	}
}

// GetList 获取商品库存
func (s *InventorySwService) GetList(c ctx.Context, params inventorySw.ParamsReq) ([]inventorySw.RespDetail, error) {
	response, err := inventorySw.Post(c, params)
	if err != nil || len(response.Data) == 0 {
		return []inventorySw.RespDetail{}, err
	}
	var returnSlice []inventorySw.RespDetail
	// 使用map去重，只保留不同的商品组合
	productMap := make(map[baiqiuSdk.ProductQuery]bool)
	for _, val := range response.Data {
		if val.StoreCode == "R098" {
			//过滤掉JD大仓
			continue
		}
		product := baiqiuSdk.ProductQuery{
			ProductCode: val.StkCode,
			Color:       val.Color,
		}
		productMap[product] = true
		returnSlice = append(returnSlice, val)
	}

	// 将去重后的商品转换为slice
	var uniqueProducts []baiqiuSdk.ProductQuery
	for product := range productMap {
		uniqueProducts = append(uniqueProducts, product)
	}

	// 分批查询（最多只需要查询一次，如果去重后数量<=20）
	batchSize := 20
	var allImageList []baiqiuSdk.Product

	for i := 0; i < len(uniqueProducts); i += batchSize {
		end := i + batchSize
		if end > len(uniqueProducts) {
			end = len(uniqueProducts)
		}

		batch := uniqueProducts[i:end]
		imageList, _ := baiqiuSdk.NewBatchProductClient().BatchSearch(c, baiqiuSdk.BatchSearchRequest{
			Products: batch,
		})
		if imageList.Data != nil && len(imageList.Data.Products) > 0 {
			allImageList = append(allImageList, imageList.Data.Products...)
		}
	}

	// 使用查询到的图片数据构建映射
	if len(allImageList) > 0 {
		var imageMap = make(map[string]string)
		for _, image := range allImageList {
			imageMap[image.ProductCode+"*"+image.Color] = image.ImagePath
		}
		for i, val := range returnSlice {
			if imagePath, ok := imageMap[val.StkCode+"*"+val.Color]; ok {
				returnSlice[i].PicUrl = imagePath
			}
		}
	}

	return returnSlice, nil
}

func (s *InventorySwService) InsertKeywords(c ctx.Context, keyword, KeywordType string, enterpriseId, memberId int64) error {
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

func (s *InventorySwService) GetSearchKeywords(c ctx.Context, enterpriseId, memberId int64) ([]microScrmModel.InventorySearchModel, error) {
	return microScrmDao.NewInventorySearchDao().GetKeywords(c, enterpriseId, memberId)

}

func (s *InventorySwService) DelKeywords(c ctx.Context, enterpriseId, memberId int64, ids []int64, delType string) error {
	whereMap := map[string]interface{}{
		"enterprise_id": enterpriseId,
		"member_id":     memberId,
	}
	if delType != "" {
		whereMap["keyword_type"] = delType
	} else {
		whereMap["id"] = ids
	}
	return microScrmDao.NewInventorySearchDao().DeleteData(c, whereMap)
}

func (s *InventorySwService) CheckKeywordsNum(c ctx.Context, enterpriseId, memberId int64, keywordType string) error {
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
