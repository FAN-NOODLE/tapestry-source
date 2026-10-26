package microScrmDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"time"
)

type MemberLearnRecord struct {
	dao       dao.Dao
	tableName string
}

func NewMemberLearnRecord() *MemberLearnRecord {
	return &MemberLearnRecord{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.MemberLearnRecordModel{}.TableName(),
	}
}

//GetLimitLearnRecord 获取未处理的积分记录
func (d *MemberLearnRecord) GetLimitLearnRecord(c app.Ctx, enterpriseId int64) (list []microScrmModel.MemberLearnRecordModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Limit(100).
		Find(&list).Error
	return
}

//DeleteLearnRecord 删除已处理的记录
func (d *MemberLearnRecord) DeleteLearnRecord(c app.Ctx, enterpriseId int64, ids []int64) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("id in ?", ids).
		UpdateColumn("deleted_at", time.Now().Unix()).Error

	return err
}

func (d *MemberLearnRecord) CreateInfo(c ctx.Context, enterpriseId int64, info microScrmModel.MemberLearnRecordModel) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&info).Error
}
