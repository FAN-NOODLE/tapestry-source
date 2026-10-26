package excitationDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
)

type MemberDataRecord struct {
	dao   dao.Dao
	model *excitationModels.ExcitationMemberDataRecordModel
}

func NewMemberDataRecord() *MemberDataRecord {
	return &MemberDataRecord{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(excitationModels.ExcitationMemberDataRecordModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *MemberDataRecord) WithTx(tx *gorm.DB) *MemberDataRecord {
	d.dao.DB = tx
	return d
}

func (d *MemberDataRecord) GetOne(c app.Ctx, query interface{}, args ...interface{}) (excitationModels.ExcitationMemberDataRecordModel, error) {
	var result excitationModels.ExcitationMemberDataRecordModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *MemberDataRecord) CreateData(c app.Ctx, data excitationModels.ExcitationMemberDataRecordModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *MemberDataRecord) CreateOne(c app.Ctx, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

// GetMemberLevelRecord 获取导购等级历史进程
func (d *MemberDataRecord) GetMemberLevelRecord(c app.Ctx, enterpriseId int64, memberId []int64, order string) (res []excitationResp.LevelRecordResp, err error) {
	err = d.dao.WithContext(c).Model(d.model).
		Select("created_at,level,member_id").
		Where("enterprise_id = ?", enterpriseId).
		Where("member_id in ?", memberId).
		Where("level != ''").
		//Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Order(order).
		Find(&res).Error
	return res, err
}

// GetMemberMaxLevelRecord 获取导购最大等级记录
func (d *MemberDataRecord) GetMemberMaxLevelRecord(c app.Ctx, enterpriseId int64, memberId []int64) (res []excitationResp.LevelRecordResp, err error) {
	err = d.dao.WithContext(c).Model(d.model).
		Select("created_at,level,member_id,CAST(REPLACE(level, 'LV', '') AS UNSIGNED) level_num").
		Where("enterprise_id = ?", enterpriseId).
		Where("member_id in ?", memberId).
		Where("level != ''").
		Order("level_num desc,created_at asc").
		Find(&res).Error
	return res, err
}
