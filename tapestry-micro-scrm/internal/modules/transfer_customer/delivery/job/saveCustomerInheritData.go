package job

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/models"
	customerModels "micro-scrm/internal/app/models/customer"
	customerSvc "micro-scrm/internal/app/service/customer"
	"micro-scrm/internal/config"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/redis"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/workwechat"
	"micro-scrm/internal/pkg/xjob"
	"time"
)

// SaveCustomerInheritData 离职待继承人员信息入库
type SaveCustomerInheritData struct {
	inheritService *customerSvc.Inherit
	name           string
}

func NewSaveCustomerInheritData() *SaveCustomerInheritData {
	return &SaveCustomerInheritData{
		name:           fmt.Sprintf("%s::离职待继承外部联系人数据同步Job", config.CfgData.App),
		inheritService: customerSvc.NewInherit(),
	}
}
func (job SaveCustomerInheritData) Run(ctx ctx.Context, jobHandler *xjob.JobHandler) error {

	// 获取所有品牌的企业ID
	enterpriseInfoList, enterpriseInfoErr := orm.BossService.Enterprise.GetAllInfo()
	if enterpriseInfoErr != nil {
		log.Info(ctx).Msgf("%s::获取所有品牌的信息时报错，错误为：%v", job.name, enterpriseInfoErr)
		return enterpriseInfoErr
	}

	// var wg sync.WaitGroup
	for _, enterpriseInfo := range enterpriseInfoList {
		// wg.Add(1)
		job.Exec(ctx, enterpriseInfo.ID)

	}
	// wg.Wait()

	err := FillCustomerInheritData().Run(ctx, jobHandler)
	if err != nil {
		return err
	}
	fmt.Printf("%s,任务执行结束\n", job.name)
	return nil
}

func (job SaveCustomerInheritData) Exec(ctx ctx.Context, enterpriseId int64) {
	syncKey := "sync_inherit_time_" + fmt.Sprintf("%d", enterpriseId)
	redis.Client().SetNX(ctx, syncKey, time.Now().Unix(), 0)
	defer func() {
		redis.Client().Del(ctx, syncKey)
	}()

	params := work.GetUnassignedListRequest{
		PageSize: 500,
	}
	wechat, err := workwechat.NewWorkWeChat(ctx, enterpriseId)
	if err != nil {
		return
	}
	for {
		// 所有目前离职待分配的外部联系人列表
		var newList []string
		var userIds []string
		var externalIds []string
		var newMap = make(map[string]work.UnassignedList, 0)
		infoData, err := wechat.GetUnassignedList(params)
		if err != nil {
			log.Info(ctx).Msgf("%s::企微获取待继承人员数据错误，错误为：%v", job.name, err)
			return
		}
		for _, value := range infoData.Info {
			tmpStr := value.HandoverUserid + value.ExternalUserid + fmt.Sprintf("%d", value.DimissionTime)
			newList = append(newList, tmpStr)
			externalIds = append(externalIds, value.ExternalUserid)
			userIds = append(userIds, value.HandoverUserid)
			newMap[tmpStr] = value
		}
		userIds = helper.ArrayUniq(userIds)
		// 获取库中已有的数据
		list, _ := job.inheritService.GetListByUserIds(ctx, enterpriseId, userIds)
		var dbList []string
		var dbMap = make(map[string]int64, 0)
		for _, db := range list {
			if helper.InArray[string](db.ExternalUserid, externalIds) {
				tmpStr := db.HandoverUserid + db.ExternalUserid + fmt.Sprintf("%d", db.DimissionTime)
				dbList = append(dbList, tmpStr)
				dbMap[tmpStr] = db.ID
			}
		}
		// 如果数据不在数据库中则需要新增，整理新增的数据
		insertList := helper.GenericArrayDiff(newList, dbList)
		var insertData = make([]customerModels.CustomerInheritModel, 0)
		if len(insertList) > 0 {
			for _, t := range insertList {
				if val, ok := newMap[t]; ok {
					insertData = append(insertData, customerModels.CustomerInheritModel{
						MicroBaseModel: models.CreateMicroBase(enterpriseId, app.UID()),
						HandoverUserid: val.HandoverUserid,
						ExternalUserid: val.ExternalUserid,
						DimissionTime:  val.DimissionTime,
					})
				}
			}
		}
		//如果在数据库中则刷新更新时间，整理需要刷新时间的数据
		updateList := helper.ArrayIntersect(newList, dbList)
		updateIds := make([]int64, 0)
		if len(updateList) > 0 {
			for _, t := range updateList {
				if id, ok := dbMap[t]; ok {
					updateIds = append(updateIds, id)
				}
			}
		}
		if len(insertData) > 0 {
			if err = job.inheritService.InsertNewData(ctx, insertData); err != nil {
				log.Info(ctx).Msgf("%s::离职待分配人员数据库写入错误，错误为：%v", job.name, err)
				return
			}
		}
		if len(updateIds) > 0 {
			if err = job.inheritService.UpdateByIds(ctx, enterpriseId, updateIds); err != nil {
				log.Info(ctx).Msgf("%s:已存在的离职待分配人员数据时间更新错误，错误为：%v", job.name, err)
				return
			}
		}

		if infoData.IsLast {
			break
		}
		params.Cursor = infoData.NextCursor
	}

	if err = job.inheritService.DelInvalidData(ctx, enterpriseId); err != nil {
		log.Info(ctx).Msgf("%s::移除失效离职继承数据错误，错误为：%v", job.name, err)
		return
	}

	if err != nil {
		return
	}

}
