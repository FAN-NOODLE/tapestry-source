package microBossDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
	"micro-scrm/internal/pkg/db/gormx"
	"time"
)

type EnterpriseAdmin struct {
	dao       dao.Dao
	tableName string
}

func NewEnterpriseAdmin() *EnterpriseAdmin {
	return &EnterpriseAdmin{
		dao:       dao.Dao{DB: gormx.GetBossInstance()},
		tableName: microBossModel.EnterpriseAdminModel{}.TableName(),
	}
}

func (d *EnterpriseAdmin) GetInfo(c app.Ctx, data microBossModel.EnterpriseAdminModel) (microBossModel.EnterpriseAdminModel, error) {
	var result microBossModel.EnterpriseAdminModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(data).First(&result).Error
	return result, err
}

func (d *EnterpriseAdmin) GetMany(c app.Ctx, query interface{}) ([]microBossModel.EnterpriseAdminModel, error) {
	var result []microBossModel.EnterpriseAdminModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query).Find(&result).Error
	return result, err
}

// GetActiveByMobiles 按员工手机号批量查找当前企业管理员。
func (d *EnterpriseAdmin) GetActiveByMobiles(c app.Ctx, enterpriseID int64, mobiles []string) ([]microBossModel.EnterpriseAdminModel, error) {
	if len(mobiles) == 0 {
		return nil, nil
	}
	var result []microBossModel.EnterpriseAdminModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("id, enterprise_id, mobile").
		Where("enterprise_id = ? AND deleted_at = 0 AND mobile IN ?", enterpriseID, mobiles).
		Order("id ASC").Find(&result).Error
	return result, err
}

func (d *EnterpriseAdmin) UpdateMobileByOldMobile(c app.Ctx, oldMobile, mobile string, enterpriseId int64) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"mobile":        oldMobile,
		"enterprise_id": enterpriseId,
	}).Updates(map[string]interface{}{
		"mobile": mobile,
	}).Error
}

func (d *EnterpriseAdmin) RemoveByMobile(c app.Ctx, mobile string, enterpriseId int64) error {
	return d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"mobile":        mobile,
		}).
		Update("deleted_at", time.Now().Unix()).Error
}

func (d *EnterpriseAdmin) BatchRemoveByMobile(c app.Ctx, mobiles []string, enterpriseId int64) error {
	return d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"mobile":        mobiles,
			"deleted_at":    0,
		}).
		Update("deleted_at", time.Now().Unix()).Error
}
