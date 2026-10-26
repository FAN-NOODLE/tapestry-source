package activity

import (
	"database/sql"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/uid"
	"time"
)

// Service 活动服务
type Service struct {
	repo Repository
}

// SyncActivityData 同步活动数据
func (s *Service) SyncActivityData(c ctx.Context, request *SyncRequest) error {
	switch request.Method {
	case common.DdcEventEnum.ActivityAdd:
		return s.AddActivity(c, &request.Data.ShopActivityEntity, request.Data.SkuList)
	case common.DdcEventEnum.ActivityUpdate:
		return s.UpdateActivity(c, &request.Data.ShopActivityEntity, request.Data.SkuList)
	case common.DdcEventEnum.ActivityDelete:
		return s.DeleteActivity(c, &request.Data.ShopActivityEntity)
	}
	return nil
}

// AddActivity 新增活动
func (s *Service) AddActivity(c ctx.Context, activity *ShopActivityEntity, skuList []*ShopActivitySku) error {
	if activity.EnterpriseID == 0 || activity.ActivityID == 0 || len(skuList) == 0 {
		log.Info(c).Msgf("企业id,活动id和spu都不能为空, activity:%v, skuList:%v", activity, skuList)
		return nil
	}
	dbActivity, err := s.GetActivity(c, activity.EnterpriseID, activity.ActivityID)
	if err != nil {
		return err
	}
	// 数据库中已存在活动，不做任何处理
	if dbActivity != nil {
		return nil
	}
	err = s.AddActivityToDB(c, activity)
	if err != nil {
		return err
	}
	err = s.AddActivitySku(c, skuList, activity)
	if err != nil {
		return err
	}
	return nil
}

// UpdateActivity 更新活动信息
func (s *Service) UpdateActivity(c ctx.Context, activity *ShopActivityEntity, skuList []*ShopActivitySku) error {
	dbActivity, err := s.GetActivity(c, activity.EnterpriseID, activity.ActivityID)
	if err != nil {
		return err
	}
	if dbActivity == nil {
		return s.AddActivity(c, activity, skuList) // 数据库未查询到活动数据，走新增逻辑
	}
	activity.ID = dbActivity.ID
	activity.CreatedAt = dbActivity.CreatedAt
	err = s.UpdateActivityToDB(c, activity)
	if err != nil {
		return err
	}
	err = s.UpdateActivitySkuToDB(c, activity, skuList)
	if err != nil {
		return err
	}
	return nil
}

// DeleteActivity 删除活动数据
func (s *Service) DeleteActivity(c ctx.Context, activity *ShopActivityEntity) error {
	dbActivity, err := s.GetActivity(c, activity.EnterpriseID, activity.ActivityID)
	if err != nil {
		return err
	}
	if dbActivity == nil {
		return nil
	}
	dbActivity.DeletedAt = time.Now().Unix()
	err = s.UpdateActivityToDB(c, dbActivity)
	if err != nil {
		return err
	}
	return nil
}

// GetActivity 获取活动信息
func (s *Service) GetActivity(c ctx.Context, enterpriseID int64, activityID int64) (*ShopActivityEntity, error) {
	data, err := s.repo.GetActivity(c, enterpriseID, activityID)
	if err != nil && err != sql.ErrNoRows {
		log.Info(c).Msgf("查询活动信息失败, err:%s, enterprise:%d, activity:%d", err.Error(), enterpriseID, activityID)
		return nil, err
	}
	return data, nil
}

// AddActivityToDB 添加活动数据到数据库
func (s *Service) AddActivityToDB(c ctx.Context, data *ShopActivityEntity) error {
	data.ID = uid.ID(c)
	if data.CreatedAt == 0 {
		data.CreatedAt = time.Now().Unix()
	}
	data.UpdatedAt = data.CreatedAt
	_, err := s.repo.AddActivity(c, data)
	if err != nil {
		log.Info(c).Msgf("添加活动数据失败, err:%s, activity:%v", err.Error(), data)
		return err
	}
	return nil
}

// AddActivitySku 添加活动spu数据
func (s *Service) AddActivitySku(c ctx.Context, skuList []*ShopActivitySku, activity *ShopActivityEntity) error {
	if len(skuList) == 0 {
		return nil
	}
	for _, v := range skuList {
		v.ID = uid.ID(c)
		v.ShopActivityID = activity.ID
		v.EnterpriseID = activity.EnterpriseID
		v.CreatedAt = time.Now().Unix()
		v.UpdatedAt = v.CreatedAt
	}
	_, err := s.repo.AddActivitySku(c, skuList)
	if err != nil {
		log.Info(c).Msgf("添加活动spu数据失败, err:%s, skuList:%v", err.Error(), skuList)
		return err
	}
	return nil
}

// UpdateActivityToDB 更新活动数据到数据库
func (s *Service) UpdateActivityToDB(c ctx.Context, data *ShopActivityEntity) error {
	data.UpdatedAt = time.Now().Unix()
	_, err := s.repo.UpdateActivity(c, data)
	if err != nil {
		log.Info(c).Msgf("更新活动数据失败, err:%s, activity:%v", err.Error(), data)
		return err
	}
	return nil
}

// UpdateActivitySkuToDB 更新活动spu数据
func (s *Service) UpdateActivitySkuToDB(c ctx.Context, activity *ShopActivityEntity, skuList []*ShopActivitySku) error {
	dbSkuList, err := s.repo.GetActivitySkuList(c, activity.EnterpriseID, activity.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Info(c).Msgf("获取活动spu列表失败, err:%s, activity:%v", err.Error(), activity)
		return err
	}
	var requestSkuMap = make(map[int64]*ShopActivitySku)
	for _, v := range skuList {
		requestSkuMap[v.SkuID] = v
	}
	var updateSkuList []*ShopActivitySku
	for _, v := range dbSkuList {
		if newSku, ok := requestSkuMap[v.SkuID]; ok {
			newSku.ID = v.ID
			newSku.ShopActivityID = v.ShopActivityID
			newSku.EnterpriseID = v.EnterpriseID
			newSku.CreatedAt = v.CreatedAt
			newSku.UpdatedAt = time.Now().Unix()
			updateSkuList = append(updateSkuList, newSku)
			delete(requestSkuMap, v.SkuID)
		} else {
			// 删除请求数据中没有得sku
			v.DeletedAt = time.Now().Unix()
			updateSkuList = append(updateSkuList, v)
		}
	}
	var newSkuList []*ShopActivitySku
	for _, v := range requestSkuMap {
		newSkuList = append(newSkuList, v)
	}
	if len(newSkuList) > 0 {
		err = s.AddActivitySku(c, newSkuList, activity)
		if err != nil {
			log.Info(c).Msgf("新增活动sku数据失败, err:%s, data:%v", err.Error(), newSkuList)
		}
	}
	// 删除和更新sku数据
	if len(updateSkuList) > 0 {
		for _, v := range updateSkuList {
			_, err := s.repo.UpdateActivitySku(c, v)
			if err != nil {
				log.Info(c).Msgf("更新sku数据失败, err:%s, data:%v", err.Error(), v)
			}
		}
	}

	return nil
}

// AllActivity 查询所有的活动
func (s Service) AllActivity(c ctx.Context, request *ShopActivityEntity) ([]*ShopActivityEntity, error) {
	return s.repo.AllActivity(c, request)
}

// SkuListByActivity 通过活动id 查询
func (s Service) SkuListByActivity(c ctx.Context, request SkuListByActivityIdsReq) (int64, []*GoodsListItem, error) {
	offset := (request.Page - 1) * request.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(request.Size)}
	return s.repo.SkuListByActivityIds(c, request, pageVo)
}
