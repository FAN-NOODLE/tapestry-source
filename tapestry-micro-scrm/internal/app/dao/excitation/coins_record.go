package excitationDao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
)

type CoinsRecord struct {
	dao       dao.Dao
	model     *microScrmModel.ExcitationCoinsRecordModel
	tableName string
}

func NewCoinsRecord() *CoinsRecord {
	return &CoinsRecord{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		model:     new(microScrmModel.ExcitationCoinsRecordModel),
		tableName: microScrmModel.ExcitationCoinsRecordModel{}.TableName(),
	}
}

// WithTx 切换db连接,用于事物
func (d *CoinsRecord) WithTx(tx *gorm.DB) *CoinsRecord {
	d.dao.DB = tx
	return d
}

func (d *CoinsRecord) ReleaseTx() {
	d.dao.DB = gormx.GetInstance()
}

func (d *CoinsRecord) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.ExcitationCoinsRecordModel, error) {
	var result microScrmModel.ExcitationCoinsRecordModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *CoinsRecord) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.ExcitationCoinsRecordModel, error) {
	var result []microScrmModel.ExcitationCoinsRecordModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *CoinsRecord) GetPagination(c app.Ctx, cond excitationReq.CoinsRecordListReq) ([]microScrmModel.ExcitationCoinsRecordModel, int64, error) {
	var (
		total int64
		resp  []microScrmModel.ExcitationCoinsRecordModel
		err   error
	)

	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, cond.EnterpriseId)).Model(d.model).
		Scopes(func(db *gorm.DB) *gorm.DB {
			// 导购维度
			if len(cond.MemberIdArr) > 0 {
				db.Where("member_id in ?", cond.MemberIdArr)
			}
			// 金币来源
			if cond.SourceType > 0 {
				db.Where("source_type = ?", cond.SourceType)
			}
			// 金币变动类型
			if cond.ChangeType > 0 {
				db.Where(fmt.Sprintf("coins %s 0", helper.RuleThree(cond.ChangeType == excitationReq.ChangeTypeInc, ">", "<")))
			}
			// 获得时间
			if cond.StartTime > 0 && cond.EndTime > 0 {
				db.Where("created_at between ? and ?", cond.StartTime, cond.EndTime)
			}
			return db
		})
	query.Count(&total)

	offset := (cond.Page - 1) * cond.Size
	err = query.Order("created_at DESC").Order("id DESC").Limit(cond.Size).Offset(offset).Scan(&resp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return resp, total, err
}

func (d *CoinsRecord) CreateData(c app.Ctx, enterpriseId int64, data microScrmModel.ExcitationCoinsRecordModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Create(&data).Error
	return err
}

func (d *CoinsRecord) CreateInBatchesData(c app.Ctx, enterpriseId int64, data []microScrmModel.ExcitationCoinsRecordModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *CoinsRecord) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *CoinsRecord) DeleteData(c app.Ctx, enterpriseId int64, query interface{}) error {
	var model microScrmModel.ExcitationCoinsRecordModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).Delete(&model).Error
	return err
}

// GetSpendCoinsByMemberIds 获取总的消费记录
func (d *CoinsRecord) GetSpendCoinsByMemberIds(c app.Ctx, enterpriseId int64, memberIds []int64) (result []excitationResp.SpendCoinsResp, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).
		Select("member_id, sum(coins) as spend_coins").
		Where("member_id in ?", memberIds).
		Where("coins < 0").
		Group("member_id").Find(&result).Error
	return result, err
}
