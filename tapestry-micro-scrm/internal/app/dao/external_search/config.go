package externalSearchDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
)

type ConfigDao struct {
	dao.EnterpriseRepo[externalSearchModel.ExternalSearchConfigModel]
	model     *externalSearchModel.ExternalSearchConfigModel
	tableName string
}

func NewConfigDao() *ConfigDao {
	return &ConfigDao{
		EnterpriseRepo: dao.EnterpriseRepo[externalSearchModel.ExternalSearchConfigModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(externalSearchModel.ExternalSearchConfigModel),
		tableName: externalSearchModel.ExternalSearchConfigModel{}.TableName(),
	}
}

func (d *ConfigDao) GetMapByKeys(c ctx.Context, enterpriseId int64, keys []string) (map[string]string, error) {
	var list []externalSearchModel.ExternalSearchConfigModel
	err := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).
		Where("key in ?", keys).
		Find(&list).Error
	var resMap = make(map[string]string, 0)
	for _, model := range list {
		resMap[model.Key] = model.Value
	}
	return resMap, err
}
