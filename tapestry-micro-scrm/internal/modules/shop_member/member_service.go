package shop_member

import (
	"database/sql"
	"micro-scrm/internal/pkg/helper"
	"strconv"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/lock"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/pkg/common"
)

type Service struct {
	Repo Repository
}

// SyncShopMemberLevel 同步商城会员等级数据
func (s *Service) SyncShopMemberLevel(c ctx.Context, request *ShopMemberLevelEntity, event common.DdcEvent) (bool, error) {
	if request.LevelID == 0 || request.Level == 0 || request.LevelName == "" || request.EnterpriseID == 0 {
		return false, nil
	}
	key := helper.Md5V(strconv.FormatInt(request.EnterpriseID, 10) + strconv.FormatInt(request.LevelID, 10) + event.String())
	mtx, err := lock.Mutex(
		c,
		lock.Key(key),
		lock.TTL(3*time.Second),
		lock.Retry(lock.NoRetry()),
	)
	if err != nil {
		log.Error(c).Msgf("同步商城会员等级数据生成锁失败-LevelID=%d，错误信息：%+v", request.LevelID, err)
		return false, nil
	}
	lockErr := mtx.Lock()
	if lockErr != nil {
		log.Error(c).Msgf("同步商城会员等级数据加成锁失败-LevelID=%d，错误信息：%+v", request.LevelID, lockErr)
		return false, nil
	}
	defer mtx.Unlock()
	// 1. 根据level_id查询数据
	levelInfo, err := s.Repo.GetShopMemberLevel(c, request.EnterpriseID, request.LevelID)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	// 2. 数据为空，新增数据：level_id为空新增shop_member_level表
	if levelInfo.LevelID == 0 {
		_, err := s.Repo.AddShopMemberLevel(c, request)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	// 3. 数据不为空，判断ddc同步事件
	// 3.1 新增事件
	if event == common.DdcEventEnum.ShopMemberLevelAdd {
		return true, nil
	}
	// 3.2 更新事件，更新shop_member_level表数据
	request.ID = levelInfo.ID
	_, err = s.Repo.UpdateShopMemberLevel(c, request)
	if err != nil {
		log.Info(c).Msgf("更新商城会员等级数据失败,err:%s,level_id:%d, data:%+v", err.Error(), levelInfo.LevelID, request)
		return false, err
	}
	log.Debug(c).Msgf("更新商城会员等级数据成功: levelInfo:%v", levelInfo)
	return true, nil
}

// LevelList 等级列表
func (s *Service) LevelList(c ctx.Context, enterpriseID int64, levelIDs []int64) ([]*ShopMemberLevelEntity, error) {
	list, err := s.Repo.GetLevelList(c, enterpriseID, levelIDs)
	if err != nil {
		log.Info(c).Msgf("查询会员等级列表失败, err:%s, enterpriseID:%d", err.Error(), enterpriseID)
		return list, err
	}
	return list, nil
}
