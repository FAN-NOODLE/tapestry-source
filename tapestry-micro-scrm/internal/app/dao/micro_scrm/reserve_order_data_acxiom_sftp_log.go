package microScrmDao

import (
	"errors"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type ReserveOrderDataAcxiomSftpLog struct {
	dao       dao.Dao
	tableName string
}

func NewReserveOrderDataAcxiomSftpLogDao() *ReserveOrderDataAcxiomSftpLog {
	return &ReserveOrderDataAcxiomSftpLog{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.ReserveOrderDataAcxiomSftpLogModel{}.TableName(),
	}
}

// GetInfo 获取sftp信息记录
func (d *ReserveOrderDataAcxiomSftpLog) GetInfo(c ctx.Context, enterpriseId int64, fileType int, fileTime int64, pullType string) (data microScrmModel.ReserveOrderDataAcxiomSftpLogModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).Where("enterprise_id=? and file_type=? and file_time=? ", enterpriseId, fileType, fileTime).Scopes(func(db *gorm.DB) *gorm.DB {
		if pullType != "" {
			db.Where("push_type=?", pullType)
		}
		return db
	}).First(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

// CreateLog 操作日志记录
func (d *ReserveOrderDataAcxiomSftpLog) CreateLog(c ctx.Context, data microScrmModel.ReserveOrderDataAcxiomSftpLogModel) error {
	return d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
}

// UpdateLogStatus 操作日志更改状态为已完成
func (d *ReserveOrderDataAcxiomSftpLog) UpdateLogStatus(c ctx.Context, enterpriseId, id int64) error {
	return d.dao.WithContext(c).Table(d.tableName).Where("id=? and enterprise_id=?", id, enterpriseId).Update("status", 1).Error
}

func (d *ReserveOrderDataAcxiomSftpLog) GetExistsFile(c ctx.Context, enterpriseId int64, fileType int, pullType string) (data []microScrmModel.ReserveOrderDataAcxiomSftpLogModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).Where("enterprise_id=? and file_type=?", enterpriseId, fileType).Scopes(func(db *gorm.DB) *gorm.DB {
		if pullType != "" {
			db.Where("push_type=?", pullType)
		}
		return db
	}).Find(&data).Error
	return data, nil
}
