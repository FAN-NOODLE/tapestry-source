package excitationDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
	"time"
)

type Medal struct {
	dao   dao.Dao
	model *excitationModels.ExcitationMedalModel
}

func NewMedal() *Medal {
	return &Medal{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationMedalModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *Medal) WithTx(tx *gorm.DB) *Medal {
	d.dao.DB = tx
	return d
}

func (d *Medal) GetPagination(c app.Ctx, cond excitationReq.MedalList) ([]excitationModels.ExcitationMedalModel, int64, error) {
	var (
		total int64
		resp  []excitationModels.ExcitationMedalModel
		err   error
	)

	query := d.dao.WithContext(c).Model(d.model).Where("enterprise_id = ?", cond.EnterpriseId).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if cond.IsEffective > 0 {
				// 有效筛选
				nowTime := time.Now().Unix()
				if cond.IsEffective == 1 {
					db.Where("rule_start_time <= ? and rule_end_time > ?", nowTime, nowTime).Or("rule_type = ?", excitationModels.MedalRuleTypeLong)
				}
				// 无效筛选
				if cond.IsEffective == 2 {
					db.Where("rule_start_time >= ? and rule_end_time < ?", nowTime, nowTime).Or("rule_type = ?", excitationModels.MedalRuleTypeOnlyOne)
				}
			}
			if cond.Name != "" {
				db.Where("name like ?", "%"+cond.Name+"%")
			}
			if cond.Source != 0 {
				db.Where("source = ?", cond.Source)
			}
			if cond.RuleDimensionId != 0 {
				db.Where("rule_dimension_id = ?", cond.RuleDimensionId)
			}
			if cond.RuleCategoryId != 0 {
				db.Where("rule_category_id = ?", cond.RuleCategoryId)
			}
			if cond.OtherRewards != 0 {
				db.Where("other_rewards = ?", cond.OtherRewards)
			}
			if cond.StartTime > 0 && cond.EndTime > 0 {
				db.Where("created_at between ? and ?", cond.StartTime, cond.EndTime)
			}
			if cond.RuleStartTime > 0 && cond.RuleEndTime > 0 {
				db.Where("rule_start_time >= ? and rule_end_time <= ?", cond.RuleStartTime, cond.RuleEndTime).
					Or("rule_type = ?", excitationModels.MedalRuleTypeLong)
			}
			if cond.Status > 0 {
				db.Where("status = ? ", cond.Status)
			}
			return db
		})

	query.Count(&total)

	offset := (cond.Page - 1) * cond.Size
	err = query.Order("id DESC,created_at DESC").Limit(cond.Size).Offset(offset).Find(&resp).Error

	return resp, total, err
}

func (d *Medal) GetOne(c app.Ctx, query interface{}, args ...interface{}) (excitationModels.ExcitationMedalModel, error) {
	var result excitationModels.ExcitationMedalModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Preload("MedalRule").Find(&result).Error
	return result, err
}

func (d *Medal) GetUnscopedOne(c app.Ctx, query interface{}, args ...interface{}) (excitationModels.ExcitationMedalModel, error) {
	var result excitationModels.ExcitationMedalModel
	err := d.dao.WithContext(c).Model(d.model).Unscoped().Where(query, args...).Find(&result).Error
	return result, err
}

func (d *Medal) GetOneDetail(c app.Ctx, query interface{}, args ...interface{}) (excitationModels.ExcitationMedalModel, error) {
	var result excitationModels.ExcitationMedalModel
	err := d.dao.WithContext(c).Model(d.model).
		Unscoped().
		Preload("MedalRule").
		Preload("MedalStyle", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Where(query, args...).Find(&result).Error
	return result, err
}

func (d *Medal) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]excitationModels.ExcitationMedalModel, error) {
	var result []excitationModels.ExcitationMedalModel
	err := d.dao.WithContext(c).Model(d.model).Preload("MedalRule").Preload("MedalStyle").Where(query, args...).Find(&result).Error
	return result, err
}

func (d *Medal) GetManyUnscoped(c app.Ctx, query interface{}, args ...interface{}) ([]excitationModels.ExcitationMedalModel, error) {
	var result []excitationModels.ExcitationMedalModel
	err := d.dao.WithContext(c).Model(d.model).
		Unscoped().
		Preload("MedalRule").
		Preload("MedalStyle", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Where(query, args...).
		Order("created_at desc").
		Find(&result).Error
	return result, err
}

func (d *Medal) CreateData(c app.Ctx, data excitationModels.ExcitationMedalModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *Medal) SaveData(c app.Ctx, data excitationModels.ExcitationMedalModel) error {
	err := d.dao.WithContext(c).Model(d.model).Save(&data).Error
	return err
}

func (d *Medal) UpdateData(c app.Ctx, query interface{}, data interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *Medal) DeleteData(c app.Ctx, query interface{}) error {
	var model excitationModels.ExcitationMedalModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Delete(&model).Error
	return err
}

func (d *Medal) ForceDelete(c app.Ctx, query interface{}) error {
	var model excitationModels.ExcitationMedalModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Unscoped().Delete(&model).Error
	return err
}
