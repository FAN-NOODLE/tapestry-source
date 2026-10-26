package excitationDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/pkg/db/gormx"
)

type TaskRecord struct {
	dao        dao.Dao
	model      *excitationModels.ExcitationTaskRecordModel
	existTable bool
}

func NewTaskRecord() *TaskRecord {
	return &TaskRecord{
		dao:        dao.Dao{DB: gormx.GetInstance()},
		model:      new(excitationModels.ExcitationTaskRecordModel),
		existTable: false,
	}
}

func (d *TaskRecord) GetOne(c app.Ctx, query interface{}, args ...interface{}) (excitationModels.ExcitationTaskRecordModel, error) {
	var result excitationModels.ExcitationTaskRecordModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *TaskRecord) CreateData(c app.Ctx, data excitationModels.ExcitationTaskRecordModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *TaskRecord) CreateInBatchesData(c app.Ctx, data []excitationModels.ExcitationTaskRecordModel) error {
	err := d.dao.WithContext(c).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *TaskRecord) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}
