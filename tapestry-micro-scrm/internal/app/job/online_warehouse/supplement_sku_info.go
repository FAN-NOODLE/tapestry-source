package online_warehouse

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/lock"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	onlineWarehouseDao "micro-scrm/internal/app/dao/online_warehouse"
	"micro-scrm/internal/app/http/request/source"
	cacheSvc "micro-scrm/internal/app/service/cache"
	sourceService "micro-scrm/internal/app/service/source"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/xjob"
	"strconv"
	"strings"
	"time"
)

type SupplementSkuInfo struct {
	name string
}

func NewSupplementSkuInfo() *SupplementSkuInfo {
	return &SupplementSkuInfo{
		name: "给老的商品通过ross接口补充sku数据",
	}
}

func (job *SupplementSkuInfo) Run(c ctx.Context, handler *xjob.JobHandler) error {
	specsArgs := handler.Runner().Specs()
	enterpriseIds := make([]int64, 0)
	if len(specsArgs) > 0 {
		for _, arg := range specsArgs {
			enterpriseId, _ := strconv.ParseInt(arg, 10, 64)
			enterpriseIds = append(enterpriseIds, enterpriseId)
		}
	} else {
		enterprises, err := orm.BossService.Enterprise.GetAllInfo()
		if err != nil || len(enterprises) == 0 {
			return err
		}
		for _, enterprise := range enterprises {
			enterpriseIds = append(enterpriseIds, enterprise.ID)
		}
	}

	var wg threading.RoutineGroup
	for _, id := range enterpriseIds {
		enterpriseId := id
		cacheRes, err := cacheSvc.NewRedisData().GetEnterpriseConfig(c, enterpriseId, "store_name")
		if err != nil {
			continue
		}
		if cacheRes.Name != "coach" {
			continue
		}
		wg.Run(func() {
			job.Handle(c, enterpriseId)
		})
	}

	wg.Wait()

	return nil
}

func (job *SupplementSkuInfo) Handle(c ctx.Context, enterpriseId int64) error {
	// 添加分布式锁避免重复发货重复执行
	key := helper.Md5V(job.name + fmt.Sprintf("%d", enterpriseId))
	mtx, err := lock.Mutex(
		c,
		lock.Key(key),
		lock.TTL(60*time.Second),
		lock.Retry(lock.NoRetry()),
	)
	if err != nil {
		return err
	}
	err = mtx.Lock()
	if err != nil {
		log.Info(c).Msgf("%s::获取锁失败，执行太频繁", job.name)
		return err
	}
	defer mtx.Unlock()
	skuDao := onlineWarehouseDao.NewOrderSkuDao()
	list, _ := skuDao.GetMany(c, enterpriseId, "dept=? AND upc_code = ?", "", "")
	if len(list) > 0 {
		for _, model := range list {
			data, _ := sourceService.NewBrandService().SearchList(c, source.BrandListReq{
				BrandCacheReq: source.BrandCacheReq{
					EnterpriseId:            enterpriseId,
					SkuItem:                 strings.ToUpper(model.SkuInfo),
					Accurate:                true,
					IncludeNoneImageProduct: true,
				},
			})
			if len(data) > 0 {
				for _, row := range data {
					var skuNumber string
					var imageUrl string
					var upcCode string
					if len(row.UpcList) > 0 {
						skuNumber = row.UpcList[0].SkuNumber
						upcCode = row.UpcList[0].UpcCode
					}
					if len(row.BrandSourceList) > 0 {
						imageUrl = row.BrandSourceList[0].URL
					}
					_ = skuDao.UpdateData(c, enterpriseId, map[string]interface{}{
						"id": model.ID,
					}, map[string]interface{}{
						"dept":        row.Department,
						"style_color": row.StyleCode + " " + row.Color,
						"sku_number":  skuNumber,
						"sku_name":    row.ProductName,
						"image_url":   imageUrl,
						"upc_code":    upcCode,
					})
					break
				}
			}
		}
	}
	return nil
}
