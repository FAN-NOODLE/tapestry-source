package membersJob

import (
	"errors"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"go.uber.org/atomic"
	"micro-scrm/internal/app"
	drainageContactDao "micro-scrm/internal/app/dao/drainage"
	membersDao "micro-scrm/internal/app/dao/members"
	externalUserDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	appJob "micro-scrm/internal/app/job"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	cacheSvc "micro-scrm/internal/app/service/cache"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/microwork"
	"micro-scrm/internal/pkg/xjob"
	"strconv"
	"sync"
)

type SyncRelationDataToCdp struct {
	name string
}

func NewSyncRelationDataToCdp() *SyncRelationDataToCdp {
	return &SyncRelationDataToCdp{
		name: "同步外部联系人关系",
	}
}

func (job SyncRelationDataToCdp) Run(c ctx.Context, jobHandler *xjob.JobHandler) error {
	log.Info(c).Msgf("%s,开始执行", job.name)
	enterpriseIds, err := appJob.GetEnterpriseIds(jobHandler)
	if err != nil {
		return err
	}
	for _, enterpriseId := range enterpriseIds {
		err = job.handle(c, enterpriseId)
		if err != nil {
			log.Err(c, err).Msgf("%s,处理企业%d异常", job.name, enterpriseId)
		}
	}

	return nil
}

func (job SyncRelationDataToCdp) handle(c ctx.Context, enterpriseId int64) error {
	// 获取当前企业的所有员工
	config, err := cacheSvc.NewRedisData().GetEnterpriseConfig(c, enterpriseId, "store_name")
	if err != nil {
		return err
	}
	isCoach := false
	if config.Name == "coach" {
		isCoach = true
	}
	isSw := false
	if config.Name == "sw" {
		isSw = true
	}
	if !isCoach && !isSw {
		return nil
	}
	successNum := atomic.NewUint64(0)
	memberDepartmentMap := make(map[string]string)
	storeDepartmentInfo, err := microWorkDao.NewStoreDepartment().GetMany(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
	})
	thirdIdNameMap := make(map[int]string)
	for _, info := range storeDepartmentInfo {
		if _, ok := thirdIdNameMap[int(info.WechatDepartmentID)]; !ok {
			thirdIdNameMap[int(info.WechatDepartmentID)] = info.ThirdStoreCode
		}
	}
	// 获取所有导购所在的门店信息
	memberLists, _ := membersDao.NewMembers().GetAllMemberInfo(c, enterpriseId)
	for _, info := range memberLists {
		if len(info.MemberDepartments) >= 1 {
			for _, department := range info.MemberDepartments {
				if department.DeletedAt == 0 { // 导购门店信息初始化
					departmentName := department.Department.Name
					memberDepartmentMap[info.Userid] = departmentName
					thirdId := department.Department.ThirdID
					if name, ok := thirdIdNameMap[thirdId]; ok {
						memberDepartmentMap[info.Userid] = name
					}

				}
			}
		}
	}
	stateNameMap := make(map[int64]string)
	if isSw {
		stateList, _ := drainageContactDao.NewDrainageContactsDao().GetMany(c, map[string]interface{}{
			"enterprise_id": enterpriseId,
		})
		for _, model := range stateList {
			stateNameMap[model.ID] = model.Name
		}
	}
	redisKey := fmt.Sprintf("sync_relation_id:%d", enterpriseId)
	idStr := app.Redis().Get(c, redisKey).Val()
	id, _ := strconv.ParseInt(idStr, 10, 64)
	cdpClient := microwork.NewCdp()
	p, _ := ants.NewPool(20)
	defer p.Release()
	wg := sync.WaitGroup{}
	for {
		list, err := externalUserDao.NewUserRelation().GetManyPaginationById(c, "enterprise_id= ? and id > ? and relation_status=0 and deleted_at=0", enterpriseId, id)
		if err != nil {
			return err
		}
		index := len(list)
		if index == 0 {
			break
		}
		id = list[index-1].ID
		for _, model := range list {
			wg.Add(1)
			info := model
			if info.Unionid == "" {
				continue
			}
			stateName := ""
			if info.State != "" && isSw {
				stateId, _ := strconv.ParseInt(info.State, 10, 64)
				stateName = stateNameMap[stateId]
			}
			_ = p.Submit(func() {
				defer wg.Done()
				if storeCode, ok := memberDepartmentMap[info.UserID]; ok {
					saRequest := microwork.SaAssignmentReq{
						EnterpriseID: enterpriseId,
						Params: microwork.SaAssignmentParamsReq{
							SaAccount:      info.UserID,
							FirstBindDate:  helper.GetFormatDateByParams(info.CreatedAt),
							LatestBindDate: helper.GetFormatDateByParams(int64(info.Createtime)),
							BindTime:       helper.GetFormatDateByParams(info.CreatedAt),
							UnionId:        info.Unionid,
							StoreCode:      storeCode,
							Source:         stateName,
						},
					}
					_, saErr := cdpClient.SaAssignment(c, saRequest)
					if saErr != nil {
						var resErr *errorx.Error
						errors.As(saErr, &resErr)
						data := microScrmModel.SyncCdpLogModel{
							MicroBaseModel: models.MicroBaseModel{
								BaseModel: models.BaseModel{
									ID: app.UID(),
								},
								EnterpriseID: enterpriseId,
							},
							RelationID: info.ID,
							Result:     fmt.Sprintf("调用CDP接口失败：msg:%v", resErr.String()),
						}
						_ = externalUserDao.NewSyncCdpLogDao().CreateData(c, data)
						// 插入记录数据表，同步cdp失败
					} else {
						successNum.Add(1)
					}
				}

			})
		}
		app.Redis().Set(c, redisKey, id, 0)
	}
	wg.Wait()
	fmt.Printf("推送成功的接口数有%v\n", successNum)
	return nil
}
