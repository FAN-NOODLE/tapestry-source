package excitationDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
)

type MedalRecord struct {
	dao   dao.Dao
	model *excitationModels.ExcitationMedalRecordModel
}

func NewMedalRecord() *MedalRecord {
	return &MedalRecord{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationMedalRecordModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *MedalRecord) WithTx(tx *gorm.DB) *MedalRecord {
	d.dao.DB = tx
	return d
}

func (d *MedalRecord) GetPagination(c app.Ctx, cond excitationReq.MedalRecordList) ([]excitationResp.MedalRecordList, int64, error) {
	var (
		total int64
		resp  []excitationResp.MedalRecordList
		err   error
	)

	query := d.dao.WithContext(c).Model(d.model).
		Where("excitation_medal.enterprise_id = ?", cond.EnterpriseId).
		Select("excitation_medal_record.id,excitation_medal_record.created_at,excitation_medal_record.member_id," +
			"excitation_medal_record.ranking,excitation_medal_record.ranking_value,excitation_medal_rule.description," +
			"excitation_medal.name as rule_name,excitation_medal.source,excitation_medal.invite_top," +
			"excitation_medal.rule_dimension_id,excitation_medal.rule_category_id,excitation_medal.integral," +
			"excitation_medal.item_name,excitation_medal_style.medal_name").
		Joins("left join excitation_medal on excitation_medal.id = excitation_medal_record.medal_id").
		Joins("left join excitation_medal_style on excitation_medal_style.id = excitation_medal_record.medal_style_id").
		Joins("left join excitation_medal_rule on excitation_medal.rule_category_id = excitation_medal_rule.id").
		Scopes(func(db *gorm.DB) *gorm.DB {
			if cond.MedalId > 0 {
				db.Where("excitation_medal_record.medal_id = ?", cond.MedalId)
			}
			if len(cond.MemberIds) > 0 {
				db.Where("excitation_medal_record.member_id in ?", cond.MemberIds)
			}
			if cond.MedalRuleName != "" {
				db.Where("excitation_medal.name like ?", "%"+cond.MedalRuleName+"%")
			}
			if cond.Source > 0 {
				db.Where("excitation_medal.source = ?", cond.Source)
			}
			if cond.RuleDimensionId > 0 {
				db.Where("excitation_medal.rule_dimension_id = ?", cond.RuleDimensionId)
			}
			if cond.MedalName != "" {
				db.Where("excitation_medal_style.medal_name like ?", "%"+cond.MedalName+"%")
			}
			if cond.StartTime > 0 && cond.EndTime > 0 {
				db.Where("excitation_medal_record.created_at between ? and ?", cond.StartTime, cond.EndTime)
			}
			return db
		})
	query.Count(&total)

	offset := (cond.Page - 1) * cond.Size
	err = query.Order("excitation_medal_record.id DESC").Limit(cond.Size).Offset(offset).Scan(&resp).Error

	return resp, total, err
}

func (d *MedalRecord) GetOne(c app.Ctx, query interface{}, args ...interface{}) (excitationModels.ExcitationMedalRecordModel, error) {
	var result excitationModels.ExcitationMedalRecordModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *MedalRecord) GetMany(c app.Ctx, query interface{}) ([]excitationModels.ExcitationMedalRecordModel, error) {
	var result []excitationModels.ExcitationMedalRecordModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Find(&result).Error
	return result, err
}

func (d *MedalRecord) CreateData(c app.Ctx, data excitationModels.ExcitationMedalRecordModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *MedalRecord) CreateInBatchesData(c app.Ctx, data []excitationModels.ExcitationMedalRecordModel) error {
	err := d.dao.WithContext(c).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *MedalRecord) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *MedalRecord) DeleteData(c app.Ctx, query interface{}) error {
	var model excitationModels.ExcitationMedalRecordModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Delete(&model).Error
	return err
}

func (d *MedalRecord) ForceDelete(c app.Ctx, query interface{}) error {
	var model excitationModels.ExcitationMedalRecordModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Unscoped().Delete(&model).Error
	return err
}

func (d *MedalRecord) GetReceived(c app.Ctx, enterpriseId int64, memberId int64) ([]excitationModels.ExcitationMedalRecordModel, error) {
	var result []excitationModels.ExcitationMedalRecordModel
	err := d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id = ? and member_id = ? and is_receive = ?", enterpriseId, memberId, excitationModels.MedalRecordIsReceiveTrue).
		Order("created_at desc").
		Find(&result).Error
	return result, err
}

func (d *MedalRecord) GetGained(c app.Ctx, enterpriseId int64, memberId int64) ([]excitationModels.ExcitationMedalRecordModel, error) {
	var result []excitationModels.ExcitationMedalRecordModel
	err := d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id = ? and member_id = ?", enterpriseId, memberId).
		Order("created_at desc").
		Find(&result).Error
	return result, err
}

func (d *MedalRecord) GetRecordByMemberIdAndMedalId(c app.Ctx, enterpriseId int64, memberId, medalId, startTime, endTime int64) ([]excitationModels.ExcitationMedalRecordModel, error) {
	var result []excitationModels.ExcitationMedalRecordModel
	err := d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id = ? and member_id = ? and medal_id = ?", enterpriseId, memberId, medalId).
		Where("created_at between ? and ?", startTime, endTime).
		Find(&result).Error
	return result, err
}

func (d *MedalRecord) GetReceivedNumGroupStyleId(c app.Ctx, enterpriseId, startTime, endTime int64) ([]excitationModels.ReceivedNumGroupStyleId, error) {
	var result []excitationModels.ReceivedNumGroupStyleId
	err := d.dao.WithContext(c).Model(d.model).
		Select("medal_style_id,count(*) as total").
		Where("enterprise_id = ?", enterpriseId).
		Where("created_at between ? and ?", startTime, endTime).
		Group("medal_style_id").
		Scan(&result).Error
	return result, err
}

func (d *MedalRecord) GetLastMemberReceived(c app.Ctx, enterpriseId int64, memberId int64) (excitationModels.ExcitationMedalRecordModel, error) {
	var result excitationModels.ExcitationMedalRecordModel
	err := d.dao.WithContext(c).Model(d.model).
		Where("enterprise_id = ? and member_id = ?", enterpriseId, memberId).
		Order("created_at desc").
		First(&result).Error
	return result, err
}

// GetExistsByMedalIds 获取是否存在某规则的奖章
func (d *MedalRecord) GetExistsByMedalIds(c app.Ctx, enterpriseId int64, medalIds []int64) (result []int64, err error) {
	err = d.dao.WithContext(c).Model(d.model).
		Select("medal_id").
		Where("enterprise_id = ? and medal_id in ?", enterpriseId, medalIds).
		Group("medal_id").
		Pluck("medal_id", &result).Error
	return
}
