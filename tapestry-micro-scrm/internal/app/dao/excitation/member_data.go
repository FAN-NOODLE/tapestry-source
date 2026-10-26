package excitationDao

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
	excitationModels "micro-scrm/internal/app/models/excitation"
	membersSvc "micro-scrm/internal/app/service/member"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"time"
)

type MemberData struct {
	dao   dao.Dao
	model *excitationModels.ExcitationMemberDataModel
}

func NewMemberData() *MemberData {
	return &MemberData{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationMemberDataModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *MemberData) WithTx(tx *gorm.DB) *MemberData {
	d.dao.DB = tx
	return d
}

// GetList 后台查询列表
func (d *MemberData) GetList(c app.Ctx, cond excitationReq.MemberDataListReq) ([]excitationResp.MemberDataListResp, int64, error) {
	var (
		total int64
		resp  []excitationResp.MemberDataListResp
		err   error
	)

	query := d.dao.WithContext(c).Model(d.model).
		Select("id, member_id, point_value, point_updated_at, level, coins, spend_coins, total_coins, level_id, level_updated_at, medal_num, medal_updated_at").
		Where(map[string]interface{}{
			"enterprise_id": cond.EnterpriseId,
		})

	if cond.MemberUserid != "" {
		memberIds, errs := membersSvc.NewMembersService().GetMemberIdsByUserid(c, cond.EnterpriseId, cond.MemberUserid)
		if err != nil {
			return resp, total, errs
		}
		query.Where("member_id in (?)", memberIds)
	}
	if cond.MemberName != "" {
		memberIds, errs := membersSvc.NewMembersService().GetMemberIdsByName(c, cond.EnterpriseId, cond.MemberName)
		if err != nil {
			return resp, total, errs
		}
		query.Where("member_id in (?)", memberIds)
	}
	if cond.MemberThirdID > 0 {
		var thirdIds []int
		thirdIds = append(thirdIds, cond.MemberThirdID)
		memberIds, errs := membersSvc.NewMembersService().GetMemberIdsByThirdIds(c, cond.EnterpriseId, thirdIds)
		if err != nil {
			return resp, total, errs
		}
		query.Where("member_id in (?)", memberIds)
	}

	if cond.Level != "" {
		query.Where("level = ?", cond.Level)
	}
	if cond.PointUpdatedStartTime > 0 && cond.PointUpdatedEndTime > 0 {
		query.Where("point_updated_at >= ? and point_updated_at < ?", cond.PointUpdatedStartTime, cond.PointUpdatedEndTime)
	}
	if cond.LevelUpdatedStartTime > 0 && cond.LevelUpdatedEndTime > 0 {
		query.Where("level_updated_at >= ? and level_updated_at < ?", cond.LevelUpdatedStartTime, cond.LevelUpdatedEndTime)
	}
	if cond.MedalUpdatedStartTime > 0 && cond.MedalUpdatedEndTime > 0 {
		query.Where("medal_updated_at >= ? and medal_updated_at < ?", cond.MedalUpdatedStartTime, cond.MedalUpdatedEndTime)
	}
	query.Count(&total)

	if cond.PointValue == "asc" || cond.PointValue == "desc" {
		query.Order("CAST(point_value AS SIGNED) " + cond.PointValue)
	} else if cond.Coins == "asc" || cond.Coins == "desc" {
		query.Order("coins " + cond.Coins)
	} else if cond.SpendCoins == "asc" || cond.SpendCoins == "desc" {
		query.Order("spend_coins " + cond.SpendCoins)
	} else if cond.TotalCoins == "asc" || cond.TotalCoins == "desc" {
		query.Order("total_coins " + cond.TotalCoins)
	} else if cond.MedalNum == "asc" || cond.MedalNum == "desc" {
		query.Order("medal_num " + cond.MedalNum)
	}

	offset := (cond.Page - 1) * cond.Size
	err = query.Order("id DESC").
		Limit(cond.Size).Offset(offset).
		Find(&resp).Error

	return resp, total, err
}

func (d *MemberData) DataIndex(c app.Ctx, enterpriseId int64, memberId int64) (excitationModels.ExcitationMemberDataModel, error) {
	var (
		resp excitationModels.ExcitationMemberDataModel
		err  error
	)

	query := d.dao.WithContext(c).Model(d.model).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"member_id":     memberId,
		})
	err = query.Find(&resp).Error

	return resp, err
}

func (d *MemberData) GetLevelCount(c app.Ctx, enterpriseId int64, levelType int) int64 {

	var count int64

	d.dao.WithContext(c).Model(d.model).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
		}).
		Where("level_type < ? and level_type > 0", levelType).
		Count(&count)

	return count
}

func (d *MemberData) GetMemberCount(c app.Ctx, enterpriseId int64) int64 {
	var count int64

	d.dao.WithContext(c).Model(d.model).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
		}).
		Count(&count)

	return count
}

func (d *MemberData) GetOne(c app.Ctx, query interface{}, args ...interface{}) (excitationModels.ExcitationMemberDataModel, error) {
	var result excitationModels.ExcitationMemberDataModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *MemberData) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).Updates(data).Error
	return err
}

func (d *MemberData) CreateData(c app.Ctx, data excitationModels.ExcitationMemberDataModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *MemberData) CreateOne(c app.Ctx, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *MemberData) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]excitationModels.ExcitationMemberDataModel, error) {
	var result []excitationModels.ExcitationMemberDataModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

// GetManyLimit 分页批量获取导购等级列表
func (d *MemberData) GetManyLimit(c app.Ctx, enterpriseId int64, pageLimit int, beginID int64) ([]excitationModels.ExcitationMemberDataModel, error) {
	var result []excitationModels.ExcitationMemberDataModel
	err := d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id = ? ", enterpriseId).
		Where("id > ? ", beginID).
		Limit(pageLimit).
		Find(&result).Error
	return result, err
}

// GetMemberIdsByQuery 根据条件获取所有导购id
func (d *MemberData) GetMemberIdsByQuery(c app.Ctx, query interface{}, args ...interface{}) (result []int64, err error) {
	err = d.dao.WithContext(c).Model(d.model).Where(query, args...).Pluck("member_id", &result).Error
	return result, err
}

// CoinsExpireByMemberIds  金币过期
func (d *MemberData) CoinsExpireByMemberIds(c app.Ctx, enterpriseId int64, memberIds []int64) error {
	return d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id = ? AND member_id in ?", enterpriseId, memberIds).
		UpdateColumns(map[string]interface{}{
			"coins":            0,
			"coins_updated_at": time.Now().Unix(),
		}).Error
}

// SaveCoins 金币更新,更新完返回真实的更新金币金额
func (d *MemberData) SaveCoins(c ctx.Context, params excitationReq.UpdateCoinsReq) (float64, error) {
	var coinsChange = params.Coins
	info, err := d.GetOne(c, "enterprise_id = ? and member_id = ?", params.EnterpriseId, params.MemberId)
	if err != nil {
		return coinsChange, err
	}
	if info.ID > 0 {
		if info.Coins+params.Coins < 0 {
			//金币不能小于0
			coinsChange = -info.Coins //如果扣减后金币为0 则扣减的金币为剩余金币
		}
		updateData := map[string]interface{}{
			"coins":            gorm.Expr(fmt.Sprintf("coins + ?"), helper.MathRounded(coinsChange, 1)),
			"coins_updated_at": time.Now().Unix(),
		}
		if coinsChange > 0 {
			updateData["total_coins"] = helper.MathRounded(info.TotalCoins+coinsChange, 1)
		} else if coinsChange < 0 {
			updateData["spend_coins"] = helper.MathRounded(info.SpendCoins+coinsChange, 1)
		}
		err = d.UpdateData(c,
			map[string]interface{}{
				"enterprise_id": params.EnterpriseId,
				"id":            info.ID,
			}, updateData)
		if err != nil {
			return coinsChange, err
		}
	}
	return coinsChange, nil
}
