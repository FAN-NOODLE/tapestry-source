package externalSearchJob

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	externalSearchDao "micro-scrm/internal/app/dao/external_search"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
	cacheSvc "micro-scrm/internal/app/service/cache"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/pkg/xjob"
)

const (
	extensionBackfillBatchSize = 1000
	extensionBackfillPause     = 20 * time.Millisecond
)

type BackfillExternalExtensions struct {
	customerDao     *externalSearchDao.ExternalCustomerExtension
	relationshipDao *externalSearchDao.ExternalRelationshipExtension
	extraCkDao      *externalSearchDao.ExternalExtraCk
	relationshipCk  *externalSearchDao.ExternalRelationshipCk
}

func NewBackfillExternalExtensions() *BackfillExternalExtensions {
	return &BackfillExternalExtensions{
		customerDao:     externalSearchDao.NewExternalCustomerExtension(),
		relationshipDao: externalSearchDao.NewExternalRelationshipExtension(),
		extraCkDao:      externalSearchDao.NewExternalExtraCk(),
		relationshipCk:  externalSearchDao.NewExternalRelationshipCk(),
	}
}

func (job *BackfillExternalExtensions) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIds, err := extensionBackfillEnterpriseIds(handler)
	if err != nil {
		return err
	}
	// 企业之间串行执行，避免多个千万级回填同时占用 MySQL 和 CK。
	for _, enterpriseId := range enterpriseIds {
		config, configErr := cacheSvc.NewRedisData().GetEnterpriseConfig(c, enterpriseId, "store_name")
		if configErr != nil {
			return fmt.Errorf("获取企业配置失败 enterprise_id=%d: %w", enterpriseId, configErr)
		}
		if config.Name != "coach" {
			continue
		}
		if err = job.backfillEnterprise(c, enterpriseId); err != nil {
			return err
		}
	}
	return nil
}

func extensionBackfillEnterpriseIds(handler *xjob.JobHandler) ([]int64, error) {
	specs := handler.Runner().Specs()
	if len(specs) > 0 {
		ids := make([]int64, 0, len(specs))
		for _, spec := range specs {
			id, err := strconv.ParseInt(spec, 10, 64)
			if err != nil || id <= 0 {
				return nil, fmt.Errorf("非法企业ID: %s", spec)
			}
			ids = append(ids, id)
		}
		return ids, nil
	}
	enterprises, err := orm.BossService.Enterprise.GetAllInfo()
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(enterprises))
	for _, enterprise := range enterprises {
		ids = append(ids, enterprise.ID)
	}
	return ids, nil
}

func (job *BackfillExternalExtensions) backfillEnterprise(c ctx.Context, enterpriseId int64) error {
	fmt.Printf("\n扩展表回填开始 enterprise_id=%d\n", enterpriseId)
	if err := job.backfillCustomerExtensions(c, enterpriseId); err != nil {
		return fmt.Errorf("客户扩展表回填失败 enterprise_id=%d: %w", enterpriseId, err)
	}
	if err := job.backfillRelationshipExtensions(c, enterpriseId); err != nil {
		return fmt.Errorf("关系扩展表回填失败 enterprise_id=%d: %w", enterpriseId, err)
	}
	fmt.Printf("扩展表回填完成 enterprise_id=%d\n", enterpriseId)
	return nil
}

func (job *BackfillExternalExtensions) backfillCustomerExtensions(c ctx.Context, enterpriseId int64) error {
	lastId, err := loadExtensionBackfillCursor(c, enterpriseId, "customer_last_id")
	if err != nil {
		return err
	}
	var scanned, written int64
	for {
		rows, queryErr := job.extraCkDao.GetExtensionBackfillBatch(c, enterpriseId, lastId, extensionBackfillBatchSize)
		if queryErr != nil {
			return queryErr
		}
		if len(rows) == 0 {
			fmt.Printf("客户扩展表回填完成 enterprise_id=%d,last_id=%d,scanned=%d,written=%d\n", enterpriseId, lastId, scanned, written)
			return saveExtensionBackfillCursor(c, enterpriseId, "customer", lastId, "completed")
		}
		data := buildCustomerExtensionBackfillData(rows, time.Now().Unix())
		if err = job.customerDao.UpsertBackfillData(c, enterpriseId, data); err != nil {
			return err
		}
		lastId = rows[len(rows)-1].ID
		scanned += int64(len(rows))
		written += int64(len(data))
		if err = saveExtensionBackfillCursor(c, enterpriseId, "customer", lastId, "running"); err != nil {
			return err
		}
		if err = waitExtensionBackfill(c); err != nil {
			return err
		}
	}
}

func (job *BackfillExternalExtensions) backfillRelationshipExtensions(c ctx.Context, enterpriseId int64) error {
	lastId, err := loadExtensionBackfillCursor(c, enterpriseId, "relationship_last_id")
	if err != nil {
		return err
	}
	var scanned, written int64
	for {
		rows, queryErr := job.relationshipCk.GetExtensionBackfillBatch(c, enterpriseId, lastId, extensionBackfillBatchSize)
		if queryErr != nil {
			return queryErr
		}
		if len(rows) == 0 {
			fmt.Printf("关系扩展表回填完成 enterprise_id=%d,last_id=%d,scanned=%d,written=%d\n", enterpriseId, lastId, scanned, written)
			return saveExtensionBackfillCursor(c, enterpriseId, "relationship", lastId, "completed")
		}
		data := buildRelationshipExtensionBackfillData(rows, time.Now().Unix())
		if err = job.relationshipDao.UpsertData(c, enterpriseId, data); err != nil {
			return err
		}
		lastId = rows[len(rows)-1].ID
		scanned += int64(len(rows))
		written += int64(len(data))
		if err = saveExtensionBackfillCursor(c, enterpriseId, "relationship", lastId, "running"); err != nil {
			return err
		}
		if err = waitExtensionBackfill(c); err != nil {
			return err
		}
	}
}

func buildCustomerExtensionBackfillData(rows []externalSearchModel.ExternalExtraModel, now int64) []externalSearchModel.ExternalCustomerExtensionModel {
	dataMap := make(map[string]externalSearchModel.ExternalCustomerExtensionModel)
	for _, row := range rows {
		if row.UserUnionid == "" || !hasCustomerExtensionValue(row) {
			continue
		}
		item := dataMap[row.UserUnionid]
		item.Unionid = row.UserUnionid
		item.ValidReserveAt = maxInt64(item.ValidReserveAt, int64(row.ValidReserveAt))
		item.SuccessReserveAt = maxInt64(item.SuccessReserveAt, int64(row.SuccessReserveAt))
		item.TaskFinishAt = maxInt64(item.TaskFinishAt, int64(row.TaskFinishAt))
		item.TaskCreatedAt = maxInt64(item.TaskCreatedAt, row.TaskCreatedAt)
		if int64(row.CouponUpdatedAt) >= item.CouponUpdatedAt {
			item.CouponStatus = row.CouponStatus
			item.CouponBy30daysStatus = row.CouponBy30daysStatus
			item.BirthdayCouponStatus = row.BirthdayCouponStatus
			item.RenewMembershipCouponStatus = row.RenewMembershipCouponStatus
			item.WelcomeCouponStatus = row.WelcomeCouponStatus
			item.UpgradeCouponStatus = row.UpgradeCouponStatus
			item.CouponUpdatedAt = int64(row.CouponUpdatedAt)
		}
		item.CreatedAt = now
		item.UpdatedAt = maxInt64(item.UpdatedAt, int64(row.UpdatedAt))
		if item.UpdatedAt == 0 {
			item.UpdatedAt = now
		}
		dataMap[row.UserUnionid] = item
	}
	result := make([]externalSearchModel.ExternalCustomerExtensionModel, 0, len(dataMap))
	for _, item := range dataMap {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Unionid < result[j].Unionid })
	return result
}

func hasCustomerExtensionValue(row externalSearchModel.ExternalExtraModel) bool {
	return row.ValidReserveAt > 0 || row.SuccessReserveAt > 0 || row.TaskFinishAt > 0 || row.TaskCreatedAt > 0 ||
		row.CouponUpdatedAt > 0 ||
		row.CouponStatus > 0 || row.CouponBy30daysStatus > 0 || row.BirthdayCouponStatus > 0 ||
		row.RenewMembershipCouponStatus > 0 || row.WelcomeCouponStatus > 0 || row.UpgradeCouponStatus > 0
}

func buildRelationshipExtensionBackfillData(rows []externalSearchModel.ExternalRelationshipModel, now int64) []externalSearchModel.ExternalRelationshipExtensionModel {
	dataMap := make(map[int64]externalSearchModel.ExternalRelationshipExtensionModel)
	for _, row := range rows {
		if row.ID <= 0 || row.InteractTime <= 0 && row.ActiveTime <= 0 {
			continue
		}
		item := dataMap[row.ID]
		item.RelationID = row.ID
		item.Unionid = row.Unionid
		item.UserID = row.UserID
		item.InteractTime = maxInt64(item.InteractTime, row.InteractTime)
		item.ActiveTime = maxInt64(item.ActiveTime, row.ActiveTime)
		item.CreatedAt = now
		item.UpdatedAt = maxInt64(item.UpdatedAt, int64(row.UpdatedAt))
		if item.UpdatedAt == 0 {
			item.UpdatedAt = now
		}
		dataMap[row.ID] = item
	}
	result := make([]externalSearchModel.ExternalRelationshipExtensionModel, 0, len(dataMap))
	for _, item := range dataMap {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].RelationID < result[j].RelationID })
	return result
}

func maxInt64(left, right int64) int64 {
	if left > right {
		return left
	}
	return right
}

func extensionBackfillProgressKey(enterpriseId int64) string {
	return fmt.Sprintf("customer:external-extension-backfill:%d", enterpriseId)
}

func loadExtensionBackfillCursor(c ctx.Context, enterpriseId int64, field string) (int64, error) {
	value, err := app.Redis().HGet(c, extensionBackfillProgressKey(enterpriseId), field).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	lastId, err := strconv.ParseInt(value, 10, 64)
	if err != nil || lastId < 0 {
		return 0, fmt.Errorf("非法回填游标 enterprise_id=%d,field=%s,value=%s", enterpriseId, field, value)
	}
	return lastId, nil
}

func saveExtensionBackfillCursor(c ctx.Context, enterpriseId int64, prefix string, lastId int64, status string) error {
	return app.Redis().HSet(c, extensionBackfillProgressKey(enterpriseId), map[string]interface{}{
		prefix + "_last_id":    lastId,
		prefix + "_status":     status,
		prefix + "_updated_at": time.Now().Unix(),
	}).Err()
}

func waitExtensionBackfill(c ctx.Context) error {
	timer := time.NewTimer(extensionBackfillPause)
	defer timer.Stop()
	select {
	case <-c.Done():
		return c.Err()
	case <-timer.C:
		return nil
	}
}
