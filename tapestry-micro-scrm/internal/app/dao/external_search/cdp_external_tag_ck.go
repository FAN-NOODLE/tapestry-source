package externalSearchDao

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
	"micro-scrm/internal/pkg/db/gormx"
)

type CdpExternalTagCk struct {
	dao            dao.Dao
	tableName      string
	localTableName string
}

func NewCdpExternalTagCk() *CdpExternalTagCk {
	return &CdpExternalTagCk{
		dao:            dao.Dao{DB: gormx.GetInstance()},
		tableName:      externalSearchModel.CdpExternalTag{}.TableName(),
		localTableName: externalSearchModel.CdpExternalTag{}.TableName() + "_local",
	}
}

func (d *CdpExternalTagCk) BatchSaveData(c app.Ctx, enterpriseId int64, newData []externalSearchModel.CdpExternalTag) error {
	if len(newData) > 0 {
		// 批量插入数据
		err := gormx.GetCkInstance().WithContext(c).Transaction(func(tx *gorm.DB) error {
			return tx.Scopes(func(db *gorm.DB) *gorm.DB {
				return db.Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.localTableName))
			}).CreateInBatches(newData, len(newData)).Error
		})
		if err != nil {
			return err
		}
		log.Info(c).Msgf("【CDP_TAG】，插入条数:%d\n", len(newData))
	}

	return nil
}

// MergeData  ck数据合并
func (d *CdpExternalTagCk) MergeData(c app.Ctx, enterpriseId int64) error {
	tableName := fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.localTableName)
	return gormx.GetCkInstance().WithContext(c).Debug().Exec(fmt.Sprintf("OPTIMIZE TABLE %s FINAL DEDUPLICATE", tableName)).Error
}

func (d *CdpExternalTagCk) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (info externalSearchModel.CdpExternalTag, err error) {
	err = gormx.GetCkInstance().WithContext(c).Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)).
		Where(query, args...).Find(&info).Error
	return
}

func (d *CdpExternalTagCk) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]externalSearchModel.CdpExternalTag, error) {
	list := make([]externalSearchModel.CdpExternalTag, 0)
	err := gormx.GetCkInstance().WithContext(c).Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)).
		Where(query, args...).Find(&list).Error
	return list, err
}

// TruncateTable   清空表
func (d *CdpExternalTagCk) TruncateTable(c app.Ctx, enterpriseId int64) error {
	localTableName := fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.localTableName)
	err := gormx.GetCkInstance().WithContext(c).Debug().
		Exec(fmt.Sprintf("TRUNCATE table %s on CLUSTER default_cluster", localTableName)).Error
	if err != nil {
		return err
	}
	tableName := fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)
	err = gormx.GetCkInstance().WithContext(c).Debug().
		Exec(fmt.Sprintf("TRUNCATE table %s on CLUSTER default_cluster", tableName)).Error
	if err != nil {
		return err
	}
	return nil
}
