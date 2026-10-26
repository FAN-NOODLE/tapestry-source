package externalSearchDao

import (
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
	swCrmCustomerModel "micro-scrm/internal/app/models/sw_crm_customer"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"strings"
	"time"
)

type ExternalTagCkSw struct {
	dao            dao.Dao
	tableName      string
	localTableName string
}

func NewExternalTagCkSw() *ExternalTagCkSw {
	return &ExternalTagCkSw{
		dao:            dao.Dao{DB: gormx.GetInstance()},
		tableName:      externalSearchModel.ExternalTagModel{}.TableName(),
		localTableName: externalSearchModel.ExternalTagModel{}.TableName() + "_local",
	}
}

// SaveAtValuesData 顾客季度消费更新
func (d *ExternalTagCkSw) SaveAtValuesData(c app.Ctx, enterpriseId int64, dataList []swCrmCustomerModel.CrmCustomerExtensionModel) error {
	var unionIds []string
	for _, newInfo := range dataList {
		if newInfo.AtValues != "" {
			unionIds = append(unionIds, newInfo.UnionID)
		}
	}
	if len(unionIds) == 0 {
		return nil
	}
	var oldList []externalSearchModel.ExternalTagModel
	err := gormx.GetCkInstance().WithContext(c).Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.tableName)).
		Where("unionid in (?) AND tag_type = ?", unionIds, externalSearchModel.SWAtAmount).Find(&oldList).Error
	if err != nil {
		return err
	}
	var oldKeySlice []string
	// 循环查出来的非邀约数据，让有效数据合并
	for _, oldInfo := range oldList {
		oldKeySlice = append(oldKeySlice, fmt.Sprintf("%s-%s", oldInfo.Unionid, oldInfo.TagKey))
	}
	var newData []externalSearchModel.ExternalTagModel // 待更新的数据
	for _, newInfo := range dataList {
		if newInfo.AtValues != "" {
			newAtSlice := strings.Split(newInfo.AtValues, "|")
			for _, newAtStr := range newAtSlice {
				newTmpAt := strings.Split(newAtStr, ":")
				if len(newTmpAt) == 2 {
					if len(oldKeySlice) == 0 || !helper.InArray(fmt.Sprintf("%s-%s", newInfo.UnionID, newTmpAt[0]), oldKeySlice) {
						newData = append(newData, externalSearchModel.ExternalTagModel{
							Unionid:   newInfo.UnionID,
							TagType:   externalSearchModel.SWAtAmount,
							TagKey:    newTmpAt[0],
							TagValue:  newTmpAt[1],
							CreatedAt: time.Now().Unix(),
							UpdatedAt: time.Now().Unix(),
						})
					}
				}
			}
		}
	}

	if len(newData) > 0 {
		if len(newData) == 1 {
			//这里有gorm针对clickhouse insert的bug，如果clickhouse字段中有数组类型且一次只插一个元素就会报参数数量错误的问题
			newData = append(newData, newData[0])
		}
		// 批量插入数据
		err := gormx.GetCkInstance().WithContext(c).Transaction(func(tx *gorm.DB) error {
			return tx.Scopes(func(db *gorm.DB) *gorm.DB {
				return db.Table(fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, d.localTableName))
			}).Create(&newData).Error
		})
		if err != nil {
			return err
		}
		fmt.Printf("【SW】顾客季度消费更新，更新条数:%d\n", len(newData))
	}
	return nil
}
