package source

import (
	"encoding/json"
	"fmt"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	BaozunCacheDao "micro-scrm/internal/app/dao/source"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	sourceResp "micro-scrm/internal/app/http/response/source"
	"micro-scrm/internal/app/models"
	sourceModel "micro-scrm/internal/app/models/source"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/redis"
	"micro-scrm/internal/pkg/sdk/baozun"
	"micro-scrm/internal/pkg/sdk/talkie"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

type BrandService struct {
	BaozunCacheDao *BaozunCacheDao.BaozunCacheDao
}

func NewBrandService() *BrandService {
	return &BrandService{
		BaozunCacheDao: BaozunCacheDao.NewBaozunCacheDao(),
	}
}

func (s *BrandService) SearchList(c app.Ctx, req sourceRequest.BrandListReq) ([]sourceResp.BrandProductInfo, error) {
	var data baozun.DataResponse
	var respProductList = make([]sourceResp.BrandProductInfo, 0)
	if req.StoreCode == "" && req.Mobile != "" {
		memberInfo, _ := membersDao.NewMembers().GetMemberInfoAndDepartmentsByMobile(c, req.EnterpriseId, req.Mobile)
		if memberInfo.DepartmentThirdId > 0 {
			storeList, _ := microWorkDao.NewStoreDepartment().GetAllByWechatDepartmentIds(c, req.EnterpriseId, []int{int(memberInfo.DepartmentThirdId)})
			if len(storeList) > 0 {
				req.StoreCode = storeList[0].ThirdStoreCode
			}
		}
	}
	err := baozun.SearchResources(c, baozun.SearchParams{
		EnterpriseId:            req.EnterpriseId,
		Upc:                     req.UpcCode,
		UpcSlice:                req.UpcCodes,
		ProductName:             req.ProductName,
		ProductNo:               req.SkuItem,
		SkuColor:                req.SkuColor,
		IsLike:                  !req.Accurate,
		StoreCode:               req.StoreCode,
		Size:                    req.Size,
		LastCode:                req.LastCode,
		IncludeNoneImageProduct: req.IncludeNoneImageProduct,
	}, &data)
	if len(data.Data) == 0 {
		return respProductList, nil
	}
	var resourceKeys []string
	for _, product := range data.Data {
		for resourceKey := range product.ExtData.Resource {
			resourceKeys = append(resourceKeys, resourceKey)
		}
	}
	var cacheMap = make(map[string]sourceModel.BaozunCacheModel)
	var existsKeys []string
	var cacheIds []int64
	if len(resourceKeys) > 0 {
		cacheList, err := BaozunCacheDao.NewBaozunCacheDao().GetMany(c, req.EnterpriseId, "resource_key in ?", resourceKeys)
		if err != nil {
			return respProductList, err
		}
		for _, cache := range cacheList {
			cacheMap[cache.ResourceKey] = cache
			existsKeys = append(existsKeys, cache.ResourceKey)
			cacheIds = append(cacheIds, cache.ID)
		}
	}
	// 匹配文件名中的数字部分的正则表达式
	re := regexp.MustCompile(`([0-9]+)\.jpg$`)
	for _, product := range data.Data {
		var sourceList = make([]sourceResp.BrandSourceInfo, 0)
		for resourceKey, imageInfo := range product.ExtData.Resource {
			parsedURL, err := url.Parse(imageInfo.Url)
			if err != nil {
				return respProductList, err
			}

			tmpSourceInfo := sourceResp.BrandSourceInfo{
				ResourceKey:  resourceKey,
				EnterpriseID: req.EnterpriseId,
				Name:         imageInfo.Name,
				Type:         helper.RuleThree(imageInfo.Type == 1, sourceModel.TypeImg, sourceModel.TypeVideo),
				URL:          fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, parsedURL.Path),
				Size:         int64(imageInfo.Meta.Size),
				Ext:          imageInfo.Meta.Type,
				IsMain:       helper.RuleThree(strings.Contains(imageInfo.Name, "_a0"), sourceModel.BrandIsMainImage, 0),
				Sort:         1000, //默认排序取1000，这样匹配不到的图片都可以拍到后面去
			}
			matches := re.FindStringSubmatch(imageInfo.Name)
			if len(matches) > 1 {
				tmpSourceInfo.Sort = helper.StringToInt(matches[1]) // 第一个捕获组中的匹配结果
			}
			if val, ok := cacheMap[resourceKey]; ok {
				tmpSourceInfo.MediaID = val.MediaId
				tmpSourceInfo.MediaExpireTime = val.MediaExpireTime
				tmpSourceInfo.MediaStatus = val.MediaStatus
				tmpSourceInfo.DownloadTimes = val.DownloadTimes
				tmpSourceInfo.ForwardingTimes = val.ForwardingTimes

			}
			sourceList = append(sourceList, tmpSourceInfo)
		}
		var upcCode string
		var skuTagPrice string
		var upcList []sourceResp.UpcListResp
		if len(product.Variants) > 0 {
			skuTagPrice = product.Variants[0].Properties.SkuTagPrice
			if len(product.Variants[0].Fragments) > 0 {
				upcCode = product.Variants[0].Fragments[0].Properties.Upc
			}
			for _, variant := range product.Variants {
				if len(variant.Fragments[0].Properties.Upc) > 0 {
					tmpUpc := sourceResp.UpcListResp{
						UpcCode:     variant.Fragments[0].Properties.Upc,
						VariantCode: variant.VariantCode,
						SkuNumber:   variant.Fragments[0].Properties.SkuNumber,
					}
					upcList = append(upcList, tmpUpc)
				}

			}
		}
		var storePrice float64
		if len(product.ExtData.Prices) > 0 && len(product.ExtData.Prices[0].Prices) > 0 {
			storePrice = product.ExtData.Prices[0].Prices[0].Value
		}
		var color string
		var styleCode string
		var department string
		var silhouetteDesc string
		var shortDescriptionZH string
		var longDescriptionZH string
		var productNameEn string
		if len(product.Master.Fragments) > 0 {
			tmProperties := product.Master.Fragments[0].Properties
			color = tmProperties.Color
			styleCode = tmProperties.StyleCode
			department = tmProperties.Department
			silhouetteDesc = tmProperties.SilhouetteDesc
			shortDescriptionZH = tmProperties.ShortDescriptionZhCN
			longDescriptionZH = tmProperties.LongDescriptionZhCN
			productNameEn = tmProperties.ProductNameEn
		}
		sort.Slice(sourceList, func(i, j int) bool {
			return sourceList[i].Sort < sourceList[j].Sort
		})
		tmpProductInfo := sourceResp.BrandProductInfo{
			ID:                 product.Master.ID,
			EnterpriseID:       req.EnterpriseId,
			ProductName:        helper.RuleThree(product.Master.Properties.ProductName != "", product.Master.Properties.ProductName, productNameEn),
			SkuItem:            product.Master.Properties.ProductNo,
			UpcCode:            upcCode,
			UpcList:            upcList,
			Color:              color,
			StyleCode:          styleCode,
			Department:         department,
			SkuTagPrice:        skuTagPrice,
			StorePrice:         storePrice,
			StoreCode:          req.StoreCode,
			SilhouetteDesc:     silhouetteDesc,
			ShortDescriptionZH: shortDescriptionZH,
			LongDescriptionZH:  longDescriptionZH,
			ProductCode:        product.Master.ProductCode,
			BrandSourceList:    sourceList,
		}
		respProductList = append(respProductList, tmpProductInfo)
	}
	err = s.SaveCache(c, req.EnterpriseId, respProductList, existsKeys, []string{})
	if err != nil {
		log.Info(c).Msgf("素材库::品牌素材缓存失败:%v", err)
	}
	return respProductList, nil

}

func (s *BrandService) GetUpcListBySkuItem(c app.Ctx, enterpriseId int64, skuItem string) ([]string, error) {
	var upcList = make([]string, 0)
	var data baozun.DataResponse
	_ = baozun.SearchResources(c, baozun.SearchParams{
		EnterpriseId:            enterpriseId,
		SkuItem:                 skuItem,
		Size:                    100,
		IncludeNoneImageProduct: true,
	}, &data)
	if len(data.Data) == 0 {
		return []string{}, nil
	}
	for _, product := range data.Data {
		if len(product.Variants) > 0 {
			for _, variant := range product.Variants {
				if len(variant.Fragments[0].Properties.Upc) > 0 {
					upcList = append(upcList, variant.Fragments[0].Properties.Upc)
				}
			}
		}
	}
	return helper.ArrayUniq(upcList), nil
}

func (s *BrandService) SaveCache(c app.Ctx, enterpriseId int64, list []sourceResp.BrandProductInfo, existsKeys, existsProductCodes []string) error {
	var insertList []sourceModel.BaozunCacheModel
	var insertProductList []sourceModel.BaozunCollectModel
	for _, product := range list {
		for _, val := range product.BrandSourceList {
			if !helper.InArray(val.ResourceKey, existsKeys) {
				existsKeys = append(existsKeys, val.ResourceKey)
				tmpInsertData := sourceModel.BaozunCacheModel{
					MicroBaseModel: models.MicroBaseModel{
						BaseModel:    models.CreateBase(app.UID()),
						EnterpriseID: enterpriseId,
					},
					IsMain:      helper.RuleThree(strings.Contains(val.Name, "_a0"), sourceModel.BrandIsMainImage, 0),
					Type:        val.Type,
					Name:        val.Name,
					URL:         val.URL,
					ProductName: product.ProductName,
					SkuItem:     product.SkuItem,
					UpcCode:     product.UpcCode,
					Color:       product.Color,
					ResourceKey: val.ResourceKey,
					Size:        val.Size,
					Ext:         val.Ext,
				}
				insertList = append(insertList, tmpInsertData)
			}
			if len(existsProductCodes) > 0 && !helper.InArray(product.ProductCode, existsProductCodes) {
				existsProductCodes = append(existsProductCodes, product.ProductCode)
				tempInsertProductData := sourceModel.BaozunCollectModel{
					MicroBaseModel: models.MicroBaseModel{
						BaseModel:    models.CreateBase(app.UID()),
						EnterpriseID: enterpriseId,
					},
					ProductCode: product.ProductCode,
					CollectNum:  0,
				}
				insertProductList = append(insertProductList, tempInsertProductData)
			}
		}
	}
	if len(insertList) > 0 {
		err := s.BaozunCacheDao.CreateInBatchesData(c, insertList)
		if err != nil {
			return err
		}
	}
	if len(insertProductList) > 0 {
		err := BaozunCacheDao.NewBaozunCollectDao().CreateInBatchesData(c, insertProductList)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BrandService) CacheInfo(c app.Ctx, req sourceRequest.BrandCacheReq) (data sourceResp.BrandCacheResp, err error) {
	productList, err := s.SearchList(c, sourceRequest.BrandListReq{
		BrandCacheReq: req,
		Size:          1,
	})
	if err != nil {
		return data, err
	}

	return buildBrandCacheResp(productList), nil
}

func buildBrandCacheResp(productList []sourceResp.BrandProductInfo) sourceResp.BrandCacheResp {
	if len(productList) == 0 || len(productList[0].BrandSourceList) == 0 {
		return sourceResp.BrandCacheResp{}
	}

	product := productList[0]
	selected := product.BrandSourceList[0]
	for _, source := range product.BrandSourceList {
		if source.IsMain == sourceModel.BrandIsMainImage {
			selected = source
			break
		}
	}

	return sourceResp.BrandCacheResp{
		IsMain:      selected.IsMain,
		Type:        selected.Type,
		Name:        selected.Name,
		URL:         selected.URL,
		ProductName: product.ProductName,
		SkuItem:     product.SkuItem,
		UpcCode:     product.UpcCode,
		Color:       product.Color,
		ResourceKey: selected.ResourceKey,
		Size:        selected.Size,
		Ext:         selected.Ext,
		SourceNum:   len(product.BrandSourceList),
	}
}

func (s *BrandService) CollectSearchList(c app.Ctx, req sourceRequest.BrandListReq) ([]sourceResp.BrandProductInfo, error) {
	var data baozun.DataResponse
	var respProductList = make([]sourceResp.BrandProductInfo, 0)
	memberInfo, err := membersDao.NewMembers().GetMemberInfoAndDepartmentsByMobile(c, req.EnterpriseId, req.Mobile)
	if err != nil {
		return nil, err
	}
	collectList, err := BaozunCacheDao.NewCollectRelationDao().GetMany(c, map[string]interface{}{
		"enterprise_id": memberInfo.EnterpriseId,
		"member_id":     memberInfo.MembersId,
		"type":          sourceModel.TypeProduct,
		"is_collect":    sourceModel.SourceCollectYes,
	})
	if err != nil {
		return nil, err
	}
	collectProductCodes := make([]string, 0)
	if len(collectList) > 0 {
		for _, val := range collectList {
			collectProductCodes = append(collectProductCodes, val.ProductCode)
		}
	}
	if req.StoreCode == "" {
		if memberInfo.DepartmentThirdId > 0 {
			storeList, _ := microWorkDao.NewStoreDepartment().GetAllByWechatDepartmentIds(c, req.EnterpriseId, []int{int(memberInfo.DepartmentThirdId)})
			if len(storeList) > 0 {
				req.StoreCode = storeList[0].ThirdStoreCode
			}
		}
	}
	err = baozun.SearchResources(c, baozun.SearchParams{
		EnterpriseId: req.EnterpriseId,
		Upc:          req.UpcCode,
		UpcSlice:     req.UpcCodes,
		ProductName:  req.ProductName,
		ProductNo:    req.SkuItem,
		IsLike:       !req.Accurate,
		StoreCode:    req.StoreCode,
		ProductCode:  req.ProductCode,
		Size:         req.Size,
		LastCode:     req.LastCode,
	}, &data)
	if len(data.Data) == 0 {
		return respProductList, nil
	}
	var resourceKeys []string
	var productCodes []string
	for _, product := range data.Data {
		for resourceKey := range product.ExtData.Resource {
			resourceKeys = append(resourceKeys, resourceKey)
		}
		productCodes = append(productCodes, product.Master.ProductCode)
	}
	var cacheMap = make(map[string]sourceModel.BaozunCacheModel)
	var existsKeys []string
	var existsProductCodes []string
	var cacheIds []int64
	if len(resourceKeys) > 0 {
		cacheList, err := BaozunCacheDao.NewBaozunCacheDao().GetMany(c, req.EnterpriseId, "resource_key in ?", resourceKeys)
		if err != nil {
			return respProductList, err
		}
		for _, cache := range cacheList {
			cacheMap[cache.ResourceKey] = cache
			existsKeys = append(existsKeys, cache.ResourceKey)
			cacheIds = append(cacheIds, cache.ID)
		}
	}
	collectProductMap := make(map[string]int)
	if len(productCodes) > 0 {
		list, err := BaozunCacheDao.NewBaozunCollectDao().GetMany(c, map[string]interface{}{
			"enterprise_id": memberInfo.EnterpriseId,
			"product_code":  productCodes,
		})
		if err != nil {
			return respProductList, err
		}
		if len(list) > 0 {
			for _, val := range list {
				existsProductCodes = append(existsProductCodes, val.ProductCode)
				collectProductMap[val.ProductCode] = val.CollectNum
			}
		}
	}
	// 匹配文件名中的数字部分的正则表达式
	re := regexp.MustCompile(`([0-9]+)\.jpg$`)
	for _, product := range data.Data {
		var sourceList = make([]sourceResp.BrandSourceInfo, 0)
		for resourceKey, imageInfo := range product.ExtData.Resource {
			parsedURL, err := url.Parse(imageInfo.Url)
			if err != nil {
				return respProductList, err
			}

			tmpSourceInfo := sourceResp.BrandSourceInfo{
				ResourceKey:  resourceKey,
				EnterpriseID: req.EnterpriseId,
				Name:         imageInfo.Name,
				Type:         helper.RuleThree(imageInfo.Type == 1, sourceModel.TypeImg, sourceModel.TypeVideo),
				URL:          fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, parsedURL.Path),
				Size:         int64(imageInfo.Meta.Size),
				Ext:          imageInfo.Meta.Type,
				IsMain:       helper.RuleThree(strings.Contains(imageInfo.Name, "_a0"), sourceModel.BrandIsMainImage, 0),
				Sort:         1000, //默认排序取1000，这样匹配不到的图片都可以拍到后面去
			}
			matches := re.FindStringSubmatch(imageInfo.Name)
			if len(matches) > 1 {
				tmpSourceInfo.Sort = helper.StringToInt(matches[1]) // 第一个捕获组中的匹配结果
			}
			if val, ok := cacheMap[resourceKey]; ok {
				tmpSourceInfo.MediaID = val.MediaId
				tmpSourceInfo.MediaExpireTime = val.MediaExpireTime
				tmpSourceInfo.MediaStatus = val.MediaStatus
				tmpSourceInfo.DownloadTimes = val.DownloadTimes
				tmpSourceInfo.ForwardingTimes = val.ForwardingTimes

			}
			sourceList = append(sourceList, tmpSourceInfo)
		}
		var upcCode string
		var skuTagPrice string
		var upcList []sourceResp.UpcListResp
		if len(product.Variants) > 0 {
			skuTagPrice = product.Variants[0].Properties.SkuTagPrice
			if len(product.Variants[0].Fragments) > 0 {
				upcCode = product.Variants[0].Fragments[0].Properties.Upc
			}
			for _, variant := range product.Variants {
				if len(variant.Fragments[0].Properties.Upc) > 0 {
					tmpUpc := sourceResp.UpcListResp{
						UpcCode:     variant.Fragments[0].Properties.Upc,
						VariantCode: variant.VariantCode,
						SkuNumber:   variant.Fragments[0].Properties.SkuNumber,
					}
					upcList = append(upcList, tmpUpc)
				}

			}
		}
		var storePrice float64
		if len(product.ExtData.Prices) > 0 && len(product.ExtData.Prices[0].Prices) > 0 {
			storePrice = product.ExtData.Prices[0].Prices[0].Value
		}
		var color string
		var styleCode string
		var department string
		var silhouetteDesc string
		var shortDescriptionZH string
		var longDescriptionZH string
		var productNameEn string
		if len(product.Master.Fragments) > 0 {
			tmProperties := product.Master.Fragments[0].Properties
			color = tmProperties.Color
			styleCode = tmProperties.StyleCode
			department = tmProperties.Department
			silhouetteDesc = tmProperties.SilhouetteDesc
			shortDescriptionZH = tmProperties.ShortDescriptionZhCN
			longDescriptionZH = tmProperties.LongDescriptionZhCN
			productNameEn = tmProperties.ProductNameEn
		}
		sort.Slice(sourceList, func(i, j int) bool {
			return sourceList[i].Sort < sourceList[j].Sort
		})
		isCollect := 0
		if helper.InArray(product.Master.ProductCode, collectProductCodes) {
			isCollect = 1
		}
		collectNum := 0
		if val, ok := collectProductMap[product.Master.ProductCode]; ok {
			collectNum = val
		}
		tmpProductInfo := sourceResp.BrandProductInfo{
			ID:                 product.Master.ID,
			EnterpriseID:       req.EnterpriseId,
			ProductName:        helper.RuleThree(product.Master.Properties.ProductName != "", product.Master.Properties.ProductName, productNameEn),
			SkuItem:            product.Master.Properties.ProductNo,
			UpcCode:            upcCode,
			UpcList:            upcList,
			Color:              color,
			StyleCode:          styleCode,
			Department:         department,
			SkuTagPrice:        skuTagPrice,
			StorePrice:         storePrice,
			StoreCode:          req.StoreCode,
			SilhouetteDesc:     silhouetteDesc,
			ShortDescriptionZH: shortDescriptionZH,
			LongDescriptionZH:  longDescriptionZH,
			ProductCode:        product.Master.ProductCode,
			BrandSourceList:    sourceList,
			IsCollect:          isCollect,
			CollectNum:         collectNum,
		}
		respProductList = append(respProductList, tmpProductInfo)
	}
	err = s.SaveCache(c, req.EnterpriseId, respProductList, existsKeys, existsProductCodes)
	if err != nil {
		log.Info(c).Msgf("素材库::品牌素材缓存失败:%v", err)
	}
	return respProductList, nil
}

func (s *BrandService) CollectList(c app.Ctx, req sourceRequest.BrandListReq, keywords string) ([]sourceResp.BrandProductInfo, error) {
	var data baozun.DataResponse
	var respProductList = make([]sourceResp.BrandProductInfo, 0)
	memberInfo, err := membersDao.NewMembers().GetMemberInfoAndDepartmentsByMobile(c, req.EnterpriseId, req.Mobile)
	if err != nil {
		return nil, err
	}
	collectList, err := BaozunCacheDao.NewCollectRelationDao().GetCollectProduct(c, memberInfo.EnterpriseId, memberInfo.MembersId, keywords)
	if err != nil {
		return nil, err
	}
	collectProductCodes := make([]string, 0)
	collectProductMap := make(map[string]int)
	collectTimeMap := make(map[string]int64, 0)
	if len(collectList) == 0 {
		return respProductList, nil
	}
	if len(collectList) > 0 {
		for _, val := range collectList {
			collectProductCodes = append(collectProductCodes, val.ProductCode)
			collectTimeMap[val.ProductCode] = val.UpdatedAt
		}
		baozunCollect, err := BaozunCacheDao.NewBaozunCollectDao().GetMany(c, map[string]interface{}{
			"enterprise_id": memberInfo.EnterpriseId,
			"product_code":  collectProductCodes,
		})
		if err != nil {
			return nil, err
		}
		for _, val := range baozunCollect {
			collectProductMap[val.ProductCode] = val.CollectNum
		}

	}
	if req.StoreCode == "" {
		if memberInfo.DepartmentThirdId > 0 {
			storeList, _ := microWorkDao.NewStoreDepartment().GetAllByWechatDepartmentIds(c, req.EnterpriseId, []int{int(memberInfo.DepartmentThirdId)})
			if len(storeList) > 0 {
				req.StoreCode = storeList[0].ThirdStoreCode
			}
		}
	}
	err = baozun.SearchResources(c, baozun.SearchParams{
		EnterpriseId: req.EnterpriseId,
		Upc:          req.UpcCode,
		UpcSlice:     req.UpcCodes,
		ProductName:  req.ProductName,
		ProductNo:    req.SkuItem,
		IsLike:       !req.Accurate,
		StoreCode:    req.StoreCode,
		ProductCode:  collectProductCodes,
		Size:         req.Size,
		LastCode:     req.LastCode,
	}, &data)
	if len(data.Data) == 0 {
		return respProductList, nil
	}
	var resourceKeys []string
	for _, product := range data.Data {
		for resourceKey := range product.ExtData.Resource {
			resourceKeys = append(resourceKeys, resourceKey)
		}
	}
	var cacheMap = make(map[string]sourceModel.BaozunCacheModel)

	var cacheIds []int64
	if len(resourceKeys) > 0 {
		cacheList, err := BaozunCacheDao.NewBaozunCacheDao().GetMany(c, req.EnterpriseId, "resource_key in ?", resourceKeys)
		if err != nil {
			return respProductList, err
		}
		for _, cache := range cacheList {
			cacheMap[cache.ResourceKey] = cache

			cacheIds = append(cacheIds, cache.ID)
		}
	}
	// 匹配文件名中的数字部分的正则表达式
	re := regexp.MustCompile(`([0-9]+)\.jpg$`)
	for _, product := range data.Data {
		var sourceList = make([]sourceResp.BrandSourceInfo, 0)
		for resourceKey, imageInfo := range product.ExtData.Resource {
			parsedURL, err := url.Parse(imageInfo.Url)
			if err != nil {
				return respProductList, err
			}

			tmpSourceInfo := sourceResp.BrandSourceInfo{
				ResourceKey:  resourceKey,
				EnterpriseID: req.EnterpriseId,
				Name:         imageInfo.Name,
				Type:         helper.RuleThree(imageInfo.Type == 1, sourceModel.TypeImg, sourceModel.TypeVideo),
				URL:          fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, parsedURL.Path),
				Size:         int64(imageInfo.Meta.Size),
				Ext:          imageInfo.Meta.Type,
				IsMain:       helper.RuleThree(strings.Contains(imageInfo.Name, "_a0"), sourceModel.BrandIsMainImage, 0),
				Sort:         1000, //默认排序取1000，这样匹配不到的图片都可以拍到后面去
			}
			matches := re.FindStringSubmatch(imageInfo.Name)
			if len(matches) > 1 {
				tmpSourceInfo.Sort = helper.StringToInt(matches[1]) // 第一个捕获组中的匹配结果
			}
			if val, ok := cacheMap[resourceKey]; ok {
				tmpSourceInfo.MediaID = val.MediaId
				tmpSourceInfo.MediaExpireTime = val.MediaExpireTime
				tmpSourceInfo.MediaStatus = val.MediaStatus
				tmpSourceInfo.DownloadTimes = val.DownloadTimes
				tmpSourceInfo.ForwardingTimes = val.ForwardingTimes

			}
			sourceList = append(sourceList, tmpSourceInfo)
		}
		var upcCode string
		var skuTagPrice string
		var upcList []sourceResp.UpcListResp
		if len(product.Variants) > 0 {
			skuTagPrice = product.Variants[0].Properties.SkuTagPrice
			if len(product.Variants[0].Fragments) > 0 {
				upcCode = product.Variants[0].Fragments[0].Properties.Upc
			}
			for _, variant := range product.Variants {
				if len(variant.Fragments[0].Properties.Upc) > 0 {
					tmpUpc := sourceResp.UpcListResp{
						UpcCode:     variant.Fragments[0].Properties.Upc,
						VariantCode: variant.VariantCode,
						SkuNumber:   variant.Fragments[0].Properties.SkuNumber,
					}
					upcList = append(upcList, tmpUpc)
				}

			}
		}
		var storePrice float64
		if len(product.ExtData.Prices) > 0 && len(product.ExtData.Prices[0].Prices) > 0 {
			storePrice = product.ExtData.Prices[0].Prices[0].Value
		}
		var color string
		var styleCode string
		var department string
		var silhouetteDesc string
		var shortDescriptionZH string
		var longDescriptionZH string
		var productNameEn string
		if len(product.Master.Fragments) > 0 {
			tmProperties := product.Master.Fragments[0].Properties
			color = tmProperties.Color
			styleCode = tmProperties.StyleCode
			department = tmProperties.Department
			silhouetteDesc = tmProperties.SilhouetteDesc
			shortDescriptionZH = tmProperties.ShortDescriptionZhCN
			longDescriptionZH = tmProperties.LongDescriptionZhCN
			productNameEn = tmProperties.ProductNameEn
		}
		sort.Slice(sourceList, func(i, j int) bool {
			return sourceList[i].Sort < sourceList[j].Sort
		})
		isCollect := 0
		if helper.InArray(product.Master.ProductCode, collectProductCodes) {
			isCollect = 1
		}
		collectNum := 0
		collectAt := int64(0)
		if val, ok := collectProductMap[product.Master.ProductCode]; ok {
			collectNum = val
			collectAt = collectTimeMap[product.Master.ProductCode]
		}
		tmpProductInfo := sourceResp.BrandProductInfo{
			ID:                 product.Master.ID,
			EnterpriseID:       req.EnterpriseId,
			ProductName:        helper.RuleThree(product.Master.Properties.ProductName != "", product.Master.Properties.ProductName, productNameEn),
			SkuItem:            product.Master.Properties.ProductNo,
			UpcCode:            upcCode,
			UpcList:            upcList,
			Color:              color,
			StyleCode:          styleCode,
			Department:         department,
			SkuTagPrice:        skuTagPrice,
			StorePrice:         storePrice,
			StoreCode:          req.StoreCode,
			SilhouetteDesc:     silhouetteDesc,
			ShortDescriptionZH: shortDescriptionZH,
			LongDescriptionZH:  longDescriptionZH,
			ProductCode:        product.Master.ProductCode,
			BrandSourceList:    sourceList,
			IsCollect:          isCollect,
			CollectNum:         collectNum,
			CollectAt:          collectAt,
		}
		if keywords != "" {

		}
		respProductList = append(respProductList, tmpProductInfo)
	}
	sort.Slice(respProductList, func(i, j int) bool {
		return respProductList[i].CollectAt > respProductList[j].CollectAt
	})
	return respProductList, nil
}

// TalkieImageSearch Talkie图片搜索
func (s *BrandService) TalkieImageSearch(c app.Ctx, req sourceRequest.TalkieImageSearchReq) ([]sourceResp.TalkieImageSearchRegion, error) {
	var resp []sourceResp.TalkieImageSearchRegion

	// 调用 Talkie 图片搜索接口
	talkieResp, err := talkie.ImageSearch(c, talkie.ImageSearchRequest{
		ImageURL: req.ImageURL,
		Num:      req.Num,
	})
	if err != nil {
		return resp, err
	}

	// 收集所有 ProductID，用于批量查询宝尊商品信息
	productIDSet := make(map[string]struct{})
	for _, region := range talkieResp.Data {
		for _, pic := range region.PicList {
			if pic.ProductID != "" {
				productIDSet[pic.ProductID] = struct{}{}
			}
		}
	}

	// 先从 Redis 缓存中获取商品信息，收集未命中的 ProductID
	productInfoMap := make(map[string]*baozun.ProductBriefCache)
	var missedIDs []string
	for productID := range productIDSet {
		cacheKey := fmt.Sprintf(baozun.ProductInfoRedisKey, req.EnterpriseId, productID)
		cacheStr := redis.Client().Get(c, cacheKey).Val()
		if cacheStr != "" {
			var cache baozun.ProductBriefCache
			if json.Unmarshal([]byte(cacheStr), &cache) == nil {
				productInfoMap[productID] = &cache
				continue
			}
		}
		missedIDs = append(missedIDs, productID)
	}

	// 未命中缓存的部分，用 ProductCode 批量查询宝尊商品信息
	if len(missedIDs) > 0 {
		var data baozun.DataResponse
		err := baozun.SearchResources(c, baozun.SearchParams{
			EnterpriseId:            req.EnterpriseId,
			ProductCode:             missedIDs,
			IsLike:                  false, // 精准匹配
			IncludeNoneImageProduct: true,
			Size:                    len(missedIDs),
		}, &data)
		if err != nil {
			log.Info(c).Msgf("Talkie图片搜索-宝尊商品信息批量查询失败, err: %v", err)
		}
		if err == nil && len(data.Data) > 0 {
			for _, product := range data.Data {
				productCode := product.Master.ProductCode
				var cache baozun.ProductBriefCache
				cache.ProductNo = product.Master.Properties.ProductNo
				cache.Name = product.Master.Properties.ProductName
				if len(product.Variants) > 0 {
					if len(product.Variants[0].Fragments) > 0 {
						cache.Upc = product.Variants[0].Fragments[0].Properties.Upc
					}
				}
				// 取 a0 主图地址
				for _, imageInfo := range product.ExtData.Resource {
					if strings.Contains(imageInfo.Name, "_a0") {
						parsedURL, parseErr := url.Parse(imageInfo.Url)
						if parseErr == nil {
							cache.Url = fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, parsedURL.Path)
						}
						break
					}
				}
				// 取商品颜色
				if len(product.Master.Fragments) > 0 {
					cache.Color = product.Master.Fragments[0].Properties.Color
					if cache.Name == "" {
						cache.Name = product.Master.Fragments[0].Properties.ProductNameEn
					}
				}
				// 写入 Redis 缓存，1天
				cacheKey := fmt.Sprintf(baozun.ProductInfoRedisKey, req.EnterpriseId, productCode)
				if cacheBytes, mErr := json.Marshal(cache); mErr == nil {
					redis.Client().Set(c, cacheKey, string(cacheBytes), 24*time.Hour)
				}
				productInfoMap[productCode] = &cache
			}
		}
	}

	// 转换响应数据
	if len(talkieResp.Data) > 0 {
		for _, region := range talkieResp.Data {
			tmpRegion := sourceResp.TalkieImageSearchRegion{
				Region:  region.Region,
				PicList: make([]sourceResp.TalkieImageSearchPicInfo, 0),
			}
			for _, pic := range region.PicList {
				tmpPic := sourceResp.TalkieImageSearchPicInfo{
					ProductID:  pic.ProductID,
					PicName:    pic.PicName,
					CategoryID: pic.CategoryID,
					Score:      pic.Score,
				}
				// 从 ProductID 中截取 SKU（去掉颜色部分）
				if idx := strings.Index(pic.ProductID, "_"); idx > 0 {
					tmpPic.Sku = pic.ProductID[:idx]
				}
				if info, ok := productInfoMap[pic.ProductID]; ok && info != nil {
					tmpPic.Name = info.Name
					tmpPic.Upc = info.Upc
					tmpPic.Url = info.Url
					tmpPic.Color = info.Color
				}
				tmpRegion.PicList = append(tmpRegion.PicList, tmpPic)
			}
			resp = append(resp, tmpRegion)
		}
	}

	return resp, nil
}
