package microScrmDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type OrderAchievementSftpLog struct {
	dao       dao.Dao
	tableName string
}

func NewOrderAchievementSftpLog() *OrderAchievementSftpLog {
	return &OrderAchievementSftpLog{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.OrderAchievementSftpLogModel{}.TableName(),
	}
}

func (d OrderAchievementSftpLog) GetLogList(c app.Ctx, enterpriseId int64) ([]microScrmModel.OrderAchievementSftpLogModel, error) {
	var (
		result []microScrmModel.OrderAchievementSftpLogModel
		err    error
	)
	err = d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id", enterpriseId).
		Find(&result).Error

	return result, err
}

func (d OrderAchievementSftpLog) CreateLog(c ctx.Context, data microScrmModel.OrderAchievementSftpLogModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
	return err
}

func (d OrderAchievementSftpLog) UpdateLog(c ctx.Context, query interface{}, data microScrmModel.OrderAchievementSftpLogModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).Where(query).UpdateColumns(data).Error
	return err
}

func (d OrderAchievementSftpLog) DelLog(c ctx.Context, query interface{}) error {
	err := d.dao.WithContext(c).Table(d.tableName).Where(query).Delete(&microScrmModel.OrderAchievementSftpLogModel{}).Error
	return err
}
