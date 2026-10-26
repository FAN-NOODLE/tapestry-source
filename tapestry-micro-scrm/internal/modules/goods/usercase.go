package goods

import (
	"database/sql"
	"encoding/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/uid"
	"time"
)

// Service 商品服务
type Service struct {
	repo Repository
}

// SkuList 列表
func (s Service) SkuList(c ctx.Context, cond *ListReq) (int64, []*SkuItem, error) {
	offset := (cond.Page - 1) * cond.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(cond.Size)}
	return s.repo.SkuList(c, pageVo, cond)
}

// SyncGoodsChange 同步商品变更事件数据
func (s *Service) SyncGoodsChange(c ctx.Context, request *SyncChangeRequest) error {
	var spuIDList []int64
	var skuIDList []int64
	var skuMap = make(map[int64]*SyncChangeSku)
	var spuMap = make(map[int64]*SyncChangeItem)
	for _, v := range request.Data.List {
		spuIDList = append(spuIDList, v.SpuID)
		spuMap[v.SpuID] = v
		for _, vv := range v.SkuList {
			skuIDList = append(skuIDList, vv.ID)
			skuMap[vv.ID] = vv
		}
	}
	dbSpuList, err := s.repo.QuerySpuList(c, spuIDList)
	if err != nil {
		log.Info(c).Msgf("spu列表查询失败, err:%s, spuIDList:%v", err.Error(), spuIDList)
		return err
	}
	err = s.UpdateSpuChanges(c, dbSpuList, spuMap)
	if err != nil {
		return err
	}
	dbSkuList, err := s.repo.QuerySkuList(c, skuIDList)
	if err != nil {
		log.Info(c).Msgf("sku列表查询失败, err:%s, skuIDList:%v", err.Error(), skuIDList)
		return err
	}
	err = s.UpdateSkuChanges(c, dbSkuList, skuMap)
	if err != nil {
		return err
	}
	return nil
}

// SyncGoodsDelete 同步商品删除事件
func (s *Service) SyncGoodsDelete(c ctx.Context, request *SyncDeleteRequest) error {
	var deleteSkuList []*SkuEntity
	now := time.Now().Unix()
	var err error
	if request.Data.IsDeleteSpu {
		// 删除spu
		spuInfo, err := s.QuerySpuInfo(c, request.Data.EnterpriseID, request.Data.SpuID)
		if err != nil {
			return err
		}
		if spuInfo == nil {
			return nil
		}
		spuInfo.DeletedAt = now
		_, err = s.repo.UpdateSpu(c, spuInfo)
		if err != nil {
			log.Info(c).Msgf("删除spu失败, err:%s, spu:%v", err.Error(), spuInfo)
			return err
		}
		// 获取spu下的所有sku列表，进行删除
		deleteSkuList, err = s.repo.QuerySkuListOfSpu(c, request.Data.EnterpriseID, request.Data.SpuID)
		if err != nil && err != sql.ErrNoRows {
			log.Info(c).Msgf("查询sku数据失败, err:%s, spu:%v", err.Error(), spuInfo)
			return err
		}
		if err == sql.ErrNoRows {
			return nil
		}
	} else {
		deleteSkuList, err = s.repo.QuerySkuList(c, request.Data.SkuIDs)
		if err != nil && err != sql.ErrNoRows {
			log.Info(c).Msgf("查询sku数据失败, err:%s, skuIDs:%v", err.Error(), request.Data.SkuIDs)
			return err
		}
		if err == sql.ErrNoRows {
			return nil
		}
	}
	if len(deleteSkuList) == 0 {
		return nil
	}
	// 删除sku列表
	for _, v := range deleteSkuList {
		v.DeletedAt = now
		_, err := s.repo.UpdateSku(c, v)
		if err != nil {
			log.Info(c).Msgf("删除sku数据失败: err:%s, sku:%v", err.Error(), v)
		}
	}
	return nil
}

// UpdateSpuChanges 更新spu更新数据
func (s *Service) UpdateSpuChanges(c ctx.Context, dbSpuList []*SpuEntity, requestSpuMap map[int64]*SyncChangeItem) error {
	for _, v := range dbSpuList {
		if spuInfo, ok := requestSpuMap[v.SpuID]; ok {
			v.Status = spuInfo.SpuStatus
			v.UpdatedAt = time.Now().Unix()
		}
		_, err := s.repo.UpdateSpu(c, v)
		if err != nil {
			log.Info(c).Msgf("更新spu数据失败, err:%s, spu:%v", err.Error(), v)
			return err
		}
	}
	return nil
}

// UpdateSkuChanges 更新sku变更信息
func (s *Service) UpdateSkuChanges(c ctx.Context, dbSkuList []*SkuEntity, requestSkuMap map[int64]*SyncChangeSku) error {
	for _, v := range dbSkuList {
		if skuInfo, ok := requestSkuMap[v.SkuID]; ok {
			v.Status = skuInfo.Status
			v.SourceID = skuInfo.SourceID
			v.Stock = skuInfo.Stock
			v.SellPrice = skuInfo.Price
			v.SignPrice = skuInfo.SignPrice
			v.UpdatedAt = time.Now().Unix()
		}
		_, err := s.repo.UpdateSku(c, v)
		if err != nil {
			log.Info(c).Msgf("更新sku数据失败, err:%s, sku:%v", err.Error(), v)
			return err
		}
	}
	return nil
}

// SyncActivityData 同步活动数据
func (s *Service) SyncGoodsData(c ctx.Context, request *SyncRequest) error {
	switch request.Method {
	case common.DdcEventEnum.GoodsBase:
		return s.syncGoodsBase(c, request)
	case common.DdcEventEnum.GoodsChange:

	}
	//case common.DdcEventEnum.ActivityUpdate:
	//
	//}
	return nil
}

// syncGoodsBase 同步商品基本信息
func (s *Service) syncGoodsBase(c ctx.Context, request *SyncRequest) error {
	spuEntity, skuList, err := s.translateGoodsData(c, request)
	if err != nil {
		log.Info(c).Msgf("解析商品数据失败, err:%s, request:%v", err.Error(), request)
		return err
	}
	if spuEntity.EnterpriseID == 0 || spuEntity.SpuID == 0 {
		log.Info(c).Msgf("spu企业id和spuID不能为空, spuEntity:%v", spuEntity)
		return errs.SpuBaseDataErr
	}
	dbSpu, err := s.QuerySpuInfo(c, spuEntity.EnterpriseID, spuEntity.SpuID)
	if err != nil {
		return err
	}
	if dbSpu == nil {
		dbSpu = spuEntity
		err = s.AddSpuData(c, dbSpu)
		if err != nil {
			return err
		}
		err = s.AddSkuData(c, dbSpu, skuList)
		if err != nil {
			return err
		}
		return nil
	}
	err = s.updateSpuInfo(c, dbSpu, spuEntity)
	if err != nil {
		return err
	}
	err = s.updateSkuList(c, dbSpu, skuList)
	if err != nil {
		return err
	}
	return nil
}

// updateSpuInfo 更新spu信息
func (s *Service) updateSpuInfo(c ctx.Context, dbSpu *SpuEntity, requestSpu *SpuEntity) error {
	requestSpu.ID = dbSpu.ID
	requestSpu.CreatedAt = dbSpu.CreatedAt
	requestSpu.UpdatedAt = time.Now().Unix()
	_, err := s.repo.UpdateSpu(c, requestSpu)
	if err != nil {
		log.Info(c).Msgf("更新spu数据失败, err:%s, spu:%v", err.Error(), requestSpu)
		return err
	}
	return nil
}

// updateSkuInfo 更新sku信息
func (s *Service) updateSkuInfo(c ctx.Context, dbSku *SkuEntity, requestSku *SkuEntity) error {
	requestSku.ID = dbSku.ID
	requestSku.CreatedAt = dbSku.CreatedAt
	requestSku.UpdatedAt = time.Now().Unix()
	_, err := s.repo.UpdateSku(c, requestSku)
	if err != nil {
		log.Info(c).Msgf("更新spu数据失败, err:%s, spu:%v", err.Error(), requestSku)
		return err
	}
	return nil
}

func (s *Service) updateSkuList(c ctx.Context, spuEntity *SpuEntity, skuList []*SkuEntity) error {
	dbSkuList, err := s.repo.QuerySkuListOfSpu(c, spuEntity.EnterpriseID, spuEntity.SpuID)
	if err != nil {
		log.Info(c).Msgf("查询spu下sku列表失败, err:%s, spu:%v", err.Error(), spuEntity)
		return err
	}
	var requestSkuMap = make(map[int64]*SkuEntity)
	for _, v := range skuList {
		requestSkuMap[v.SkuID] = v
	}
	for _, v := range dbSkuList {
		if requestSku, ok := requestSkuMap[v.SkuID]; ok {
			err = s.updateSkuInfo(c, v, requestSku)
			if err != nil {
				return err
			}
			delete(requestSkuMap, v.SkuID)
		} else {
			v.DeletedAt = time.Now().Unix()
			_, err = s.repo.UpdateSku(c, v)
			if err != nil {
				return err
			}
		}
	}
	var newSkuList []*SkuEntity
	for _, v := range requestSkuMap {
		newSkuList = append(newSkuList, v)
	}
	err = s.AddSkuData(c, spuEntity, newSkuList)
	if err != nil {
		return err
	}
	return nil
}

// QuerySpuInfo 从数据库查询spu信息
func (s *Service) QuerySpuInfo(c ctx.Context, enterpriseID int64, spuID int64) (*SpuEntity, error) {
	spu, err := s.repo.QuerySpu(c, enterpriseID, spuID)
	if err != nil && err != sql.ErrNoRows {
		log.Info(c).Msgf("查询spu数据失败,err:%s, enterpriseID:%d, spuID:%d", err.Error(), enterpriseID, spuID)
		return nil, err
	}
	return spu, nil
}

// AddSkuData 添加sku数据
func (s *Service) AddSkuData(c ctx.Context, spuEntity *SpuEntity, skuList []*SkuEntity) error {
	now := time.Now().Unix()
	for _, v := range skuList {
		v.ID = uid.ID(c)
		v.CreatedAt = now
		v.UpdatedAt = now
		_, err := s.repo.InsertSku(c, v)
		if err != nil {
			log.Info(c).Msgf("添加sku数据失败, err:%s, sku:%v", err.Error(), v)
			return err
		}
	}
	return nil
}

// AddSpuData 新增spu数据
func (s *Service) AddSpuData(c ctx.Context, entity *SpuEntity) error {
	entity.ID = uid.ID(c)
	entity.CreatedAt = time.Now().Unix()
	entity.UpdatedAt = entity.CreatedAt
	_, err := s.repo.InsertSpu(c, entity)
	if err != nil {
		log.Info(c).Msgf("添加spu失败, err:%s, entity:%v", err.Error(), entity)
		return err
	}
	return nil
}

// TranslateGoodsData 从商品同步请求参数中提取spu和sku数据
func (s Service) translateGoodsData(c ctx.Context, request *SyncRequest) (*SpuEntity, []*SkuEntity, error) {
	spuEntity, _ := s.translateGoodsSpu(c, request)
	skuList, _ := s.translateGoodsSku(c, request)
	log.Info(c).Msgf("spu:%v, sku:%v", spuEntity, skuList)
	return spuEntity, skuList, nil
}

func (s Service) translateGoodsSku(c ctx.Context, request *SyncRequest) ([]*SkuEntity, error) {
	attributeMap := s.getSpuAttributeMap(c, request.Data.Attributes) // 生成attributeMap
	var skuList []*SkuEntity
	for _, v := range request.Data.Sku.List {
		sku := v.SkuEntity
		var attributeList = make([]*StoreAttribute, 0)
		// 将请求数据中的SkuAttribute转换成数据中存储的格式
		for _, vv := range v.Attributes {
			attributeList = append(attributeList, &StoreAttribute{
				AttributeName:      attributeMap[vv.AttributeID],
				AttributeItemValue: vv.AttributeValue,
			})
		}
		bytes, _ := json.Marshal(attributeList)
		sku.SkuAttribute = string(bytes)
		sku.SkuPicture = v.Pictures.Image
		skuList = append(skuList, &sku)
	}
	return skuList, nil
}

// translateGoodsSpu 从请求参数中解析spu数据
func (s Service) translateGoodsSpu(c ctx.Context, request *SyncRequest) (*SpuEntity, error) {
	spuEntity := request.Data.SpuEntity
	spuEntity.Picture = request.Data.DefaultAttributes.Pictures.Image // 取pictures中的主图
	spuEntity.EnterpriseID = request.Data.DefaultAttributes.EnterpriseId
	attributeMap := s.getSpuAttributeMap(c, request.Data.Attributes) // 生成attributeMap
	defaultAttributes := request.Data.DefaultAttributes.Attributes
	var attributeList = make([]*StoreAttribute, 0)

	// 将请求数据中的SkuAttribute转换成数据中存储的格式
	for _, v := range defaultAttributes {
		attributeList = append(attributeList, &StoreAttribute{
			AttributeName:      attributeMap[v.AttributeID],
			AttributeItemValue: v.AttributeValue,
		})
	}
	bytes, _ := json.Marshal(attributeList)
	spuEntity.DefaultAttribute = string(bytes)
	return &spuEntity, nil
}

// getSpuAttributeMap 获取spu属性集合
func (s Service) getSpuAttributeMap(_ ctx.Context, attributes []*AttributeItem) map[int64]string {
	var attributeMap = make(map[int64]string)
	for _, v := range attributes {
		attributeMap[v.ID] = v.Name
	}
	return attributeMap
}

// 查询商品的 SPU list
// TODO 有些参数数据库并没有使用
func (s Service) SpuList(c ctx.Context, request *SpuListRequest) (int64, []*SpuEntity, error) {
	offset := (request.Page - 1) * request.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(request.Size)}
	return s.repo.SpuList(c, pageVo, request)
}
