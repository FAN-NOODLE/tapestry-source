package microScrmDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type SyncCdpLogDao struct {
	dao       dao.Dao
	model     *microScrmModel.SyncCdpLogModel
	tableName string
}

func NewSyncCdpLogDao() *SyncCdpLogDao {
	return &SyncCdpLogDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		model:     new(microScrmModel.SyncCdpLogModel),
		tableName: microScrmModel.SyncCdpLogModel{}.TableName(),
	}
}

// CreateData 保存记录
func (d *SyncCdpLogDao) CreateData(c app.Ctx, data microScrmModel.SyncCdpLogModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}
