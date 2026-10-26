package externalSearchJob

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	externalSearchDao "micro-scrm/internal/app/dao/external_search"
	cacheSvc "micro-scrm/internal/app/service/cache"
	customerSvc "micro-scrm/internal/app/service/customer"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/pkg/notice"
	"micro-scrm/internal/pkg/xjob"
	"strconv"
	"time"
)

type UpdateMaterializeView struct {
	name              string
	intervalInSeconds int64
}

func NewUpdateMaterializeView() *UpdateMaterializeView {
	return &UpdateMaterializeView{
		name:              "外部联系人大宽表数据更新",
		intervalInSeconds: 600, //数据更新间隔
	}
}

func (job *UpdateMaterializeView) Run(c ctx.Context, handler *xjob.JobHandler) error {
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

func (job *UpdateMaterializeView) Handle(c ctx.Context, enterpriseId int64) {
	// 外部联系人数据更新
	externalExtraCkDao := externalSearchDao.NewExternalExtraCk()
	startTime := time.Now().Unix()
	err := externalExtraCkDao.UpdateExternalExtraCk(c, enterpriseId, job.intervalInSeconds)
	if err != nil {
		fmt.Printf("\nexternal_extra数据更新失败，原因：%v", err)
		notice.NewWechat().SendTextMessage(c, fmt.Sprintf("%s异常：%s，%s", job.name, "external_extra", err.Error()))
		return
	}
	fmt.Printf("\n顾客数据更新插入时间:%ds\n", time.Now().Unix()-startTime)
	// ck数据合并
	mergeStartTime := time.Now().Unix()
	err = customerSvc.NewExternalSearch().ExtraDataMerge(c, enterpriseId, false)
	// 外部联系人
	if err != nil {
		fmt.Printf("顾客数据数据合并失败，原因：%v", err)
		return
	}
	fmt.Printf("\n顾客数据合并时间:%ds\n", time.Now().Unix()-mergeStartTime)

	// 外部联系人导购关系数据更新
	externalRelationshipCkDao := externalSearchDao.NewExternalRelationshipCk()
	startTime = time.Now().Unix()
	err = externalRelationshipCkDao.UpdateExternalRelationshipCk(c, enterpriseId, job.intervalInSeconds)
	if err != nil {
		fmt.Printf("\nexternal_relationship数据更新失败，原因：%v", err)
		notice.NewWechat().SendTextMessage(c, fmt.Sprintf("%s异常：%s，%s", job.name, "external_relationship", err.Error()))
		return
	}
	fmt.Printf("\n顾客关系数据更新插入时间:%ds\n", time.Now().Unix()-startTime)
	// ck数据合并
	mergeStartTime = time.Now().Unix()
	err = customerSvc.NewExternalSearch().RelationshipDataMerge(c, enterpriseId, false)
	// 外部联系人
	if err != nil {
		fmt.Printf("\n顾客关系数据合并失败，原因：%v", err)
		return
	}
	fmt.Printf("\n顾客关系数据合并时间:%ds\n", time.Now().Unix()-mergeStartTime)

	fmt.Printf("%s成功", job.name)
}
