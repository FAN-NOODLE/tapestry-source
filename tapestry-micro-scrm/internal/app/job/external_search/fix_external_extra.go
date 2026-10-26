package externalSearchJob

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	externalSearchDao "micro-scrm/internal/app/dao/external_search"
	cacheSvc "micro-scrm/internal/app/service/cache"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/pkg/notice"
	"micro-scrm/internal/pkg/xjob"
	"strconv"
)

type FixExternalExtra struct {
	name string
}

func NewFixExternalExtra() *FixExternalExtra {
	return &FixExternalExtra{
		name: "补偿修复extra和relationShip遗漏的数据",
	}
}

func (job *FixExternalExtra) Run(c ctx.Context, handler *xjob.JobHandler) error {
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

func (job *FixExternalExtra) Handle(c ctx.Context, enterpriseId int64) {
	// 外部联系人数据更新
	externalExtraCkDao := externalSearchDao.NewExternalExtraCk()
	err := externalExtraCkDao.FixExternalExtra(c, enterpriseId)
	if err != nil {
		fmt.Printf("\nexternal_extra数据修复失败，原因：%v", err)
		notice.NewWechat().SendTextMessage(c, fmt.Sprintf("%s异常：%s，%s", job.name, "external_extra", err.Error()))
		return
	}

	relationCkDao := externalSearchDao.NewExternalRelationshipCk()
	err = relationCkDao.FixRelationShip(c, enterpriseId)
	if err != nil {
		fmt.Printf("\nrelation_ship数据修复失败，原因：%v", err)
		notice.NewWechat().SendTextMessage(c, fmt.Sprintf("%s异常：%s，%s", job.name, "relation_ship", err.Error()))
		return
	}

	fmt.Printf("%s成功", job.name)
}
