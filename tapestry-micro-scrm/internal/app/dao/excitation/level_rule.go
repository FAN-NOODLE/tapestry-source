package excitationDao

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
)

type LevelRule struct {
	dao   dao.Dao
	model *excitationModels.ExcitationLevelRuleModel
}

func NewLevelRule() *LevelRule {
	return &LevelRule{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationLevelRuleModel),
	}
}

func (d *LevelRule) WithTx(tx *gorm.DB) *LevelRule {
	d.dao.DB = tx
	return d
}

func (d *LevelRule) Create(c app.Ctx, data excitationModels.ExcitationLevelRuleModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	// WithTx 切换db连接,用于事物
	return err
}

func (d *LevelRule) CreateInBatchesData(c app.Ctx, data []excitationModels.ExcitationLevelRuleModel) error {
	err := d.dao.WithContext(c).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *LevelRule) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *LevelRule) DeleteData(c app.Ctx, query interface{}) error {
	var model excitationModels.ExcitationLevelRuleModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Delete(&model).Error
	return err
}

func (d *LevelRule) GetDetail(c app.Ctx, cond excitationReq.LevelRuleDetailReq) ([]excitationResp.LevelRuleDetailResp, error) {
	var (
		resp []excitationResp.LevelRuleDetailResp
		err  error
	)

	query := d.dao.WithContext(c).Model(d.model).
		Where(map[string]interface{}{
			"enterprise_id": cond.EnterpriseId,
		})

	err = query.Order("min_value ASC").Find(&resp).Error

	return resp, err
}

// GetDetailList 导购等级规则信息
func (d *LevelRule) GetDetailList(c app.Ctx, cond excitationReq.LevelRuleDetailReq) ([]excitationResp.LevelRuleDetailInfoResp, error) {
	var (
		resp []excitationResp.LevelRuleDetailInfoResp
		err  error
	)

	query := d.dao.WithContext(c).Model(d.model).
		Where(map[string]interface{}{
			"enterprise_id": cond.EnterpriseId,
		})

	err = query.Order("min_value ASC").Find(&resp).Error

	return resp, err
}

//GetMaxDetail 查询等级最大的那一个
func (d *LevelRule) GetMaxDetail(c app.Ctx, enterpriseId int64, levelType int) (excitationResp.LevelRuleDetailResp, error) {
	var (
		resp excitationResp.LevelRuleDetailResp
		err  error
	)
	query := d.dao.WithContext(c).Model(d.model).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
		})

	err = query.Where("level_type > ?", levelType).Order("level_type ASC").Find(&resp).Error

	if resp.ID == 0 {
		query := d.dao.WithContext(c).Model(d.model).
			Where(map[string]interface{}{
				"enterprise_id": enterpriseId,
			})
		err = query.Where("level_type", levelType).Order("level_type ASC").Find(&resp).Error
		resp.MaxValue = resp.MinValue
	} else {
		//因为这里要球取最小值所以改成MinValue
		resp.MaxValue = resp.MinValue
	}

	return resp, err
}

// GetLevelInfoForPoint 根据积分查询等级规则 最大的一条
func (d *LevelRule) GetLevelInfoForPoint(c app.Ctx, enterpriseId int64, point string) (excitationResp.LevelRuleDetailInfoResp, error) {
	var (
		resp excitationResp.LevelRuleDetailInfoResp
		err  error
	)
	query := d.dao.WithContext(c).Model(d.model).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
		}).Where("min_value <= ? ", point)

	err = query.Order("min_value DESC").First(&resp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return resp, err
}

// GetLevelInfoForLevel 根据等级类型查询等级规则
func (d *LevelRule) GetLevelInfoForLevel(c app.Ctx, enterpriseId int64, levelType int) (excitationResp.LevelRuleDetailInfoResp, error) {
	var (
		resp excitationResp.LevelRuleDetailInfoResp
		err  error
	)
	query := d.dao.WithContext(c).Model(d.model).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"level_type":    levelType,
		})

	err = query.First(&resp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return resp, err
}

// GetLevelIcon 根据id查询等级图标
func (d *LevelRule) GetLevelIcon(c app.Ctx, enterpriseId int64, levelId int64) string {
	var icon string
	if levelId > 0 {
		query := d.dao.WithContext(c).Model(d.model).
			Where(map[string]interface{}{
				"enterprise_id": enterpriseId,
				"id":            levelId,
			})
		query.Pluck("level_icon", &icon)
	}
	return icon
}

// GetLevelType 根据id查询等级标识
func (d *LevelRule) GetLevelType(c app.Ctx, enterpriseId int64, levelId int64) int {
	var levelType int
	if levelId > 0 {
		query := d.dao.WithContext(c).Model(d.model).
			Unscoped().
			Where(map[string]interface{}{
				"enterprise_id": enterpriseId,
				"id":            levelId,
			})
		query.Pluck("level_type", &levelType)
	}
	return levelType
}

// GetLevelList 导购等级规则信息
func (d *LevelRule) GetLevelList(c app.Ctx, enterpriseId int64) ([]excitationResp.LevelRuleDetailInfoResp, error) {
	var (
		resp []excitationResp.LevelRuleDetailInfoResp
		err  error
	)

	query := d.dao.WithContext(c).Model(d.model).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
		})

	err = query.Order("level_type DESC").Find(&resp).Error

	return resp, err
}
